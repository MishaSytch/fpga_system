package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const PATH = "B:/Измерение_зависимости_отклика_от_растояния_излучателя"

type Measurement struct {
	Name      string
	Columns   []string    // Заголовки столбцов
	Points    []string    // Единицы измерения
	AvgData   [][]float64 // Усредненные данные
	FileCount int         // Количество обработанных файлов
}

func main() {
	rootDir := PATH
	measurements, err := ProcessRootDir(rootDir)
	if err != nil {
		log.Fatalf("Error processing directory: %v", err)
	}

	// Выводим результаты
	for _, m := range measurements {
		file, err := os.Create(PATH + "/" + m.Name + ".csv")
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer file.Close()

		m.Columns[3] = strings.ReplaceAll(m.Columns[3], ";", ",")

		file.WriteString(fmt.Sprintf(strings.Join(m.Columns, ";")) + "\n")
		file.WriteString(fmt.Sprintf(strings.Join(m.Points, ";")) + "\n")
		for row := range m.AvgData {
			file.WriteString(
				fmt.Sprintf("%v;%v;%v;%v\n",
					m.AvgData[row][0],
					m.AvgData[row][1],
					m.AvgData[row][2],
					m.AvgData[row][3],
				))
		}
	}
}

// processRootDir обрабатывает корневую директорию
func ProcessRootDir(rootPath string) ([]*Measurement, error) {
	var measurements []*Measurement

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем корневую директорию и файлы
		if path == rootPath || !info.IsDir() {
			return nil
		}

		// Проверяем, есть ли в директории CSV файлы
		hasCSV := false
		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
				hasCSV = true
				break
			}
		}

		// Если есть CSV файлы - обрабатываем как директорию с измерениями
		if hasCSV {
			measurement, err := processMeasurementDir(path)
			if err != nil {
				return fmt.Errorf("error processing %s: %v", path, err)
			}
			measurements = append(measurements, measurement)
			return filepath.SkipDir // Не заходим в поддиректории
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return measurements, nil
}

func readCSV(filePath string) ([][]float64, []string, []string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Читаем заголовки
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, nil, err
	}

	points, err := reader.Read()
	if err != nil {
		return nil, nil, nil, err
	}

	var data [][]float64
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, nil, err
		}

		row := make([]float64, len(headers))
		for i, v := range record {
			if v == "не число" {
				continue
			}
			v = strings.ReplaceAll(v, ",", ".")
			f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("error parsing float in file %s: %v", filePath, err)
			}
			row[i] = f
		}
		data = append(data, row)
	}

	return data, headers, points, nil
}

// processMeasurementDir обрабатывает директорию с измерениями
func processMeasurementDir(dirPath string) (*Measurement, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var csvFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			csvFiles = append(csvFiles, filepath.Join(dirPath, file.Name()))
		}
	}

	if len(csvFiles) == 0 {
		return nil, fmt.Errorf("no CSV files found in directory %s", dirPath)
	}

	var sumData [][]float64
	var headers []string
	var points []string
	var totalRows int

	// Читаем и суммируем данные из всех CSV файлов
	for i, file := range csvFiles {
		data, h, p, err := readCSV(file)
		if err != nil {
			return nil, fmt.Errorf("error reading %s: %v", file, err)
		}

		// Проверяем согласованность заголовков
		if i == 0 {
			headers = h
			points = p
		} else if !compareHeaders(headers, h) {
			return nil, fmt.Errorf("headers mismatch between files in %s", dirPath)
		}

		// Инициализируем sumData при первом проходе
		if i == 0 {
			sumData = make([][]float64, len(data))
			for j := range data {
				sumData[j] = make([]float64, len(data[j]))
				copy(sumData[j], data[j])
			}
			totalRows = len(data)
		} else {
			// Проверяем, что количество строк совпадает
			if len(data) != totalRows {
				totalRows = min(totalRows, len(data))
				sumData = sumData[:totalRows]
			}
			// Суммируем данные
			for j := range totalRows {
				for k := range data[j] {
					sumData[j][k] += data[j][k]
				}
			}
		}
	}

	// Вычисляем средние значения
	avgData := make([][]float64, len(sumData))
	for i := range len(sumData) {
		for j := range sumData[0] {
			avgData[i] = append(avgData[i], sumData[i][j]/float64(len(files)))
		}
	}

	return &Measurement{
		Name:      filepath.Base(dirPath),
		Columns:   headers,
		Points:    points,
		AvgData:   avgData,
		FileCount: len(csvFiles),
	}, nil
}

// compareHeaders сравнивает два набора заголовков
func compareHeaders(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
