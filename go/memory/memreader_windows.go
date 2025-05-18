//go:build windows
// +build windows

package memory

import (
	"bufio"
	"fmt"
	"fpga-ultrasound-go/ultrasignal"
	"log"
	"math/rand"
	"os"
	"time"
)

const FrameSize = 1024
const PATH1 = "B:\\fpga_data.bin"
const PATH2 = "B:\\loop_fpga_data.bin"

var lastSize int64 = 0

func CurrentData(frequency float64) string {
	go func(path string) {
		file, _ := os.Create(path)
		defer file.Close()

		ticker := time.NewTicker(ultrasignal.FreqToTime(frequency))
		defer ticker.Stop()

		counter := 0
		log.Println("Симуляция данных")
		for range ticker.C {
			var data float64
			switch {
			case counter < 50:
				data = (rand.Float64() - 0.5) * 0.05 // Центрированный шум: -0.025 до 0.025
			case counter == 50:
				data = 1.0
			case counter == 70:
				data = 0.6
			case counter > 1000:
				counter = -1
			default:
				data = (rand.Float64() - 0.5) * 0.05 // Центрированный шум
			}
			file.WriteString(fmt.Sprintf("%.3f\n", data))
			counter++
		}
		log.Println("Симуляция данных прекращена")
	}(PATH2)

	return PATH2
}

func CreateFPGAData() string {
	file, _ := os.Create(PATH1)

	defer file.Close()

	for range FrameSize {
		data := rand.Float64()
		file.Write([]byte{byte(data)})
	}
	return PATH1
}

func ReadFrame(path string, freq float64, data *[]float64) error {
	for {
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Printf("❌ Error getting file info: %v", err)
			continue
		}

		if fileInfo.Size() == lastSize {
			time.Sleep(ultrasignal.FreqToTime(freq))
			continue
		}
		break
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var val float64
		if _, err := fmt.Sscanf(scanner.Text(), "%f", &val); err == nil {
			*data = append(*data, val)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("ошибка очистки файла: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("ошибка перемещения указателя: %v", err)
	}

	return nil
}
