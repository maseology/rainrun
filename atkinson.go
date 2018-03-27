package rainrun

import "math"

// Atkinson simple storage model
// based on formulation given in: Atkinson, S.E., M. Sivapalan, N.R. Viney, R.A. Woods, 2003. Predicting space-time variability of hourly streamflow and the role of climate seasonality: Mahurangi Catchment, New Zealand. Hydrological Processes 17. pp. 2171-2193.
// original ref: Atkinson S.E., R.A. Woods, M. Sivapalan, 2002. Climate and landscape controls on water balance model complexity over changing timescales. Water Resource Research 38(12): 1314.
// additional ref: Wittenberg H., M. Sivapalan, 1999. Watershed groundwater balance equation using streamflow recession analysis and baseflow separation. Journal of Hydrology 219, pp.20-33.
type Atkinson struct {
	sto, sint, sintc, cov, kb, a, b, sbc, sfc float64
}

// New constructor
// sto: current storage; sint current interception storage; cov: fractional forest cover; kb = 1/Tcbf
func (m *Atkinson) New(soildepth, porosity, fc, wp, coverdense, intcap, kb, a, b float64) {
	if coverdense < 0. || coverdense > 1. {
		panic("Atkinson input error")
	}
	m.cov = coverdense // cover density
	m.sintc = intcap   // interception storage capacity
	m.sbc = soildepth * (porosity - wp)
	m.sfc = soildepth * (fc - wp) // originally written as Sbc * (fc - wp) / (n - wp)
	m.kb = kb                     // baseflow recession coefficient
	m.a = a                       // sub-surface flow coefficient
	m.b = b                       // sub-surface flow coefficient
}

// Update state
func (m *Atkinson) Update(p, ep float64) (float64, float64, float64) {
	g := m.sto         // saving antecedent storage
	eveg := m.cov * ep // transpiration
	if m.sto < m.sfc {
		eveg *= m.sto / m.sfc
	}
	ebs := (1. - m.cov) * ep // bare soil evaporation
	if m.sto < m.sbc {
		ebs *= m.sto / m.sbc
	}
	var qse, qss float64 // saturation excess, subsurface runoff
	if m.sto > m.sbc {
		qse = m.sto - m.sbc
	}
	if m.sto < m.sbc {
		qss = math.Pow((m.sto-m.sfc)/m.a, 1./m.b) // (Wittenberg and Sivapalan, 1999)
	} else if m.sto > m.sfc {
		qss = math.Pow((m.sbc-m.sfc)/m.a, 1./m.b) // (Wittenberg and Sivapalan, 1999)
	}

	qbf := m.sto * m.kb // baseflow
	eint := ep          // interception evaporation
	if p+m.sint < ep {
		eint = p + m.sint
	}

	a := eveg + ebs + eint // total actual ET
	var thr float64        // throughflow
	if p > m.sintc-m.sint {
		thr = p - (m.sintc - m.sint)
	}

	m.sint += p - eint - thr                    // interception water balance
	m.sto += thr - eveg - ebs - qse - qss - qbf // soil zone water balance
	g -= m.sto
	if g < 0. {
		g = 0. // not part of the Atkinson model, but closest assumption to infiltration (=rechage)
	}
	q := qse + qss + qbf // total discharge

	return a, q, g
}
