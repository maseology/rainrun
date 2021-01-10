package prep

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/maseology/goHydro/grid"
	"github.com/maseology/goHydro/tem"
	"github.com/maseology/mmio"
	"github.com/maseology/rdrr/model"
)

func CreateMetGob(gobdir, gagfp, frcFP string, cid0 int) {

	gdefFP := "M:/Peel/RDRR-PWRMM21/dat/elevation.real_SWS10.indx.gdef"
	demFP := "M:/Peel/RDRR-PWRMM21/dat/elevation.real.uhdem"

	dtb := time.Date(2010, 10, 1, 0, 0, 0, 0, time.UTC)
	dte := time.Date(2020, 9, 30, 18, 0, 0, 0, time.UTC)
	intvl := 86400 * time.Second
	nstep := int(math.Ceil(dte.Sub(dtb).Seconds() / 86400.))

	// get observation
	dts, qs := func() ([]time.Time, []float64) {
		gag, err := mmio.ReadCsvDateFloat(gagfp)
		if err != nil {
			log.Fatalf("createMetGob ReadCsvDateFloat error: %v", err)
		}
		ii, dts, qs := 0, make([]time.Time, nstep), make([]float64, nstep)
		for t := dtb; !t.After(dte); t = t.Add(intvl) {
			d := mmio.DayDate(t)
			dts[ii] = d
			if v, ok := gag[d]; ok {
				qs[ii] = v
			} else {
				qs[ii] = 0.
			}
			ii++
		}
		return dts, qs
	}()

	vals, ca := func(dts []time.Time) ([][]float64, float64) {

		// get grid definition
		fmt.Println("\ncollecting forcings..")
		gd := func() *grid.Definition {
			gd, err := grid.ReadGDEF(gdefFP, true)
			if err != nil {
				log.Fatalf("%v", err)
			}
			if len(gd.Sactives) <= 0 {
				log.Fatalf("error: grid definition requires active cells")
			}
			return gd
		}()

		// get forcings
		fmt.Println("collecting meterological data")
		var frc *model.FORC
		var err error
		npar := 2
		if len(frcFP) > 0 { // .gob
			if frc, err = model.LoadGobFORC(frcFP); err != nil {
				log.Fatalf("load forcing error: %v", err)
			}
		} else { // ORMGP_5000_YCDB.met
			npar = 4
			if frc, err = loadMetFORC(gd); err != nil {
				log.Fatalf("load forcing error: %v", err)
			}
		}

		// get catchment cells
		fmt.Println("collecting catchment cells")
		cids := func() []int {
			var dem tem.TEM
			if err := dem.New(demFP); err != nil {
				log.Fatalf(" tem.New() error: %v", err)
			}
			return dem.ContributingAreaIDs(cid0)
		}()

		// get met IDs
		fmt.Println("collecting meterological IDs")
		fmid := func(cids []int) map[int]float64 {
			mc := make(map[int]int)
			for _, cid := range cids {
				if mid, ok := frc.XR[cid]; ok {
					if _, ok := mc[mid]; ok {
						mc[mid]++
					} else {
						mc[mid] = 1
					}
				} // else {
				// 	log.Fatalf("error finding met IDs")
				// }
			}
			fmc, dnm := make(map[int]float64, len(mc)), 0.
			for i, v := range mc {
				fv := float64(v)
				fmc[i] = fv
				dnm += fv
			}
			for i, v := range fmc {
				fmc[i] = v / dnm
			}
			return fmc
		}(cids)

		dxr := make(map[time.Time]int, len(dts))
		for i, t := range dts {
			dxr[t] = i
		}
		// if npar == 2 {
		// 	p, ep := make([]float64, len(dts)), make([]float64, len(dts))
		// 	for i, t := range frc.T {
		// 		d := mmio.DayDate(t)
		// 		if ii, ok := dxr[d]; !ok {
		// 			continue
		// 		} else {
		// 			for mid, w := range fmid {
		// 				p[ii] += frc.D[0][mid][i] * w
		// 				ep[ii] += frc.D[1][mid][i] * w
		// 			}
		// 		}
		// 	}
		// 	return [][]float64{p, ep}, gd.CellArea() * float64(len(cids))
		// }
		vs := make([][]float64, len(dts))
		// n := make([]float64, len(dts))
		for i := 0; i < len(dts); i++ {
			vs[i] = make([]float64, npar)
		}
		for i, t := range frc.T {
			d := mmio.DayDate(t)
			if ii, ok := dxr[d]; !ok {
				continue
			} else {
				for mid, w := range fmid {
					for k := 0; k < npar; k++ {
						vs[ii][k] += frc.D[k][mid][i] * w
					}
				}
				// n[ii]++
			}
		}
		// for ii := range n {
		// 	for k := 0; k < npar; k++ {
		// 		vs[ii][k] /= n[ii]
		// 	}
		// }
		return vs, gd.CellArea() * float64(len(cids))
	}(dts)

	dat := make([][]float64, len(dts))
	for i := range dts {
		dat[i] = append(vals[i], qs[i]*86400./ca)
	}

	// save gob
	fmt.Println("saving GOB")
	f, err := os.Create(gobdir + mmio.FileName(gagfp, false) + ".gob")
	defer f.Close()
	if err != nil {
		log.Fatalf("createMetGob saveGob error: %v", err)
	}
	if err := gob.NewEncoder(f).Encode(dat); err != nil {
		log.Fatalf("createMetGob saveGob error: %v", err)
	}
	if err := gob.NewEncoder(f).Encode(dts); err != nil {
		log.Fatalf("createMetGob saveGob error: %v", err)
	}
}
