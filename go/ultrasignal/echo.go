package ultrasignal

import "math"

func DetectEchoes(signal []float64, threshold float64) []int {
	var indices []int
	for i, val := range signal {
		if math.Abs(val) > threshold {
			indices = append(indices, i)
		}
	}
	return indices
}
