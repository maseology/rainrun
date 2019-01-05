package rainrun

// Lumper : interface to rainfall-runfall lumped models
type Lumper interface {
	New(p ...float64)
	Update(p, ep float64) (float64, float64, float64)
	Storage() float64
	// SampleSpace([]float64) []float64
}
