package inout

import (
	"log"
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

// DT holds dates
var DT []time.Time

// DOY hold day of year
var DOY []int

// TS timestep in seconds
var TS float64

// LoadMET collect the climate data, set to a global variable
func LoadMET(fp string, print bool) {
	Nfrc, FRC, HDR = func() (int, [][]float64, *met.Header) {
		h, c, err := met.ReadMET(fp, print)
		if err != nil {
			log.Fatalln(err)
		}

		TS = h.IntervalSec()
		DT = make([]time.Time, 0, len(c.T))
		for _, t := range c.T {
			DT = append(DT, t)
		}
		sort.Slice(DT, func(i, j int) bool { return DT[i].Before(DT[j]) })

		DOY = make([]int, len(DT))
		for i, t := range DT {
			DOY[i] = t.YearDay()
		}

		afrc := make([][]float64, 0, len(DT))
		switch h.WBCD {
		case 33554486:
			for i := range DT {
				afrc = append(afrc, []float64{c.D[i][0][0], c.D[i][0][1], c.D[i][0][2], c.D[i][0][3], c.D[i][0][4]})
			}
		case 33555968:
			log.Fatalf("met.go LoadMET(): FIX CODE\n")
			// chk := func(d time.Time, i int, p string) float64 {
			// 	if _, ok := dc[d]; !ok {
			// 		log.Fatalln(d, "not included in met file")
			// 	}
			// 	if v1, ok := dc[d][i]; ok {
			// 		return v1
			// 	}
			// 	log.Fatalln(p, "not included in met file")
			// 	return math.NaN()
			// }

			// for _, d := range dt { // [timeID][0][TypeID]
			// 	v := make([]float64, 3)
			// 	v[0] = chk(d, met.AtmosphericYield, "AtmosphericYield")
			// 	v[1] = chk(d, met.AtmosphericDemand, "AtmosphericDemand")
			// 	v[2] = chk(d, met.UnitDischarge, "UnitDischarge")
			// 	afrc = append(afrc, v)
			// }
		default:
			log.Fatalf("met.go LoadMET() error: WBCD code not supported: %d\n", h.WBCD)
		}

		return len(afrc), afrc, h
	}()
}
