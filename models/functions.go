package rainrun

import "math"

func etRadToGlobal(Ke, tx, tn float64) float64 {
	// see pg 151 in DeWalle & Rango; attributed to Bristow and Campbell (1984)
	// ref: Bristow, K.L. and G.S. Campbell, 1984. On the relationship between incoming solar radiation and daily maximum and minimum temperature. Agricultural and Forest Meteorology 31(2):159--166.
	return Ke * a * (1. - math.Exp(-b*math.Pow(tx-tn, g)))
}
