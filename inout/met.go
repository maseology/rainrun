package inout

import (
	"log"
	"math"
	"sort"
	"time"

	"github.com/maseology/goHydro/met"
)

// HDR holds header info
var HDR *met.Header

// FRC holds forcing data
var FRC [][]float64

// Nfrc number of timesteps
var Nfrc int
var dt []time.Time

// TS timestep in seconds
var TS float64

// LoadMET collect the climate data, set to a global variable
func LoadMET(fp string, print bool) {
	Nfrc, FRC, HDR = func() (int, [][]float64, *met.Header) {
		h, dc, err := met.ReadMET(fp, print)
		if err != nil {
			log.Fatalln(err)
		}

		TS = h.IntervalSec()
		dt = make([]time.Time, 0, len(dc))
		for i := range dc {
			dt = append(dt, i)
		}
		sort.Slice(dt, func(i, j int) bool { return dt[i].Before(dt[j]) })

		afrc := make([][]float64, 0, len(dt))
		chk := func(d time.Time, i int, p string) float64 {
			if _, ok := dc[d]; !ok {
				log.Fatalln(d, "not included in met file")
			}
			if v1, ok := dc[d][i]; ok {
				return v1
			}
			log.Fatalln(p, "not included in met file")
			return math.NaN()
		}

		for _, d := range dt {
			v := make([]float64, 3)
			v[0] = chk(d, met.AtmosphericYield, "AtmosphericYield")
			v[1] = chk(d, met.AtmosphericDemand, "AtmosphericDemand")
			v[2] = chk(d, met.UnitDischarge, "UnitDischarge")
			afrc = append(afrc, v)
		}
		return len(afrc), afrc, h
	}()
}
