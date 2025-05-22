package ultrasignal

import (
	"gonum.org/v1/gonum/dsp/fourier"
	"math"
	"math/cmplx"
)

// ComputeEnvelope вычисляет огибающую сигнала методом скользящего максимума с экспоненциальным сглаживанием.
//
// Огибающая e[n] рассчитывается как:
//
//	e[n] = α * max(|x[n-window+1]..x[n]|) + (1 - α) * e[n-1]
//
// где α ∈ [0..1] — коэффициент сглаживания.
//
// Параметры:
//   - data: входной сигнал
//   - window: ширина окна для оценки локального максимума
//
// Возвращает:
//   - result: массив огибающей амплитуды
func ComputeEnvelope(data []float64, window int) []float64 {
	if len(data) == 0 || window <= 0 {
		return data
	}
	result := make([]float64, len(data))
	result[0] = math.Abs(data[0])
	alpha := 0.1 // коэффициент сглаживания

	for i := 1; i < len(data); i++ {
		start := max(0, i-window)
		maxVal := 0.0
		for j := start; j <= i; j++ {
			val := math.Abs(data[j])
			if val > maxVal {
				maxVal = val
			}
		}
		result[i] = alpha*maxVal + (1-alpha)*result[i-1]
	}
	return result
}

// ComputeEnvelopeHilbert вычисляет точную огибающую амплитуды сигнала с использованием
// аналитического сигнала, построенного через преобразование Гильберта.
//
// Формула огибающей:
//
//	envelope[n] = |x[n] + j * H{x[n]}| = sqrt(x[n]^2 + H{x[n]}^2)
//
// Где H{x[n]} — преобразование Гильберта.
//
// Параметры:
//   - signal: входной вещественный сигнал
//
// Возвращает:
//   - envelope: массив огибающих значений
func ComputeEnvelopeHilbert(signal []float64) []float64 {
	analytic := ComputeAnalyticSignal(signal)
	envelope := make([]float64, len(analytic))
	for i, v := range analytic {
		envelope[i] = cmplx.Abs(v)
	}
	return envelope
}

// ComputeAnalyticSignal строит аналитический сигнал (complex-valued), применяя преобразование Гильберта.
//
// Алгоритм:
//  1. Преобразуем сигнал в спектр (FFT).
//  2. Убираем отрицательные частоты (анализируем только положительные).
//  3. Удваиваем положительные частоты (кроме DC и Nyquist).
//  4. Обратным FFT получаем комплексный сигнал x[n] + j*H{x[n]}.
//
// Параметры:
//   - signal: вещественный временной сигнал
//
// Возвращает:
//   - complexSignal: комплексный аналитический сигнал
func ComputeAnalyticSignal(signal []float64) []complex128 {
	n := len(signal)

	// Преобразуем в комплексный сигнал (мнимая часть = 0)
	complexSignal := make([]complex128, n)
	for i := range signal {
		complexSignal[i] = complex(signal[i], 0)
	}

	// Быстрое преобразование Фурье
	fft := fourier.NewCmplxFFT(n)
	spectrum := fft.Coefficients(nil, complexSignal)

	// Удаление отрицательных частот
	for i := n/2 + 1; i < n; i++ {
		spectrum[i] = 0
	}

	// Удваиваем положительные частоты (кроме DC и Nyquist)
	for i := 1; i < n/2; i++ {
		spectrum[i] *= 2
	}

	// Обратное преобразование Фурье — получаем аналитический сигнал
	return fft.Sequence(nil, spectrum)
}
