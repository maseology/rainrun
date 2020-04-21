package sample

import (
	mm "github.com/maseology/mmaths"
)

//////////////// Makkink (4)
func sampleMakkink(u []float64) []float64 {
	// a := 1. // a and alpha are effectively the same
	b := mm.LinearTransform(0., .1, u[0])
	c := mm.LinearTransform(0., 5., u[1])
	alpha := mm.LinearTransform(0., 5., u[2])
	beta := mm.LinearTransform(-.0005, 0.0005, u[3])
	return []float64{b, c, alpha, beta}
}

//////////////// CCF (4)
func sampleCCF(u []float64) []float64 {
	tindex := mm.LogLinearTransform(0.0002, 0.05, u[4]) // CCF temperature index; range .0002 to 0.0005 m/°C/d -- roughly 1/10 DDF (pg.278)
	ddfc := mm.LinearTransform(0., 10., u[5])           // DDF adjustment factor based on pack density, see DeWalle and Rango, pg. 275; Ref: Martinec (1960)=1.1
	baseT := mm.LinearTransform(-5., 5., u[6])          // base/critical temperature (°C)
	tsf := mm.LinearTransform(0.1, 0.7, u[7])           // TSF (surface temperature factor), 0.1-0.5 have been used
	// ddf := mm.LinearTransform(0.001, 0.008, u[1])    // (initial) degree-day/melt factor; range .001 to .008 m/°C/d  (pg.275)
	return []float64{tindex, ddfc, baseT, tsf}
}

//////////////// GR4J (4)
func sampleGR4J(u []float64) []float64 {
	prd := mm.LinearTransform(0., 1., u[0]) // x1: "production storage" capacity (mm)
	x2 := mm.LinearTransform(-.1, .1, u[1]) // x2: water exchange coefficient (>0 for water imports, <0 for exports, =0 for no exchange)
	rte := mm.LinearTransform(0., 3., u[2]) // x3: "routing storage"/groundwater storage capacity (mm)
	x4 := mm.LinearTransform(.5, 10., u[3]) // x4: unit hydrograph time base (days)
	// qsplt := mm.LinearTransform(0., 1., u[4]) // fixed in paper as 0.9
	return []float64{prd, x2, rte, x4} //, qsplt}
}

//////////////// CCFGR4J (8)
func sampleCCFGR4J(u []float64) []float64 {
	ugr4j := sampleGR4J(u)
	uccf := sampleCCF(u[4:])
	return append(ugr4j, uccf...)
}

//////////////// MakkinkCCFGR4J (12)
func sampleMakkinkCCFGR4J(u []float64) []float64 {
	ugr4j := sampleGR4J(u)
	uccf := sampleCCF(u[4:8])
	mak := sampleMakkink(u[8:])
	return append(ugr4j, append(uccf, mak...)...)
}
