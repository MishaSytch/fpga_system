package ultrasignal

func MovingAverage(input []float64, window int) []float64 {
	if window <= 1 || len(input) == 0 {
		return input
	}
	output := make([]float64, len(input))
	half := window / 2

	for i := range input {
		sum := 0.0
		count := 0
		for j := max(0, i-half); j <= min(len(input)-1, i+half); j++ {
			sum += input[j]
			count++
		}
		output[i] = sum / float64(count)
	}
	return output
}

func ExponentialSmoothing(input []float64, alpha float64) []float64 {
	if len(input) == 0 || alpha <= 0 || alpha >= 1 {
		return input
	}
	output := make([]float64, len(input))
	output[0] = input[0]
	for i := 1; i < len(input); i++ {
		output[i] = alpha*input[i] + (1-alpha)*output[i-1]
	}
	return output
}
