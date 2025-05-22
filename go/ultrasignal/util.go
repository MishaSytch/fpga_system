package ultrasignal

import (
	"math"
	"time"
)

// Convolve вычисляет свёртку входного сигнала с ядром фильтра.
//
// Свёртка — это операция, которая для каждого сдвига i вычисляет сумму произведений
// элементов сигнала и ядра, смещённых относительно друг друга.
//
// Формула свёртки (однонаправленная, без обращения ядра):
//
//	y[i] = Σ_{j=0}^{k-1} x[i-j] * h[j],    для i-j >= 0
//
// где x — сигнал длины n, h — ядро длины k, y — выходной сигнал длины n.
//
// Здесь свёртка реализована с условием i-j >= 0, что соответствует causal свёртке,
// без учета "отрицательных" индексов.
func Convolve(signal, kernel []float64) []float64 {
	n := len(signal)
	k := len(kernel)
	output := make([]float64, n)
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			if i-j >= 0 {
				output[i] += signal[i-j] * kernel[j]
			}
		}
	}
	return output
}

// Abs вычисляет амплитуду (модуль) комплексного числа с действительной и мнимой частью.
//
// Формула:
//
//	|z| = sqrt(Re(z)^2 + Im(z)^2)
func Abs(real, imag float64) float64 {
	return math.Sqrt(real*real + imag*imag)
}

// min возвращает минимальное из двух целых чисел
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max возвращает максимальное из двух целых чисел
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// FreqToTime преобразует частоту (Гц) в период времени (time.Duration).
//
// Период рассчитывается как T = 1 / f (секунды).
// Для точного представления используется преобразование секунд в наносекунды.
//
// Если частота <= 0, возвращается 0.
func FreqToTime(frequency float64) time.Duration {
	if frequency <= 0 {
		return 0
	}
	periodSec := 1.0 / frequency
	// Переводим секунды в наносекунды
	nano := periodSec * 1e9
	return time.Duration(nano) * time.Nanosecond
}
