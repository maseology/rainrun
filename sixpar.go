package rainrun

import "math"

// SIXPAR model
// ref: Gupta V.K., S. Sorooshian, 1983. Uniqueness and Observability of Conceptual Rainfall-Runoff Model Parameters: The Percolation Process Examined. Water Resources Research 19(1). pp.269-276.
// also see: Duan, Q., S. Sorooshian, V. Gupta, 1992. Effective and Efficient Global Optimization for Conceptual Rainfall-Runoff Models. Water Resources Research 28(4). pp.1015-1031.
type SIXPAR struct {
	up, low    res
	beta, z, x float64
}

// New SIXPAR constructor
// [upCap, lowCap, upK, lowK, z, x]
func (m *SIXPAR) New(p ...float64) {
	// for TWOPAR, set pLM=0, variables pLK, pZ, pX, will have no impact
	m.up.new(p[0], p[2])  // upper reservoir: fast subsurface flow (interflow)
	m.low.new(p[1], p[3]) // update lower reservoir: slow subsurface flow (baseflow)
	m.beta = p[1] * p[3]
	m.z = p[4]
	m.x = p[5]
}

// Update state for daily inputs
func (m *SIXPAR) Update(p, ep float64) (float64, float64, float64) {
	pn := p - ep // net precipitation less ET
	var lt float64
	if m.low.cap > 0. {
		lt = 1. - m.low.storageFraction()
	}
	ut := m.up.storageFraction()
	g := math.Min(m.beta*ut*(1.+m.z*math.Pow(lt, m.x)), m.up.sto) // percolation from upper reservoir to lower reservoir (PDt=0 if pLM=0)
	q := m.up.overflow(pn - g)                                    // add rainfall and remove percolation from upper reservoir, remainder becomes saturation excess overland runoff
	m.low.update(g)                                               // add percolation to lower reservoir
	// add baseflow to runoff and update reservoirs
	q += m.up.decayExp() + m.low.decayExp() // _uk * USt + _lk * LSt 'total discharge
	return ep, q, g
}

// Storage returns total storage
func (m *SIXPAR) Storage() float64 {
	return m.up.sto + m.low.sto
}

// SampleSpace returns a hypercube from which the optimum resides
func (m *SIXPAR) SampleSpace(u []float64) []float64 {
	// const sd, n, fc = 1000.0, 0.3, 0.1
	// x1 := mm.LinearTransform(0., sd*fc, u[1])     // threshold storage (sfc=D(fc-tr))
	// x0 := x1 + mm.LinearTransform(0., sd*n, u[0]) // watershed storage (sbc=D(n-tr))
	// x2 := mm.LinearTransform(0., 1., u[2])        // coverdense
	// x3 := mm.LinearTransform(0., 0.01, u[3])      // intcap
	// x4 := mm.LinearTransform(0., 1., u[4])        // kb
	// x5 := mm.LinearTransform(0., 100., u[5])      // a
	// x6 := mm.LinearTransform(0., 1., u[6])        // b
	// return []float64{x0, x1, x2, x3, x4, x5, x6}
	return []float64{-99999.0}
}
