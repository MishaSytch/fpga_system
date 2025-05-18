package ultrasignal

import (
	"math"
	"time"
)

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

func Abs(real, imag float64) float64 {
	return math.Sqrt(real*real + imag*imag)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func FreqToTime(frequency float64) time.Duration {
	var result time.Duration
	timeLoop := 1 / frequency
	if timeLoop < 1 {
		timeLoop *= time.Second.Seconds() / time.Millisecond.Seconds()
		if timeLoop < 1 {
			timeLoop *= time.Millisecond.Seconds() / time.Microsecond.Seconds()
			if timeLoop < 1 {
				timeLoop *= time.Microsecond.Seconds() / time.Nanosecond.Seconds()
				result = time.Duration(timeLoop) * time.Nanosecond
			} else {
				result = time.Duration(timeLoop) * time.Microsecond
			}
		} else {
			result = time.Duration(timeLoop) * time.Millisecond
		}
	} else {
		result = time.Duration(timeLoop) * time.Second
	}

	return result
}
