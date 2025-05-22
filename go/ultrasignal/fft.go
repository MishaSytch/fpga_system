package ultrasignal

import (
	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/interp"
	"math"
	"math/cmplx"
)

// ComputeFFT выполняет дискретное преобразование Фурье (БПФ) для вещественного сигнала,
// возвращая амплитудный спектр в линейной частотной шкале.
//
// Формула амплитудного спектра (для частоты fₖ):
//     A(fₖ) = 2 * |Xₖ| / N  , если k ≠ 0 и k ≠ N/2
//     A(fₖ) = |Xₖ| / N      , если k = 0 или k = N/2 (при чётном N)
//
// Где:
//     Xₖ — k-й коэффициент FFT
//     N  — длина сигнала (дискретных отсчётов)
//     A(fₖ) — амплитуда в Гц (в единицах входного сигнала)
//     fₖ = k * (Fs / N) — частота, Гц
//
// Fs — частота дискретизации (sampleRate)

func ComputeFFT(input []float64, sampleRate float64) ([]float64, []float64) {
	n := len(input)
	fft := fourier.NewFFT(n)
	coeffs := fft.Coefficients(nil, input)

	halfN := n/2 + 1
	frequencies := make([]float64, halfN)
	magnitudes := make([]float64, halfN)
	step := sampleRate / float64(n)

	for i := 0; i < halfN; i++ {
		frequencies[i] = float64(i) * step
		mag := cmplx.Abs(coeffs[i]) / float64(n)

		// Удваиваем амплитуду для всех частот кроме DC и Nyquist (если есть)
		if i != 0 && i != n/2 {
			mag *= 2
		}
		magnitudes[i] = mag
	}
	return frequencies, magnitudes
}

// ComputeFFTLog строит логарифмический спектр (по частоте) из временного сигнала.
// Используется линейный спектр, интерполированный в логарифмическую сетку.
//
// Шаг логарифмической частоты (равномерный в log10 масштабе):
//     fᵢ = 10 ^ [log10(f_min) + i * Δlog],   где i = 0..(points-1)
//     Δlog = (log10(f_max) - log10(f_min)) / (points - 1)
//
// Интерполяция выполняется по амплитудному спектру из ComputeFFT.

func ComputeFFTLog(input []float64, sampleRate float64, logMin float64, logMax float64, points int) ([]float64, []float64) {
	linearFreqs, linearMag := ComputeFFT(input, sampleRate)

	logFreqs := make([]float64, points)
	logMag := make([]float64, points)
	logStep := (math.Log10(logMax) - math.Log10(logMin)) / float64(points-1)

	for i := 0; i < points; i++ {
		logFreqs[i] = math.Pow(10, math.Log10(logMin)+logStep*float64(i))
	}

	var interpSpline interp.PiecewiseLinear
	_ = interpSpline.Fit(linearFreqs, linearMag)

	for i := range logFreqs {
		logMag[i] = interpSpline.Predict(logFreqs[i])
	}

	return logFreqs, logMag
}
