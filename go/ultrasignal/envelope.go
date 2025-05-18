package ultrasignal

import (
	"gonum.org/v1/gonum/dsp/fourier"
	"math"
)

// ComputeEnvelope рассчитывает пик-огибающую с экспоненциальным сглаживанием
func ComputeEnvelope(data []float64, window int) []float64 {
	if len(data) == 0 || window <= 0 {
		return data
	}
	result := make([]float64, len(data))
	result[0] = math.Abs(data[0])
	alpha := 0.1
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

// ComputeEnvelopeHilbert рассчитывает огибающую через преобразование Гильберта
func ComputeEnvelopeHilbert(signal []float64) []float64 {
	analytic := ComputeAnalyticSignal(signal)
	envelope := make([]float64, len(analytic))
	for i, v := range analytic {
		envelope[i] = Abs(real(v), imag(v))
	}
	return envelope
}

// ComputeAnalyticSignal возвращает комплексный аналитический сигнал
func ComputeAnalyticSignal(signal []float64) []complex128 {
	n := len(signal)
	// Преобразуем сигнал в комплексный (мнимая часть 0)
	complexSignal := make([]complex128, n)
	for i := range signal {
		complexSignal[i] = complex(signal[i], 0)
	}

	fft := fourier.NewCmplxFFT(n)
	spectrum := fft.Coefficients(nil, complexSignal)

	// Подавим отрицательные частоты
	for i := n/2 + 1; i < n; i++ {
		spectrum[i] = 0
	}
	// Удвоим амплитуды положительных частот, кроме DC и Nyquist
	for i := 1; i < n/2; i++ {
		spectrum[i] *= 2
	}

	return fft.Sequence(nil, spectrum)
}
