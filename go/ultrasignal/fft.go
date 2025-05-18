package ultrasignal

import (
	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/interp"
	"math"
	"math/cmplx"
)

func ComputeFFT(input []float64, sampleRate float64) ([]float64, []float64) {
	// n := len(input)
	//	fft := fourier.NewFFT(n)
	//	coeffs := fft.Coefficients(nil, input)
	//
	//	// Вычисление спектра в линейной шкале
	//	linearFreqs := make([]float64, n/2)
	//	linearMag := make([]float64, n/2)
	//	step := sampleRate / float64(n)
	//
	//	for i := 0; i < n/2; i++ {
	//		linearFreqs[i] = float64(i) * step
	//		linearMag[i] = cmplx.Abs(coeffs[i]) / float64(n)
	//	}

	n := len(input)
	fft := fourier.NewFFT(n)
	coeffs := fft.Coefficients(nil, input)

	frequencies := make([]float64, n/2)
	magnitudes := make([]float64, n/2)
	step := sampleRate / float64(n)

	for i := 0; i < n/2; i++ {
		frequencies[i] = float64(i) * step
		magnitudes[i] = cmplx.Abs(coeffs[i]) / float64(n)
	}
	return frequencies, magnitudes
}

func ComputeFFTLog(input []float64, sampleRate float64, logMin float64, logMax float64, points int) ([]float64, []float64) {
	linearFreqs, linearMag := ComputeFFT(input, sampleRate)

	// Построение логарифмической шкалы
	logFreqs := make([]float64, points)
	logMag := make([]float64, points)
	logStep := (math.Log10(logMax) - math.Log10(logMin)) / float64(points-1)

	for i := 0; i < points; i++ {
		logFreqs[i] = math.Pow(10, math.Log10(logMin)+logStep*float64(i))
	}

	// Интерполяция спектра в логарифмические точки
	var interpSpline interp.PiecewiseLinear
	_ = interpSpline.Fit(linearFreqs, linearMag)

	for i := range logFreqs {
		logMag[i] = interpSpline.Predict(logFreqs[i])
	}

	return logFreqs, logMag
}
