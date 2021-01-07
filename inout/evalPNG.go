package inout

import (
	"fmt"

	"github.com/maseology/objfunc"

	"github.com/maseology/mmio"
	rr "github.com/maseology/rainrun/models"
)

// EvalPNG prints model output to a png
func EvalPNG(m rr.Lumper) {
	o := make([]float64, Ndt)
	s := make([]float64, Ndt)
	b := make([]float64, Ndt)
	for i, v := range FRC {
		_, r, g := m.Update(v[0], v[1])
		o[i] = v[2]
		s[i] = r
		b[i] = g
	}
	fmt.Printf(" KGE: %.3f\tNSE: %.3f\tRMSE: %.6f\tmon-wr2: %.3f\tBias: %.3f\n", objfunc.KGE(o[365:], s[365:]), objfunc.NSE(o[365:], s[365:]), objfunc.RMSE(o[365:], s[365:]), objfunc.Krause(o[365:], s[365:]), objfunc.Bias(o[365:], s[365:]))
	mmio.ObsSim("hyd.png", o[365:], s[365:])
	mmio.ObsSimFDC("fdc.png", o[365:], s[365:])
	SumHydrograph(DT, o, s, b)
	SumMonthly(DT, o, s, TS, 1.)
}
