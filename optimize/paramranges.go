package optimize

import (
	mm "github.com/maseology/mmaths"
	"github.com/maseology/montecarlo/jointdist"
)

const (
	sd = 10.
	n  = 0.3
	fc = 0.1
)

// sample returns a hypercube from which the optimum resides

//////////////// Atkinson (7)
func sampleAtkinson(u []float64) []float64 {
	x1 := mm.LinearTransform(0., sd*fc, u[1])      // threshold storage (sfc=D(fc-tr))
	x0 := x1 + mm.LinearTransform(0., sd*n, u[0])  // watershed storage (sbc=D(n-tr))
	x2 := mm.LinearTransform(0., 1., u[2])         // coverdense
	x3 := mm.LinearTransform(0., 0.01, u[3])       // intcap
	x4 := mm.LogLinearTransform(0.00001, 1., u[4]) // kb
	x5 := mm.LinearTransform(0., 100., u[5])       // a
	x6 := mm.LinearTransform(0., 1., u[6])         // b
	return []float64{x0, x1, x2, x3, x4, x5, x6}
}

//////////////// DawdyODonnell (6)
func sampleDawdyODonnell(u []float64, ts float64) []float64 {
	ksat := mm.LogLinearTransform(1e-12, 1., u[0]) * ts // ksat [m/ts]
	rs := mm.LinearTransform(0., 0.1, u[1])             // depression and interception capacity R*
	ms := mm.LinearTransform(0., 1000., u[2])           // upper soil zone capacity M*
	gs := mm.LinearTransform(0., 1000., u[3])           // lower soil zone capacity G*
	s := mm.LogLinearTransform(1e-5, 1., u[4])          // overland flow recession coefficient
	b := mm.LogLinearTransform(1e-5, 1., u[5])          // baseflow recession coefficient
	return []float64{ksat, rs, ms, gs, s, b}
}

//////////////// GR4J (5)
func sampleGR4J(u []float64) []float64 {
	sto := mm.LinearTransform(0., 10., u[0])  // storage capacity
	gw := mm.LinearTransform(0., 100., u[1])  // groundwater storage capacity
	x2 := mm.LinearTransform(-10., 10., u[4]) // water exchange coefficient
	qsplt := mm.LinearTransform(0., 1., u[3])
	x4 := mm.LinearTransform(0.5, 1., u[2]) // water exchange coefficient
	return []float64{sto, gw, x4, qsplt, x2}
}

//////////////// HBV (9)
func sampleHBV(u []float64, ts float64) []float64 {
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

//////////////// ManabeGW (5)
func sampleManabeGW(u []float64) []float64 {
	u2t, u0t := jointdist.Nested2(u[2], u[0])
	x0 := mm.LinearTransform(0., sd, u0t)
	x1 := u[1]
	x2 := mm.LinearTransform(0., sd, u2t)
	x3 := mm.LogLinearTransform(1e-10, 1., u[3])
	x4 := u[4]
	return []float64{x0, x1, x2, x3, x4}
}

//////////////// MultiLayerCapacitance (9)
func sampleMultiLayerCapacitance(u []float64) []float64 {
	const sd, n = 1000.0, 0.3
	cv := mm.LinearTransform(0., 1., u[0])
	x1 := mm.LinearTransform(0., sd, u[1])
	uj0, uj1 := jointdist.Nested2(u[2], u[3])
	x2 := mm.LinearTransform(0., n, uj0)
	fc := mm.LinearTransform(0., n, uj1)
	a := mm.LinearTransform(0., 100., u[4])
	b := mm.LinearTransform(0., 1., u[5])
	l := jointdist.SumToOne(u[6], u[7], u[8])
	return []float64{cv, x1, x2, fc, a, b, l[0], l[1], l[2]}
}

//////////////// Quinn (11)
func sampleQuinn(u []float64) []float64 {
	const sd, n, fc = 1000.0, 0.3, 0.1
	intCap := mm.LinearTransform(0., 0.1, u[0])
	impCap := mm.LinearTransform(0., 0.1, u[1])
	gwCap := mm.LinearTransform(0., 100., u[2])
	fImp := mm.LinearTransform(0., 1., u[3])
	ksat := mm.LogLinearTransform(1e-12, 1., u[4]) // ksat [m/s]
	rootZoneDepth := mm.LinearTransform(0., sd, u[5])
	porosity := mm.LinearTransform(0., n, u[6])
	fieldCap := mm.LinearTransform(0., fc, u[7])
	f := mm.LinearTransform(0., 1., u[8])
	alpha := mm.LinearTransform(0., 1., u[9])
	zwt := mm.LinearTransform(0., 10., u[10])
	return []float64{intCap, impCap, gwCap, fImp, ksat, rootZoneDepth, porosity, fieldCap, f, alpha, zwt}
}

//////////////// SIXPAR (6)
func sampleSIXPAR(u []float64) []float64 {
	upCap := mm.LinearTransform(0., 100., u[0])
	lowCap := mm.LinearTransform(0., 100., u[1])
	upK := mm.LinearTransform(0., 1., u[2])
	lowK := mm.LinearTransform(0., 1., u[3])
	z := mm.LinearTransform(0., 1., u[4])
	x := mm.LinearTransform(0., 1., u[5])
	return []float64{upCap, lowCap, upK, lowK, z, x}
}

//////////////// SPLR (5)
func sampleSPLR(u []float64) []float64 {
	r12 := mm.LinearTransform(0., 10., u[0])
	r23 := mm.LinearTransform(0., 1., u[1])
	k1 := mm.LinearTransform(0., 1., u[2])
	k2 := mm.LinearTransform(0., 1., u[3])
	k3 := mm.LinearTransform(0., 1., u[4])
	return []float64{r12, r23, k1, k2, k3}
}
