package rainrun

import "math"

// HBV model
// Bergström, S., 1976. Development and application of a conceptual runoff model for Scandinavian catchments. SMHI RHO 7. Norrköping. 134 pp.
// Bergström, S., 1992. The HBV model - its structure and applications. SMHI RH No 4. Norrköping. 35 pp
type HBV struct {
	sq, qt                                                      []float64
	fc, lp, beta, sm, suz, slz, uzl, k0, k1, k2, perc, lakefrac float64
}

// New HBV constructor
// [fc, lp, beta, uzl, k0, k1, k2, ksat, maxbas, lakeCoverFrac]
func (m *HBV) New(p ...float64) {
	if fracCheck(p[1]) || fracCheck(p[4]) || fracCheck(p[5]) || fracCheck(p[6]) {
		panic("HBV input eror")
	}
	m.fc = p[0]                         // max basin moisture storage
	m.lp = p[1]                         // soil moisture parameter
	m.beta = p[2]                       // soil moisture parameter
	m.uzl = p[3]                        // upper zone fast flow limit
	m.k0, m.k1, m.k2 = p[4], p[5], p[6] // fast, slow, and baseflow recession coefficients
	m.perc = p[7]                       // upper-to-lower zone percolation, assuming percolation rate = Ksat
	m.lakefrac = p[9]                   // lake fraction

	m.qt = TriangularTF(p[8], 0.5, 0.)               // MAXBAS: triangular weighted transfer function
	m.sq = make([]float64, len(m.qt)+1, len(m.qt)+1) // delayed runoff
}

// Update state for daily inputs
func (m *HBV) Update(pn, ep float64) (float64, float64, float64) {
	var a float64
	if m.lakefrac > 0. {
		a = m.hBVlake(pn, ep)
		// ep -= a // assume PET does not change (by commenting-out this line)
	}
	m.hBVinfiltration(pn * (1. - m.lakefrac))
	a += m.hBVet(ep)
	q, g := m.hBVrunoff()
	return a, q, g
}

func (m *HBV) hBVlake(pn, ep float64) float64 {
	m.slz += pn * m.lakefrac // assumes lakes are connected to the lower reservoir
	epl := ep * m.lakefrac
	a := epl
	if epl > m.slz {
		a = m.slz
	}
	m.slz -= a
	return a
}
func (m *HBV) hBVinfiltration(p float64) {
	i := p * math.Pow(m.sm/m.fc, m.beta)
	if i > p {
		panic("HBV error, infiltration")
	}
	m.sm += p - i // soil zone moisture storage
	m.suz += i    // upper zone moisture storage
}
func (m *HBV) hBVet(ep float64) float64 {
	etr := math.Min(1., m.sm/m.lp/m.fc) * ep
	if etr >= m.sm {
		etr = m.sm
		m.sm = 0.
	} else {
		m.sm -= etr
	}
	return etr
}
func (m *HBV) hBVrunoff() (float64, float64) {
	// soil zone accounting
	q0 := math.Max(m.k0*(m.suz-m.uzl), 0.0) // fast runoff
	m.suz -= q0
	q1 := m.k1 * m.suz // slow runoff
	m.suz -= q1        // q0 + q1 'total runoff
	q2 := m.k2 * m.slz // baseflow
	m.slz -= q2        // lower zone moisture storage

	// stream flow response function
	rgen := q0 + q1 + q2 // generated runoff
	for i := 1; i <= len(m.qt); i++ {
		m.sq[i-1] = m.sq[i] + m.qt[i-1]*rgen
	}
	q := m.sq[0]

	// percolate to lower reservoir
	g := math.Min(m.perc, m.suz)
	m.suz -= g
	m.slz += g

	return q, g
}

// Storage returns total storage
func (m *HBV) Storage() float64 {
	return m.suz + m.slz
}
