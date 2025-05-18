package ultrasignal

func CrossCorrelate(x, y []float64) []float64 {
	n := len(x)
	m := len(y)
	corr := make([]float64, n+m-1)
	for lag := -(m - 1); lag < n; lag++ {
		sum := 0.0
		for i := 0; i < m; i++ {
			j := i + lag
			if j >= 0 && j < n {
				sum += y[i] * x[j]
			}
		}
		corr[lag+m-1] = sum
	}
	return corr
}
