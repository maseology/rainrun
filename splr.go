package rainrun

// SPLR : Simple Parallel Linear Reservoir
// Buytaert, W., and K. Beven, 2011. Models as multiple working hypotheses hydrological simulation of tropical alpine . Hydrological Processes 25. pp. 1784–1799.
// 3-reservoir Tank model
type SPLR struct {
	s1, s2, s3           float64
	r12, r23, k1, k2, k3 float64
}

// New : constructor
func (m *SPLR) New(r12, r23, k1, k2, k3 float64) {
	m.r12 = r12
	m.r23 = r23
	m.k1 = k1
	m.k2 = k2
	m.k3 = k3
}

// Update : update state, returns excess
func (m *SPLR) Update(p, ep float64) (float64, float64) {
	pn, sv := p-ep, m.ts()
	u(&m.s1, m.r12*pn)
	u(&m.s2, (1.-m.r12)*m.r23*pn)
	u(&m.s3, (1.-m.r12)*(1.-m.r23)*pn)
	aet := sv - m.ts() + pn
	return aet, q(&m.s1, m.k1) + q(&m.s2, m.k2) + q(&m.s3, m.k3)
}

func (m *SPLR) ts() float64 {
	return m.s1 + m.s2 + m.s3
}

func u(s *float64, p float64) {
	*s += p
	if *s < 0. {
		*s = 0.
	}
}

func q(s *float64, k float64) float64 {
	d := k * *s
	*s -= d
	return d
}
