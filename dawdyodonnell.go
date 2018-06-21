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

// New constructor
func (m *DawdyODonnell) New(ksat, depintCap, upszCap, gwCap, olfk, bfk float64) {
	if ksat < 0. {
		panic("DawdyODonnell error, ksat < 0.0")
	}
	m.ksat = ksat
	m.depint.new(depintCap, 1., 0.)          // R; depintCap = R*
	m.ores.new(math.MaxFloat64, olfk)        // S; overland flow recession coefficient
	m.upsz.new(math.MaxFloat64, 1., upszCap) // M; upszCap = M*
	m.gwres.new(gwCap, bfk)                  // G; gwCap = G*; baseflow recession coefficient
}

// Update state
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
