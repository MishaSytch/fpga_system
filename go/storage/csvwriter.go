package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func SaveSample(filename string, data []float64) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open csv failed: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, v := range data {
		currentTime := time.Now().UTC().Format(time.RFC3339Nano)
		value := fmt.Sprintf("%0.5f", v)
		record := []string{currentTime, value}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("write csv failed: %w", err)
		}
	}

	return nil
}

func SaveSpectrum(filename string, frequencies, spectrum []float64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Данные
	for i := 0; i < len(frequencies); i++ {
		_, err := file.WriteString(fmt.Sprintf("%.6f,%.6f\n", frequencies[i], spectrum[i]))
		if err != nil {
			return err
		}
	}
	return nil
}
