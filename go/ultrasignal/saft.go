package ultrasignal

func SAFT(signals [][]float64, sampleRate, velocity float64) []float64 {
	if len(signals) == 0 {
		return nil
	}

	// Найдём минимальную длину сигнала, чтобы не выйти за границы
	n := len(signals[0])
	for _, signal := range signals {
		if len(signal) < n {
			n = len(signal)
		}
	}

	summed := make([]float64, n)
	count := make([]int, n)

	for i, signal := range signals {
		shift := i // Упрощённо: фиктивный сдвиг
		for t := 0; t < n-shift; t++ {
			summed[t] += signal[t+shift]
			count[t]++
		}
	}

	for i := range summed {
		if count[i] > 0 {
			summed[i] /= float64(count[i])
		}
	}
	return summed
}
