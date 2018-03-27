package rainrun

// Manabe reservoir
// standard form of a hydrological "bucket" model
// ref: Manabe, S., 1969. Climate and the Ocean Circulation 1: The Atmospheric Circulation and The Hydrology of the Earth's Surface. Monthly Weather Review 97(11). 739-744.
type Manabe struct {
	res
	expo, minsto float64
}

// New constructor
func (m *Manabe) New(capacity, fexposed, minSto float64) {
	if capacity < 0. || minSto < 0. || minSto > capacity || fexposed < 0.0 {
		panic("Manabe parameter error")
	}
	m.cap = capacity
	m.minsto = minSto
	m.expo = fexposed
}

// UpdateExposure : change area of reservoir exposed to evaporative forcings
func (m *Manabe) UpdateExposure(newExposure float64) {
	if newExposure < 0. {
		panic("UpdateExposure error")
	}
	if m.expo > 0. { // for example, in cases of seasonal LAI changes, storage capacity will also change
		fsto := newExposure / m.expo
		m.updateCapacity(fsto)
	}
	m.expo = newExposure
}

// updateCapacity : changes reservoir capacity
// if this causes and reservoir overflow (ie, ChangeFactor < 1), it will be determined with the next state update
func (m *Manabe) updateCapacity(changeFactor float64) {
	m.cap *= changeFactor
}

// Update state
func (m *Manabe) Update(p, ep, perc float64) (float64, float64, float64) {
	q := m.res.Update(p)
	if m.sto == 0. {
		return 0., q, 0.
	}
	if ep < mingtzero && perc < mingtzero {
		return 0., q, 0.
	}
	a, g := m.lossDirect(ep, perc)
	// a, g := m.lossExponential(ep, perc, 86400.)
	return a, q, g
}

func (m *Manabe) lossDirect(ep, perc float64) (float64, float64) {
	var a, g float64
	epx := m.expo * ep * m.res.StorageFraction() // effective PE
	if m.sto <= m.minsto {
		if ep == 0. {
			return 0., 0.
		}
		if epx >= m.sto {
			a = m.sto
			m.sto = 0.
		} else {
			a = epx
			m.sto -= epx
		}
	} else {
		if (epx + perc) > (m.sto - m.minsto) {
			fperc := perc / (epx + perc)
			g = fperc * (m.sto - m.minsto)
			a = m.sto - m.minsto - g
			epx -= a // reset to remaining available PE
			if epx >= m.minsto {
				a += m.minsto
				m.sto = 0.
			} else {
				a += epx
				m.sto = m.minsto - epx
			}
		} else {
			a = epx
			g = perc
			m.sto -= (a + g)
		}
	}
	return a, g
}

func (m *Manabe) lossExponential(ep, perc, ts float64) (float64, float64) {
	var a, g float64
	// first compute (direct) drainage
	sFree := m.sto - m.minsto
	if sFree > 0. {
		if sFree <= perc {
			m.sto = m.minsto
			g = sFree
		} else {
			m.sto -= perc
			g = perc
		}
	}
	// next compute AET
	a = m.res.DecayExp2(ep/ts, ts)
	return a, g
}
