package rainrun

import "math"

// Quinn simple storage model
// ref: Quinn P.F., K.J. Beven, 1993. Spatial and temporal predictions of soil moisture dynamics, runoff, variable source areas and evapotranspiration for Plynlimon, mid-Wales. Hydrological Processes 7. pp.425-448.
// used in early formulations of TOPMODEL, neglecting capillary fringe
type Quinn struct {
	intc, imp, sz, grav             Manabe
	fimp, n, fc, zr, ksat, f, alpha float64
}

// New constructor
func (m *Quinn) New(intercepCap, impStoCap, gwCap, fImp, ksat, rootZoneDepth, porosity, fieldCap, f, alpha float64) {
	if fImp < 0. || fImp > 1. || fieldCap > porosity || ksat < 0. {
		panic("Quinn model input error")
	}
	m.intc.New(intercepCap, 1., 0.)
	m.imp.New(impStoCap, 1., 0.)
	m.fimp = fImp
	m.zr = rootZoneDepth
	m.ksat = ksat
	m.n = porosity
	m.fc = fieldCap
	m.sz.New(rootZoneDepth*(porosity-fieldCap), 1.-fImp, 0.)
	m.grav.New(gwCap, 1.-fImp, 0.)
	m.f = f
	m.alpha = alpha
}

// Update state
func (m *Quinn) Update(p, ep, zwt float64) (float64, float64, float64) {
	var q float64
	pn, ae := p, ep
	// interception
	if m.intc.cap > 0. {
		a1, p1, _ := m.intc.Update(pn, ep, 0.0)
		pn = p1
		ae -= a1
	}
	// impervious area
	a2, q2, _ := m.imp.Update(pn, ae, 0.0)
	q += q2 * m.fimp
	ae -= a2 * m.fimp
	// pervious area (root zone and gravity reservoir); no Hortonian mechanism
	etsz := ae * (1. - m.sz.StorageFraction()) // soilzone manabe already accounts for impervious coverage
	a3, q3, g3 := m.sz.Update(pn*(1.-m.fimp), etsz, m.ksat)
	ae -= a3
	_, q4, _ := m.grav.Update(q3+g3, 0.0, 0.0) // excess moved to gravity storage
	q += q4 * (1. - m.fimp)                    // add saturation excess runoff

	gx := zwt * (m.n - m.fc)
	if gx-m.grav.sto < m.grav.cap { // ET from gravity reservoir when nearly saturated
		if ae <= m.grav.cap-gx-m.grav.sto {
			a5, _, _ := m.grav.Update(0.0, ae, 0.0)
			ae -= a5
		} else {
			a5, _, _ := m.grav.Update(0.0, m.grav.cap-gx+m.grav.sto, 0.0)
			ae -= a5
		}
	}

	// totals
	_, _, g := m.grav.Update(0.0, 0.0, math.Min(m.grav.sto, m.alpha*m.ksat*math.Exp(-m.f*zwt))) // recharge [L/TS]; setting pf = 0 and alpha sets qv = kv
	a := ep - ae                                                                                //returns AET
	return a, q, g
}
