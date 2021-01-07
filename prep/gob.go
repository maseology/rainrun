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

func CreateMetGob() {

	gagfp := "M:/Peel/RDRR-PWRMM21/dat/obs/02HB031.csv"
	cid0 := 1104986

	gobdir := "M:/Peel/GR4J-PWRMM21/PWRMM21."
	gdefFP := "M:/Peel/RDRR-PWRMM21/dat/elevation.real_SWS10.indx.gdef"
	demFP := "M:/Peel/RDRR-PWRMM21/dat/elevation.real.uhdem"
	frcFP := "M:/Peel/RDRR-PWRMM21/PWRMM21.FORC.gob"
	dtb := time.Date(2010, 10, 1, 0, 0, 0, 0, time.UTC)
	dte := time.Date(2020, 9, 30, 18, 0, 0, 0, time.UTC)
	intvl := 86400 * time.Second
	nstep := int(math.Ceil(dte.Sub(dtb).Seconds() / 86400.))

	// get observation
	dts, qs := func() ([]time.Time, []float64) {
		gag, err := mmio.ReadCsvDateValueFlag(gagfp)
		if err != nil {
			log.Fatalf("createMetGob ReadCsvDateValueFlag error: %v", err)
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

	p, ep, ca := func(dts []time.Time) ([]float64, []float64, float64) {
		// get forcings
		fmt.Println("collecting meterological data")
		frc, err := model.LoadGobFORC(frcFP)
		if err != nil {
			log.Fatalf("load forcing error: %v", err)
		}

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
				} else {
					log.Fatalf("error finding met IDs")
				}
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
		_ = fmid

		p, ep := make([]float64, len(dts)), make([]float64, len(dts))
		dxr := make(map[time.Time]int, len(dts))
		// n, fnmid := make([]float64, len(dts)), float64(len(fmid))
		for i, t := range dts {
			dxr[t] = i
		}
		for i, t := range frc.T {
			d := mmio.DayDate(t)
			if ii, ok := dxr[d]; !ok {
				log.Fatalf("error with forcings 2")
			} else {
				for mid, w := range fmid {
					// n[ii]++
					p[ii] += frc.D[0][mid][i] * w
					ep[ii] += frc.D[1][mid][i] * w
				}
			}
		}
		return p, ep, gd.CellArea() * float64(len(cids))
	}(dts)

	dat := make([][]float64, len(dts))
	for i := range dts {
		dat[i] = []float64{p[i], ep[i], qs[i] * 86400. / ca}
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
