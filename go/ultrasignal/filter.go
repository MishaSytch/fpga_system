package ultrasignal

import (
	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/floats"
	"math"
	"math/cmplx"
)

func BandPassFilter(input, kernel []float64) []float64 {
	return Convolve(input, kernel)
}

func ComputeAFC(kernel []float64, sampleRate float64) ([]float64, []float64) {
	n := len(kernel)
	delta := make([]float64, n)
	delta[0] = 1.0
	response := Convolve(delta, kernel)

	fft := fourier.NewFFT(n)
	coeffs := fft.Coefficients(nil, response)

	frequencies := make([]float64, n/2)
	afc := make([]float64, n/2)
	step := sampleRate / float64(n)

	for i := 0; i < n/2; i++ {
		frequencies[i] = float64(i) * step
		afc[i] = cmplx.Abs(coeffs[i]) / float64(n)
	}
	return frequencies, afc
}

func FIRBandPassKernel(size int, lowCutoff, highCutoff, sampleRate float64) []float64 {
	kernel := make([]float64, size)
	mid := size / 2
	low := 2 * math.Pi * lowCutoff / sampleRate
	high := 2 * math.Pi * highCutoff / sampleRate

	for i := 0; i < size; i++ {
		n := float64(i - mid)
		if n == 0 {
			kernel[i] = (high - low) / math.Pi
		} else {
			kernel[i] = (math.Sin(high*n) - math.Sin(low*n)) / (math.Pi * n)
		}
		kernel[i] *= 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(size-1)) // Hamming
	}

	floats.Scale(1.0/floats.Sum(kernel), kernel)
	return kernel
}

func TrashHoldFilter(signal []float64, Threshold float64) []float64 {
	result := make([]float64, len(signal))
	for i, v := range signal {
		if v > Threshold {
			result[i] = v
		} else {
			result[i] = 0.0
		}
	}
	return result
}
