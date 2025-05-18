package ultrasignal

import "math"

// UltrasonicSignal представляет собой акустический сигнал и методы анализа
type UltrasonicSignal struct {
	Raw         []float64
	SampleRate  float64
	Envelope    []float64
	FFTMag      []float64
	Frequencies []float64
	EchoIndices []int
}

// ComputeEnvelope рассчитывает огибающую выбранным методом
func (s *UltrasonicSignal) ComputeEnvelope(method string, window int) {
	switch method {
	case "hilbert":
		s.Envelope = ComputeEnvelopeHilbert(s.Raw)
	case "peak":
		s.Envelope = ComputeEnvelope(s.Raw, window)
	case "smooth":
		s.Envelope = ExponentialSmoothing(ComputeEnvelope(s.Raw, window), 0.1)
	default:
		panic("unknown envelope method")
	}
}

// ComputeFFT рассчитывает спектр сигнала
func (s *UltrasonicSignal) ComputeFFT() {
	s.Frequencies, s.FFTMag = ComputeFFT(s.Raw, s.SampleRate)
}

// DetectEchoes определяет позиции эхо-сигналов по порогу
func (s *UltrasonicSignal) DetectEchoes(threshold float64) {
	s.EchoIndices = DetectEchoes(s.Envelope, threshold)
}

// GetTimeOfFlight возвращает время пролета до первого эха
func (s *UltrasonicSignal) GetTimeOfFlight() float64 {
	if len(s.EchoIndices) == 0 {
		return -1
	}
	return float64(s.EchoIndices[0]) / s.SampleRate
}

// GetTimeOfFlight рассчитывает время пролета (ToF) до первого обнаруженного эхо-сигнала
func GetTimeOfFlight(echoIndices []int, sampleRate float64) float64 {
	if len(echoIndices) == 0 {
		return -1 // Нет эха
	}
	return float64(echoIndices[0]) / sampleRate
}

// HammingWindow применяет окно Хэмминга к сигналу
func HammingWindow(signal []float64) []float64 {
	n := len(signal)
	windowed := make([]float64, n)

	for i := 0; i < n; i++ {
		hamming := 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(n-1))
		windowed[i] = signal[i] * hamming
	}

	return windowed
}
