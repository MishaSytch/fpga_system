package ultrasignal

import "math"

func PhaseVelocity(freq, thickness float64, mode string) float64 {
	switch mode {
	case "A0":
		return 3100 * (1 + 0.5*math.Exp(-freq/100e3))
	case "S0":
		return 5900 * (1 - 0.3*math.Exp(-freq/100e3))
	default:
		return 0
	}
}

func GroupVelocity(freq, thickness float64, mode string) float64 {
	v := PhaseVelocity(freq, thickness, mode)
	return v + freq*(v/100e3)
}
