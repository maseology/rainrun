package rainrun

import (
	"math"
)

// DawdyODonnell model
// ref: Dawdy, D.R., and T. O'Donnell, 1965. Mathematical Models of Catchment Behavior. Journal of Hydraulics Division, ASCE, Vol. 91, No. HY4, pp. 123-137.
//see:  pg.34 in Dooge and O'Kane (2003)
type DawdyODonnell struct {
	depint, upsz manabe
	ores, gwres  res // S, G
	gwcap, ksat  float64
}

// New DawdyODonnell constructor
// [ksat, depintCap, upszCap, gwCap, olfk, bfk]
func (m *DawdyODonnell) New(p ...float64) {
	if p[0] < 0. {
		panic("DawdyODonnell error, ksat < 0.0")
	}
	m.ksat = p[0]
	m.depint.new(p[1], 1., 0.)            // R; depintCap = R*
	m.ores.new(math.MaxFloat64, p[4])     // S; overland flow recession coefficient
	m.upsz.new(math.MaxFloat64, 1., p[2]) // M; upszCap = M*
	m.gwres.new(p[3], p[5])               // G; gwCap = G*; baseflow recession coefficient
}

// Update state for daily inputs
func (m *DawdyODonnell) Update(p, ep float64) (float64, float64, float64) {
	// fill depressions & interception (R)
	eR, q1, f := m.depint.update(p, ep, m.ksat) // set percolation rate (F) to vertical conductivity, and overflow (Q1) to S
	m.ores.update(q1)                           // to overland flow stor (S)
	// upper soil zone (M)
	_, _, d := m.upsz.update(f, 0.0, m.ksat) // add percolation; set recharge rate to vertical conductivity
	// lower soil zone (G)
	c := m.gwres.update(d)                   // (C)
	eM, _, _ := m.upsz.update(c, ep-eR, 0.0) // add lower overflow back to upper soil zone (C)
	// total flow, AET, recharge
	q := m.gwres.decayExp() + m.ores.decayExp() // Qt = Qb + Qs
	a := eM + eR                                // total ET = EM + ER
	g := c - d                                  // net recharge
	return a, q, g
}

// Storage returns total storage
func (m *DawdyODonnell) Storage() float64 {
	return m.depint.sto + m.upsz.sto + m.ores.sto + m.gwres.sto
}

// SampleSpace returns a hypercube from which the optimum resides
func (m *DawdyODonnell) SampleSpace(u []float64) []float64 {
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
