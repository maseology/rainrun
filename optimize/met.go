package main

import (
	"log"
	"math"
	"sort"
	"time"

	"github.com/maseology/goHydro/met"
)

var frc [][]float64
var nfrc int

// loadMET collect the climate data, set to a global variable
func loadMET(fp string) {
	nfrc, frc = func() (int, [][]float64) {
		frc, err := met.ReadMET(fp)
		if err != nil {
			log.Fatalln(err)
		}

		dt := make([]time.Time, 0, len(frc))
		for i := range frc {
			dt = append(dt, i)
		}
		sort.Slice(dt, func(i, j int) bool { return dt[i].Before(dt[j]) })

		afrc := make([][]float64, 0, len(dt))
		chk := func(d time.Time, i int, p string) float64 {
			if _, ok := frc[d]; !ok {
				log.Fatalln(d, "not included in met file")
			}
			if v1, ok := frc[d][i]; ok {
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
		return len(afrc), afrc
	}()
}
