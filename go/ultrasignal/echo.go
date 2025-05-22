package ultrasignal

import "math"

// DetectEchoes находит временные координаты (индексы) эхо-сигналов в переданном сигнале.
//
// Эхо определяется как локальный всплеск амплитуды, превышающий заданный порог (threshold).
//
// Условие обнаружения (пороговая фильтрация по амплитуде):
//
//	|signal[i]| > threshold
//
// Параметры:
//   - signal: одномерный временной сигнал (обычно результат корреляции).
//   - threshold: числовое значение порога, выше которого сигнал считается "эхом".
//
// Возвращает:
//   - indices: срез индексов, в которых обнаружено превышение порога (возможные эхо-отклики).
func DetectEchoes(signal []float64, threshold float64) []int {
	var indices []int
	for i, val := range signal {
		if math.Abs(val) > threshold {
			indices = append(indices, i)
		}
	}
	return indices
}
