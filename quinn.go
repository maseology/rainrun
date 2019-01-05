package rainrun

import "math"

// Quinn simple storage model
// ref: Quinn P.F., K.J. Beven, 1993. Spatial and temporal predictions of soil moisture dynamics, runoff, variable source areas and evapotranspiration for Plynlimon, mid-Wales. Hydrological Processes 7. pp.425-448.
// used in early formulations of TOPMODEL, neglecting capillary fringe
type Quinn struct {
	intc, imp, sz, grav                  manabe
	fimp, n, fc, zr, ksat, f, alpha, Zwt float64
}

// New Quinn constructor
// [intercepCap, impStoCap, gwCap, fImp, ksat, rootZoneDepth, porosity, fieldCap, f, alpha, zwt]
func (m *Quinn) New(p ...float64) {
	if p[3] < 0. || p[3] > 1. || p[7] > p[6] || p[4] < 0. {
		panic("Quinn model input error")
	}
	m.intc.new(p[0], 1., 0.)
	m.imp.new(p[1], 1., 0.)
	m.fimp = p[3]
	m.zr = p[5]
	m.ksat = p[4]
	m.n = p[6]
	m.fc = p[7]
	m.sz.new(p[5]*(p[6]-p[7]), 1.-p[3], 0.)
	m.grav.new(p[2], 1.-p[3], 0.)
	m.f = p[8]
	m.alpha = p[9]
	m.Zwt = p[10] // setting as long-term average depth to watertable
}

// Update state for daily inputs
func (m *Quinn) Update(p, ep float64) (float64, float64, float64) {
	var q float64
	pn, ae := p, ep
	// interception
	if m.intc.cap > 0. {
		a1, p1, _ := m.intc.update(pn, ep, 0.0)
		pn = p1
		ae -= a1
	}
	// impervious area
	a2, q2, _ := m.imp.update(pn, ae, 0.0)
	q += q2 * m.fimp
	ae -= a2 * m.fimp
	// pervious area (root zone and gravity reservoir); no Hortonian mechanism
	etsz := ae * (1. - m.sz.storageFraction()) // soilzone manabe already accounts for impervious coverage
	a3, q3, g3 := m.sz.update(pn*(1.-m.fimp), etsz, m.ksat)
	ae -= a3
	_, q4, _ := m.grav.update(q3+g3, 0.0, 0.0) // excess moved to gravity storage
	q += q4 * (1. - m.fimp)                    // add saturation excess runoff

	gx := m.Zwt * (m.n - m.fc)
	if gx-m.grav.sto < m.grav.cap { // ET from gravity reservoir when nearly saturated
		if ae <= m.grav.cap-gx-m.grav.sto {
			a5, _, _ := m.grav.update(0.0, ae, 0.0)
			ae -= a5
		} else {
			a5, _, _ := m.grav.update(0.0, m.grav.cap-gx+m.grav.sto, 0.0)
			ae -= a5
		}
	}

	// totals
	_, _, g := m.grav.update(0.0, 0.0, math.Min(m.grav.sto, m.alpha*m.ksat*math.Exp(-m.f*m.Zwt))) // recharge [L/TS]; setting pf = 0 and alpha sets qv = kv
	a := ep - ae                                                                                  // returns AET
	return a, q, g
}

// Storage returns total storage
func (m *Quinn) Storage() float64 {
	return m.intc.sto + m.imp.sto + m.sz.sto + m.grav.sto
}

// SampleSpace returns a hypercube from which the optimum resides
func (m *Quinn) SampleSpace(u []float64) []float64 {
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
