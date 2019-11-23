package rainrun

import (
	"math"

	"github.com/maseology/goHydro/pet"
	"github.com/maseology/goHydro/snowpack"
	"github.com/maseology/goHydro/solirrad"
	"github.com/maseology/goHydro/transfunc"
)

const (
	// Bristow-Campbell-Makkink
	a     = 1.
	b     = 0.060639679562861176
	t     = 0.
	g     = .8972864886528819
	alpha = 1.3265625764694242
	beta  = -.0003664953523919842
)

// CCFHBV model
// Bergström, S., 1976. Development and application of a conceptual runoff model for Scandinavian catchments. SMHI RHO 7. Norrköping. 134 pp.
// Bergström, S., 1992. The HBV model - its structure and applications. SMHI RH No 4. Norrköping. 35 pp
type CCFHBV struct {
	HBV
	SP snowpack.CCF
	SI *solirrad.SolIrad
}

// New CCFHBV constructor
// [fc, lp, beta, uzl, k0, k1, k2, ksat, maxbas, lakeCoverFrac, tindex, ddfc, baseT, tsf]
func (m *CCFHBV) New(p ...float64) {
	if fracCheck(p[1]) || fracCheck(p[4]) || fracCheck(p[5]) || fracCheck(p[6]) { // || fracCheck(p[9]) {
		panic("HBV input eror")
	}
	m.fc = p[0]                         // max basin moisture storage
	m.lp = p[1]                         // soil moisture parameter
	m.beta = p[2]                       // soil moisture parameter
	m.uzl = p[3]                        // upper zone fast flow limit
	m.k0, m.k1, m.k2 = p[4], p[5], p[6] // fast, slow, and baseflow recession coefficients
	m.perc = p[7]                       // upper-to-lower zone percolation, assuming percolation rate = Ksat
	m.lakefrac = 0.                     //p[9]                   // lake fraction

	m.tf = transfunc.NewTF(p[8], 0.5, 0.) // MAXBAS: triangular weighted transfer function

	// Cold-content snow melt funciton
	tindex, ddfc, baseT, tsf := p[9], p[10], p[11], p[12]
	m.SP = snowpack.NewCCF(tindex, 0.0045, ddfc, baseT, tsf)
}

// Update state
func (m *CCFHBV) Update(v []float64, doy int) (y, a, r, g float64) {
	tx, tn, r, s := v[0], v[1], v[2], v[3]

	// calculate yield
	tm := (tx + tn) / 2.
	y = m.SP.Update(r, s, tm)

	// calculate ep
	ep := func() float64 {
		tm := (tx + tn) / 2.
		Kg := etRadToGlobal(m.SI.PSIdaily(doy), tx, tn)
		return pet.Makkink(Kg, tm, 101300., alpha, beta)
	}()

	a, r, g = m.HBV.Update(y, ep)
	// a = ep
	// r = y
	return
}

func etRadToGlobal(Ke, tx, tn float64) float64 {
	// see pg 151 in DeWalle & Rango; attributed to Bristow and Campbell (1984)
	// ref: Bristow, K.L. and G.S. Campbell, 1984. On the relationship between incoming solar radiation and daily maximum and minimum temperature. Agricultural and Forest Meteorology 31(2):159--166.
	return Ke * a * (1. - math.Exp(-b*math.Pow(tx-tn, g)))
}
