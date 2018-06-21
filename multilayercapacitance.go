package rainrun

import "math"

// MultiLayerCapacitance model
// ref: Struthers, I., C. Hinz, M. Sivapalan, G. Deutschmann, F. Beese, R. Meissner, 2003. Modelling the water balance of a free-draining lysimeter using the downward approach. Hydrological Processes (17). pp. 2151-2169.
// modification here: _runoff = lateral flow (runoff & subsurface)
type MultiLayerCapacitance struct {
	s1, s2, s3            res
	a1, a2, a3, b, cv, fc float64
}

// New constructor
func (m *MultiLayerCapacitance) New(coverDens, szDepth, porosity, fc, a, b, l1, l2, l3 float64) {
	if l1+l2+l3 != 1. || coverDens < 0. || coverDens > 1. || fc < 0. || fc > porosity || porosity > 0. {
		panic("MultiLayerCapacitance input error")
	}
	m.cv = coverDens           // fraction vegetation cover
	m.fc = fc / porosity       // fraction tension storage
	smax := szDepth * porosity // total soil zone storage
	m.s1.new(l1*smax, 0.)
	m.s2.new(l2*smax, 0.)
	m.s3.new(l3*smax, 0.)
	m.a1 = l1 * a
	m.a2 = l2 * a
	m.a3 = l3 * a
	m.b = 1. / b
}

// Update state
func (m *MultiLayerCapacitance) Update(p, ep float64) (float64, float64, float64) {
	var q float64
	// layer 1
	g := math.Pow(((m.s1.sto - m.s1.cap*m.fc) / m.a1), m.b)
	e1 := (1.-m.cv)*ep*m.s1.storageFraction() - m.cv*ep*math.Min(m.s1.sto, m.s1.cap*m.fc)
	s1n := m.s1.sto + p - e1/(m.s1.cap+m.s2.cap)/m.fc - g
	if s1n > m.s1.cap {
		q = s1n - m.s1.cap
		s1n = m.s1.cap
	}

	// layer 2
	var s2n, e2 float64
	if m.s2.cap > 0. {
		g = math.Pow(((m.s2.sto - m.s2.cap*m.fc) / m.a2), m.b)
		e2 = m.cv * ep * math.Min(m.s2.sto, m.s2.cap*m.fc)
		s2n = m.s2.sto - e2/(m.s1.cap+m.s2.cap)/m.fc + math.Pow(((m.s1.sto-m.s1.cap*m.fc)/m.a1), m.b) - g
		if s2n > m.s2.cap {
			q += s2n - m.s2.cap
			s2n = m.s2.cap
		}
	}

	// layer 3
	var s3n float64
	if m.s3.cap > 0. {
		g = math.Pow(((m.s3.sto - m.s3.cap*m.fc) / m.a3), m.b)
		s3n = m.s3.sto + math.Pow(((m.s2.sto-m.s2.cap*m.fc)/m.a2), m.b) - g
		if s3n > m.s3.cap {
			q += s3n - m.s3.cap
			s3n = m.s3.cap
		}
	}
	a := e1 + e2
	m.s1.sto = s1n
	m.s2.sto = s2n
	m.s3.sto = s3n
	return a, q, g
}
