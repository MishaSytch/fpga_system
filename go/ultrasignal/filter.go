package ultrasignal

import (
	"gonum.org/v1/gonum/floats"
	"math"
)

// BandPassFilter применяет фильтр (FIR kernel) к сигналу с помощью линейной свёртки.
// Математически: y[n] = ∑ₖ x[n-k]·h[k]
// где x[n] — входной сигнал, h[k] — импульсная характеристика (kernel), y[n] — отфильтрованный сигнал.
func BandPassFilter(input, kernel []float64) []float64 {
	return Convolve(input, kernel)
}

// ComputeAFC рассчитывает амплитудно-частотную характеристику (АЧХ) фильтра.
// Использует DFT (FFT) ядра фильтра (kernel).
//
// Возвращает частоты и амплитуды для половины спектра (до Nyquist).
//
// Формула: H(f) = |FFT{h[n]}|
// Нормировка: амплитуда делится на длину ядра и удваивается (кроме DC и Nyquist).
func ComputeAFC(kernel []float64, sampleRate float64) ([]float64, []float64) {
	return ComputeFFT(kernel, sampleRate)
}

// FIRBandPassKernel генерирует ядро КИХ-фильтра полосового пропускания
// методом оконного проектирования на основе sinc-функции и окна Хэмминга.
//
// Формула (идеальный фильтр):
//
//	h[n] = (sin(ω₂·n) - sin(ω₁·n)) / (π·n), где ω₁ = 2π·f_low / f_s, ω₂ = 2π·f_high / f_s
//
// Корректировка при n = 0: h[0] = (ω₂ - ω₁)/π
//
// Далее применяется окно Хэмминга:
//
//	w[n] = 0.54 - 0.46·cos(2π·n / (N - 1))
//
// Итоговое: h[n] *= w[n]
//
// Выходное ядро нормируется по сумме.
func FIRBandPassKernel(size int, lowCutoff, highCutoff, sampleRate float64) []float64 {
	kernel := make([]float64, size)
	mid := size / 2
	omega1 := 2 * math.Pi * lowCutoff / sampleRate
	omega2 := 2 * math.Pi * highCutoff / sampleRate

	for i := 0; i < size; i++ {
		n := float64(i - mid)
		if n == 0 {
			kernel[i] = (omega2 - omega1) / math.Pi
		} else {
			kernel[i] = (math.Sin(omega2*n) - math.Sin(omega1*n)) / (math.Pi * n)
		}
		// Окно Хэмминга
		kernel[i] *= 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(size-1))
	}

	// Нормализация по амплитуде
	floats.Scale(1.0/floats.Sum(kernel), kernel)
	return kernel
}

// ThresholdFilter обнуляет значения сигнала ниже порога Threshold.
// Простейший пороговый фильтр (шумоподавление).
//
// Математически:
//
//	y[n] = x[n], если x[n] > T; иначе y[n] = 0
func ThresholdFilter(signal []float64, threshold float64) []float64 {
	result := make([]float64, len(signal))
	for i, v := range signal {
		if v > threshold {
			result[i] = v
		} else {
			result[i] = 0.0
		}
	}
	return result
}
