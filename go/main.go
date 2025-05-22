package main

import (
	"fmt"
	"fpga-ultrasound-go/memory"
	"fpga-ultrasound-go/storage"
	"fpga-ultrasound-go/ultrasignal"
	"log"
	"math"
	"os"
	"runtime"
	"time"
)

const (
	FileWithFreq        = "[AFC]"
	FileWithTime        = "[Time]"
	FilterWindow        = 5
	CurrentSampleRateHz = 1e5                      // 0.1 МГц
	SampleRateHz        = CurrentSampleRateHz * 10 // 10 МГц — типичная частота дискретизации ультразвука
	FFTKernelSize       = 1000                     // Размер окна для FFT
	FIRKernelSize       = 101                      // Нечётное число
	LowCutoffFreq       = 1e-3                     // 0.001 Гц
	HighCutoffFreq      = 1e6                      // 1 МГц
	EchoThreshold       = 0.6                      // Порог обнаружения эха
	Threshold           = 0.5
	Thickness           = 10.0 // Толщина образца в мм
	Mode                = "A0" // Модальный режим ("A0" или "S0")
)

func main() {
	logFile, err := logSettings()
	if err != nil {
		log.Printf(err.Error())
	}
	defer logFile.Close()

	go func() {
		var m runtime.MemStats
		for {
			runtime.ReadMemStats(&m)
			fmt.Printf("Количество горутин: %d\n", runtime.NumGoroutine())
			fmt.Printf("Текущий объем занятой памяти: %d\n", bToMb(m.Alloc))
			fmt.Printf("Всего выделенно памяти во время запуска: %d\n", bToMb(m.TotalAlloc))
			fmt.Printf("Объем памяти, полученный от операционной системы: %d\n", bToMb(m.TotalAlloc))
			fmt.Printf("Количество cpu: %d\n", runtime.NumCPU())
		}
	}()
	log.Println("🚀 Starting FPGA Ultrasound Data Collector...")
	var dataBuffer []float64
	var raw = make(chan []float64)

	//if runtime.GOOS == "windows" {
	//	Path = memory.CurrentData(CurrentSampleRateHz)
	//}

	go func(raw chan []float64) {
		log.Println("Чтение данных")
		for {
			var data []float64
			err := memory.ReadFrame(Path, CurrentSampleRateHz, &data)
			if err != nil {
				log.Printf("❌ Memory read error: %v", err)
				break
			}
			raw <- data
			if err := storage.SaveSample("./"+FileWithTime+"_RAW_result.csv", data); err != nil {
				log.Printf("❌ raw save error: %v", err)
			}
		}
		return
	}(raw)

	log.Println("Сбор данных")
loop:
	for {
		select {
		case data := <-raw:
			dataBuffer = append(dataBuffer, data...)
			log.Printf("Длина полученных данных: %+v", len(dataBuffer))
			if len(dataBuffer) >= FFTKernelSize {
				break loop
			}
		}
	}

	processing(dataBuffer)
}

func processing(data []float64) {
	FilePath := "./"

	data = ultrasignal.ThresholdFilter(data, Threshold)

	log.Println("1️⃣ Сглаживание с использованием скользящего среднего")
	smoothed := ultrasignal.MovingAverage(data, FilterWindow)

	log.Println("2️⃣ Применение фильтра (полосовой фильтр)")
	kernel := ultrasignal.FIRBandPassKernel(FIRKernelSize, LowCutoffFreq, HighCutoffFreq, SampleRateHz)
	filteredSignal := ultrasignal.BandPassFilter(smoothed, kernel)
	if err := storage.SaveSample(FilePath+FileWithTime+"_FIR_result.csv", filteredSignal); err != nil {
		log.Printf("❌ FIR save error: %v", err)
	}

	log.Println("3️⃣ Вычисление АЧХ фильтра")
	freqsAFC, afc := ultrasignal.ComputeAFC(filteredSignal, SampleRateHz)
	if err := storage.SaveSpectrum(FilePath+FileWithFreq+"_filter_frequency_response.csv", freqsAFC, afc); err != nil {
		log.Printf("❌ AFC save error: %v", err)
	}

	log.Println("4️⃣ Расчёт огибающей через Гильберта")
	envelopeHilbert := ultrasignal.ComputeEnvelopeHilbert(filteredSignal)
	if err := storage.SaveSample(FilePath+FileWithTime+"_Envelope_via_Hilbert.csv", envelopeHilbert); err != nil {
		log.Printf("❌ Hilbert envelope save error: %v", err)
	}

	log.Println("5️⃣ Обнаружение эхо-сигналов и расчет времени полета")
	echoIndices := ultrasignal.DetectEchoes(envelopeHilbert, EchoThreshold)
	tof := ultrasignal.GetTimeOfFlight(echoIndices, SampleRateHz)
	log.Printf("⏱️ Time of Flight: %.9f секунд", tof)

	log.Println("6️⃣ Расчёт спектра с использованием FFT")
	windowed := ultrasignal.HammingWindow(filteredSignal[:min(len(filteredSignal), FFTKernelSize)])
	frequencies, spectrum := ultrasignal.ComputeFFTLog(windowed, SampleRateHz, math.Pow(10.0, -3.0), math.Pow(10.0, 6), FFTKernelSize)
	if err := storage.SaveSpectrum(FilePath+FileWithFreq+"_Signal_spectrum.csv", frequencies, spectrum); err != nil {
		log.Printf("❌ Spectrum save error: %v", err)
	}

	log.Println("7️⃣ Расчёт фазовой и групповой скорости для каждой частоты")
	var phaseVel []float64
	var groupVel []float64
	for _, freq := range frequencies {
		phaseVelC := ultrasignal.PhaseVelocity(freq, Thickness, Mode)
		groupVelC := ultrasignal.GroupVelocity(freq, Thickness, Mode)

		phaseVel = append(phaseVel, phaseVelC)
		groupVel = append(groupVel, groupVelC)
	}

	if err := storage.SaveSample(FilePath+FileWithTime+"_PhaseVelocity"+".csv", phaseVel); err != nil {
		log.Printf("❌ Phase velocity save error: %v", err)
	}
	if err := storage.SaveSample(FilePath+FileWithTime+"_GroupVelocity"+".csv", groupVel); err != nil {
		log.Printf("❌ Group velocity save error: %v", err)
	}

	log.Println("Анализ проведен")
	time.Sleep(ultrasignal.FreqToTime(CurrentSampleRateHz))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func logSettings() (*os.File, error) {
	logFile, err := os.OpenFile("ultrasound_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("❌ Не удалось открыть лог-файл: %v\n", err)
		return &os.File{}, err
	}
	log.SetOutput(logFile)

	return logFile, nil
}
