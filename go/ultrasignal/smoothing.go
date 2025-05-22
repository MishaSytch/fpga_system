package ultrasignal

// MovingAverage вычисляет скользящее среднее сигнала с заданным окном.
// Для каждого элемента i результат — среднее значение элементов входного среза
// в интервале [i - window/2, i + window/2], учитывая границы массива.
//
// Формула:
//
//	y[i] = (1 / N_i) * Σ_{j=max(0, i - half)}^{min(len(input)-1, i + half)} x[j]
//
// где N_i — количество элементов в окне для позиции i.
//
// window — размер окна сглаживания (целое число, обычно нечётное).
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

// ExponentialSmoothing выполняет экспоненциальное сглаживание сигнала с коэффициентом alpha.
// Каждый выходной элемент — взвешенное среднее текущего входного и предыдущего сглаженного значения.
//
// Формула:
//
//	y[0] = x[0]
//	y[i] = α * x[i] + (1 - α) * y[i-1],   0 < α < 1
//
// α — коэффициент сглаживания, определяет степень влияния новых значений.
// При α близком к 1 сглаживание слабое, при α близком к 0 — сильное.
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
