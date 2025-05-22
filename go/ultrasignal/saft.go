package ultrasignal

// SAFT реализует упрощённый Synthetic Aperture Focusing Technique для суммирования нескольких сигналов с учётом задержек.
// Каждый сигнал смещается во времени (индексе) в зависимости от положения приёмника/излучателя,
// чтобы компенсировать разницу в пути распространения волны с учетом скорости звука.
//
// Входные параметры:
// - signals: срез временных сигналов от разных приёмников (предполагается одинаковая частота дискретизации).
// - sampleRate: частота дискретизации сигнала [Гц].
// - velocity: скорость распространения волны в среде [м/с].
//
// Возвращает:
// - сигнал, полученный сложением сдвинутых исходных сигналов (фокусированное усиление).
//
// Примечание: В текущей реализации сдвиг упрощён до целочисленного сдвига по индексам,
// но для физически точного SAFT сдвиги нужно рассчитывать на основе геометрии и задержек времени с интерполяцией.
func SAFT(signals [][]float64, sampleRate, velocity float64) []float64 {
	if len(signals) == 0 {
		return nil
	}

	n := len(signals[0])
	for _, signal := range signals {
		if len(signal) < n {
			n = len(signal)
		}
	}

	summed := make([]float64, n)
	count := make([]int, n)

	for i, signal := range signals {
		// В реальных задачах задержка зависит от расстояния до объекта:
		// timeDelay = distance / velocity
		// indexShift = int(timeDelay * sampleRate)
		// Здесь для примера просто используем индекс i как сдвиг.

		shift := i // Заглушка: заменить реальным вычислением задержки

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
