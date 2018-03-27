package rainrun

import "math"

// GR4J model
// Perrin C., C. Michel, V. Andreassian, 2003. Improvement of a parsimonious model for streamflow simulation. Journal of Hydrology 279. pp. 275-289.
type GR4J struct {
	sto, gw       res
	uh1, uh2      []float64
	x2, qsplt, x4 float64
}

// New contructor
func (m *GR4J) New(stocap, gwstocap, x4, unitHydrographPartition, x2 float64) {
	// StorageCapacity (x1: maximum capacity of the SMA store)
	// GroundwaterStorageCapacity (x3: reference capacity of GW store)
	// x2: water exchange coefficient (>0 for water imports, <0 for exports, =0 for no exchange)
	// x4: unit hydrograph time parameter
	// unitHydrographPartition, default = 0.9
	if x4 < 0.5 || unitHydrographPartition <= 0. || unitHydrographPartition >= 1.0 {
		panic("GR4J input error")
	}
	m.sto.New(stocap, 0.)
	m.gw.New(gwstocap, 0.)
	m.x2 = x2
	m.qsplt = unitHydrographPartition // I interpret this as a runoff coefficient
	// unit hydrographs parameterization
	m.x4 = x4
	n1f, x4d1 := math.Modf(x4)
	n2f, x4d2 := math.Modf(2. * x4)
	n1, n2 := int(n1f), int(n2f)
	if x4d1 == 0. {
		n1-- // dimension of UH1(0 To n1)
	}
	if x4d2 == 0. {
		n2-- // dimension of UH2(0 To n2)
	}
	m.uh1 = make([]float64, n1, n1)
	m.uh2 = make([]float64, n2, n2)
}

// Update state
func (m *GR4J) Update(p, ep float64) (float64, float64, float64) {
	var pn, en, es float64
	if p >= ep {
		pn = p - ep
	} else {
		en = ep - p // available PET
	}
	x1 := m.sto.cap // x1: maximum capacity of the SMA store
	d1 := math.Tanh(pn / x1)
	sf := m.sto.StorageFraction()
	ps := x1 * d1 * (1. - math.Pow(sf, 2.)) / (1. + d1*sf) // Ps: portion of rain infiltrating soils (production) store
	if en > 0. {
		d1 = math.Tanh(en / x1)
		es = m.sto.sto * d1 * (2. - sf) / (1. + d1*(1.-sf)) // Es: soil evaporation
	}
	m.sto.Update(ps - es)
	g := m.sto.sto * (1. - math.Pow(1.+math.Pow(4.*m.sto.StorageFraction()/9., 4.), -0.25)) // percolation from production zone
	if m.sto.Update(-g) < 0. {                                                              // this line must be left here such that sto is updated
		panic("GR4J error: percolation")
	}
	pr := g + pn + ps
	q9 := m.updateUH1(m.qsplt * pr)
	q1 := m.updateUH2((1. - m.qsplt) * pr)

	// x3 := m.gw.cap                                       // x3: reference capacity of GW store
	fe := m.x2 * math.Pow(m.gw.StorageFraction(), 7./2.) // catchment GW exchange; x2: water exchange coefficient (>0 for water imports, <0 for exports, =0 for no exchange)
	m.gw.Update(q9 + fe)
	qr := m.gw.sto * (1. - math.Pow(1.+math.Pow(m.gw.StorageFraction(), 4.), -0.25))
	if m.gw.Update(-qr) < 0. { // this line must be left here such that gw is updated
		panic("GR4J error: percolation")
	}

	qd := math.Max(0., q1+fe)
	return en - es, qd + qr, q9
}

func (m *GR4J) updateUH1(pr float64) float64 {
	// unit hydrograph 1 for the GR4J model
	var sh, shsv float64
	n := len(m.uh1)
	for t := 0; t <= n; t++ {
		tf := float64(t)
		if tf < m.x4 {
			sh = math.Pow(tf/m.x4, 5./2.)
		} else {
			sh = 1.
		}
		if t == n-1 {
			m.uh1[t] = pr * (sh - shsv)
		} else {
			m.uh1[t] += pr * (sh - shsv)
		}
		shsv = sh
	}
	return m.uh1[0]
}
func (m *GR4J) updateUH2(pr float64) float64 {
	// unit hydrograph 2 for the GR4J model
	var sh, shsv float64
	n := len(m.uh2)
	for t := 0; t <= n; t++ {
		tf := float64(t)
		if tf < m.x4 {
			sh = 0.5 * math.Pow(tf/m.x4, 5./2.)
		} else if tf < 2.*m.x4 {
			sh = 1. - 0.5*math.Pow(2.-tf/m.x4, 5./2.)
		} else {
			sh = 1.
		}
		if t == n-1 {
			m.uh2[t] = pr * (sh - shsv)
		} else {
			m.uh2[t] += pr * (sh - shsv)
		}
		shsv = sh
	}
	return m.uh2[0]
}
