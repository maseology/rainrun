package sample

import (
	mm "github.com/maseology/mmaths"
	"github.com/maseology/montecarlo/jointdist"
)

const (
	soildepth = 1.
	n         = 0.3
	fc        = 0.1
)

// Makkink (4)
func Makkink(u []float64) []float64 {
	// a := 1. // a and alpha are effectively the same
	b := mm.LinearTransform(0., .1, u[0])
	c := mm.LinearTransform(0., 5., u[1])
	alpha := mm.LinearTransform(0., 5., u[2])
	beta := mm.LinearTransform(-.0005, 0.0005, u[3])
	return []float64{b, c, alpha, beta}
}

// CCF (4)
func CCF(u []float64) []float64 {
	tindex := mm.LogLinearTransform(0.0002, 0.05, u[0]) // CCF temperature index; range .0002 to 0.0005 m/°C/d -- roughly 1/10 DDF (pg.278)
	ddfc := mm.LinearTransform(0., 10., u[1])           // DDF adjustment factor based on pack density, see DeWalle and Rango, pg. 275; Ref: Martinec (1960)=1.1
	baseT := mm.LinearTransform(-5., 5., u[2])          // base/critical temperature (°C)
	tsf := mm.LinearTransform(0.1, 0.7, u[3])           // TSF (surface temperature factor), 0.1-0.5 have been used
	// ddf := mm.LinearTransform(0.001, 0.008, u[1])    // (initial) degree-day/melt factor; range .001 to .008 m/°C/d  (pg.275)
	return []float64{tindex, ddfc, baseT, tsf}
}

////////////////
//// MODELS ////
////////////////

// Atkinson (7)
func Atkinson(u []float64) []float64 {
	x1 := mm.LinearTransform(0., soildepth*fc, u[1])     // threshold storage (sfc=D(fc-tr))
	x0 := x1 + mm.LinearTransform(0., soildepth*n, u[0]) // watershed storage (sbc=D(n-tr))
	x2 := mm.LinearTransform(0., 1., u[2])               // coverdense
	x3 := mm.LinearTransform(0., 0.01, u[3])             // intcap
	x4 := mm.LogLinearTransform(0.00001, 1., u[4])       // kb
	x5 := mm.LinearTransform(0., 100., u[5])             // a
	x6 := mm.LinearTransform(0., 1., u[6])               // b
	return []float64{x0, x1, x2, x3, x4, x5, x6}
}

// DawdyODonnell (6)
func DawdyODonnell(u []float64, ts float64) []float64 {
	ksat := mm.LogLinearTransform(1e-12, 1., u[0]) * ts // ksat [m/ts]
	rs := mm.LinearTransform(0., 0.1, u[1])             // depression and interception capacity R*
	ms := mm.LinearTransform(0., 1000., u[2])           // upper soil zone capacity M*
	gs := mm.LinearTransform(0., 1000., u[3])           // lower soil zone capacity G*
	s := mm.LogLinearTransform(1e-5, 1., u[4])          // overland flow recession coefficient
	b := mm.LogLinearTransform(1e-5, 1., u[5])          // baseflow recession coefficient
	return []float64{ksat, rs, ms, gs, s, b}
}

// GR4J (4)
func GR4J(u []float64) []float64 {
	prd := mm.LinearTransform(0., 1., u[0]) // x1: "production storage" capacity (mm)
	x2 := mm.LinearTransform(-.1, .1, u[1]) // x2: water exchange coefficient (>0 for water imports, <0 for exports, =0 for no exchange)
	rte := mm.LinearTransform(0., 3., u[2]) // x3: "routing storage"/groundwater storage capacity (mm)
	x4 := mm.LinearTransform(.5, 10., u[3]) // x4: unit hydrograph time base (days)
	// qsplt := mm.LinearTransform(0., 1., u[4]) // fixed in paper as 0.9
	return []float64{prd, x2, rte, x4} //, qsplt}
}

// CCFGR4J (8)
func CCFGR4J(u []float64) []float64 {
	ugr4j := GR4J(u)
	uccf := CCF(u[4:])
	return append(ugr4j, uccf...)
}

// MakkinkCCFGR4J (12)
func MakkinkCCFGR4J(u []float64) []float64 {
	ugr4j := GR4J(u)
	uccf := CCF(u[4:8])
	mak := Makkink(u[8:])
	return append(ugr4j, append(uccf, mak...)...)
}

// HBV (9)
func HBV(u []float64, ts float64) []float64 {
	fc := mm.LinearTransform(0., 1., u[0])
	lp := mm.LinearTransform(0., 1., u[1])
	beta := mm.LinearTransform(0., 10., u[2])
	uzl := mm.LinearTransform(0., 100., u[3]) // upper zone fast flow limit
	k0 := mm.LinearTransform(0., 1., u[4])
	k1 := mm.LinearTransform(0., 1., u[5])
	k2 := mm.LinearTransform(0., 1., u[6])
	perc := mm.LogLinearTransform(1e-12, 1., u[7]) * ts // ksat [m/d]
	maxbas := mm.LinearTransform(0., 10., u[8])         // days
	// lakefrac := mm.LinearTransform(0., 1., u[9])
	return []float64{fc, lp, beta, uzl, k0, k1, k2, perc, maxbas} //, lakefrac}
}

// CCFHBV (13)
func CCFHBV(u []float64, ts float64) []float64 {
	uhbv := HBV(u, ts)
	tindex := mm.LogLinearTransform(0.0002, 0.05, u[9]) // CCF temperature index; range .0002 to 0.0005 m/°C/d -- roughly 1/10 DDF (pg.278)
	ddfc := mm.LinearTransform(.01, 2.5, u[10])         // DDF adjustment factor based on pack density, see DeWalle and Rango, pg. 275; Ref: Martinec (1960)=1.1
	baseT := mm.LinearTransform(-5., 5., u[11])         // base/critical temperature (°C)
	tsf := mm.LinearTransform(0.1, 0.6, u[12])          // TSF (surface temperature factor), 0.1-0.5 have been used
	// ddf := mm.LinearTransform(0.001, 0.008, u[1])           // (initial) degree-day/melt factor; range .001 to .008 m/°C/d  (pg.275)
	uccf := []float64{tindex, ddfc, baseT, tsf}
	return append(uhbv, uccf...)
}

// ManabeGW (5)
func ManabeGW(u []float64) []float64 {
	u2t, u0t := jointdist.Nested2(u[2], u[0])
	x0 := mm.LinearTransform(0., soildepth, u0t)
	x1 := u[1]
	x2 := mm.LinearTransform(0., soildepth, u2t)
	x3 := mm.LogLinearTransform(1e-10, 1., u[3])
	x4 := u[4]
	return []float64{x0, x1, x2, x3, x4}
}

// MultiLayerCapacitance (9)
func MultiLayerCapacitance(u []float64) []float64 {
	cv := mm.LinearTransform(0., 1., u[0])
	x1 := mm.LinearTransform(0., soildepth*1000., u[1])
	uj0, uj1 := jointdist.Nested2(u[2], u[3])
	x2 := mm.LinearTransform(0., n, uj0)
	fc := mm.LinearTransform(0., n, uj1)
	a := mm.LinearTransform(0., 100., u[4])
	b := mm.LinearTransform(0., 1., u[5])
	l := jointdist.SumToOne(u[6], u[7], u[8])
	return []float64{cv, x1, x2, fc, a, b, l[0], l[1], l[2]}
}

// Quinn (11)
func Quinn(u []float64) []float64 {
	intCap := mm.LinearTransform(0., 0.1, u[0])
	impCap := mm.LinearTransform(0., 0.1, u[1])
	gwCap := mm.LinearTransform(0., 100., u[2])
	fImp := mm.LinearTransform(0., 1., u[3])
	ksat := mm.LogLinearTransform(1e-12, 1., u[4]) // ksat [m/s]
	rootZoneDepth := mm.LinearTransform(0., soildepth*1000., u[5])
	porosity := mm.LinearTransform(0., n, u[6])
	fieldCap := mm.LinearTransform(0., fc, u[7])
	f := mm.LinearTransform(0., 1., u[8])
	alpha := mm.LinearTransform(0., 1., u[9])
	zwt := mm.LinearTransform(0., 10., u[10])
	return []float64{intCap, impCap, gwCap, fImp, ksat, rootZoneDepth, porosity, fieldCap, f, alpha, zwt}
}

// SIXPAR (6)
func SIXPAR(u []float64) []float64 {
	upCap := mm.LinearTransform(0., 100., u[0])
	lowCap := mm.LinearTransform(0., 100., u[1])
	upK := mm.LinearTransform(0., 1., u[2])
	lowK := mm.LinearTransform(0., 1., u[3])
	z := mm.LinearTransform(0., 1., u[4])
	x := mm.LinearTransform(0., 1., u[5])
	return []float64{upCap, lowCap, upK, lowK, z, x}
}

// SPLR (5)
func SPLR(u []float64) []float64 {
	r12 := mm.LinearTransform(0., 10., u[0])
	r23 := mm.LinearTransform(0., 1., u[1])
	k1 := mm.LinearTransform(0., 1., u[2])
	k2 := mm.LinearTransform(0., 1., u[3])
	k3 := mm.LinearTransform(0., 1., u[4])
	return []float64{r12, r23, k1, k2, k3}
}
