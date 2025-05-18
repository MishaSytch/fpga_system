package main

import (
	"fpga-ultrasound-go/memory"
	"fpga-ultrasound-go/storage"
	"fpga-ultrasound-go/ultrasignal"
	"log"
	"math"
	"time"
)

const (
	Path                = "/home/root/data"
	FileWithFreq        = "[AFC]"
	FileWithTime        = "[Time]"
	FilterWindow        = 5
	CurrentSampleRateHz = 1e5                      // 0.1 МГц
	SampleRateHz        = CurrentSampleRateHz * 10 // 10 МГц — типичная частота дискретизации ультразвука
	FFTKernelSize       = 1000                     // размер окна для FFT
	FIRKernelSize       = 101                      // нечётное число
	LowCutoffFreq       = 1e-3                     // 0.001 Гц
	HighCutoffFreq      = 1e6                      // 1 МГц
	EchoThreshold       = 0.6                      // порог обнаружения эха
	Threshold           = 0.5
	Thickness           = 10.0 // Толщина образца в мм
	Mode                = "A0" // Модальный режим ("A0" или "S0")
)

func main() {
	log.Println("🚀 Starting FPGA Ultrasound Data Collector...")
	var dataBuffer []float64
	var raw = make(chan []float64)

	go func(raw chan []float64) {
		log.Println("Чтение данных")
		for {
			var data []float64
			err := memory.ReadFrame("", CurrentSampleRateHz, &data)
			if err != nil {
				log.Printf("❌ Memory read error: %v", err)
				break
			}
			raw <- data
			if err := storage.SaveSample(Path+FileWithTime+" RAW_result.csv", data); err != nil {
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
	data = ultrasignal.TrashHoldFilter(data, Threshold)

	log.Println("1️⃣ Сглаживание с использованием скользящего среднего")
	smoothed := ultrasignal.MovingAverage(data, FilterWindow)

	log.Println("2️⃣ Применение фильтра (полосовой фильтр)")
	kernel := ultrasignal.FIRBandPassKernel(FIRKernelSize, LowCutoffFreq, HighCutoffFreq, SampleRateHz)
	filteredSignal := ultrasignal.BandPassFilter(smoothed, kernel)
	if err := storage.SaveSample(Path+FileWithTime+" FIR_result.csv", filteredSignal); err != nil {
		log.Printf("❌ FIR save error: %v", err)
	}

	log.Println("3️⃣ Вычисление АЧХ фильтра")
	freqsAFC, afc := ultrasignal.ComputeAFC(filteredSignal, SampleRateHz)
	if err := storage.SaveSpectrum(Path+FileWithFreq+" АЧХ фильтра.csv", freqsAFC, afc); err != nil {
		log.Printf("❌ AFC save error: %v", err)
	}

	log.Println("4️⃣ Расчёт огибающей через Гильберта")
	envelopeHilbert := ultrasignal.ComputeEnvelopeHilbert(filteredSignal)
	if err := storage.SaveSample(Path+FileWithTime+" Огибающая через Гильберта.csv", envelopeHilbert); err != nil {
		log.Printf("❌ Hilbert envelope save error: %v", err)
	}

	log.Println("5️⃣ Обнаружение эхо-сигналов и расчет времени полета")
	echoIndices := ultrasignal.DetectEchoes(envelopeHilbert, EchoThreshold)
	tof := ultrasignal.GetTimeOfFlight(echoIndices, SampleRateHz)
	log.Printf("⏱️ Time of Flight: %.9f секунд", tof)

	log.Println("6️⃣ Расчёт спектра с использованием FFT")
	windowed := ultrasignal.HammingWindow(filteredSignal[:min(len(filteredSignal), FFTKernelSize)])
	frequencies, spectrum := ultrasignal.ComputeFFTLog(windowed, SampleRateHz, math.Pow(10.0, -3.0), math.Pow(10.0, 6), FFTKernelSize)
	if err := storage.SaveSpectrum(Path+FileWithFreq+" Спектр сигнала.csv", frequencies, spectrum); err != nil {
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

	if err := storage.SaveSample(Path+FileWithTime+" PhaseVelocity"+".csv", phaseVel); err != nil {
		log.Printf("❌ Phase velocity save error: %v", err)
	}
	if err := storage.SaveSample(Path+FileWithTime+" GroupVelocity"+".csv", groupVel); err != nil {
		log.Printf("❌ Group velocity save error: %v", err)
	}

	log.Println("Анализ проведен")
	time.Sleep(ultrasignal.FreqToTime(CurrentSampleRateHz))
}
