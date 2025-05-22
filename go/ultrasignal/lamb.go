package ultrasignal

import "math"

// PhaseVelocity вычисляет приближенную фазовую скорость (vₚ) в зависимости от частоты и режима.
//
// Модель основана на эмпирических зависимостях:
//   - A0: дисперсионная зависимость (медленная волна, экспоненциальное приближение снизу)
//   - S0: слабодисперсионная (приближается сверху)
//
// Формула (например, для A0):
//
//	vₚ(f) ≈ v₀ * (1 + a * exp(-f / f₀))
//
// Параметры:
//   - freq: частота [Гц]
//   - thickness: толщина пластины [м] (не используется в текущей модели, но может использоваться позже)
//   - mode: "A0" или "S0"
//
// Возвращает:
//   - фазовая скорость [м/с]
func PhaseVelocity(freq, thickness float64, mode string) float64 {
	switch mode {
	case "A0":
		return 3100 * (1 + 0.5*math.Exp(-freq/1e5)) // Дисперсия при низких частотах
	case "S0":
		return 5900 * (1 - 0.3*math.Exp(-freq/1e5)) // Быстрое насыщение сверху
	default:
		return 0
	}
}

// GroupVelocity приближённо рассчитывает групповую скорость (v_g = dω/dk)
// через численную производную по частоте с центральной разностью.
//
// Параметры:
//   - freq: текущая частота [Гц]
//   - thickness: толщина пластины [м]
//   - mode: режим волны ("A0" или "S0")
//
// Возвращает:
//   - групповая скорость [м/с]
func GroupVelocity(freq, thickness float64, mode string) float64 {
	// Центрированное численное приближение производной dω/dk
	df := 1e3 // шаг частоты для численной производной [Гц]
	v1 := PhaseVelocity(freq-df, thickness, mode)
	v2 := PhaseVelocity(freq+df, thickness, mode)

	// Волновое число: k = ω / v => dk ≈ d(ω/v) = (2πf) / v
	// => dω/dk = v_g ≈ (f2 - f1) / (k2 - k1)
	omega1 := 2 * math.Pi * (freq - df)
	omega2 := 2 * math.Pi * (freq + df)

	k1 := omega1 / v1
	k2 := omega2 / v2

	if k2 == k1 {
		return 0
	}
	return (omega2 - omega1) / (k2 - k1)
}
