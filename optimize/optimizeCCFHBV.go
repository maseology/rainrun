package optimize

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/im7mortal/UTM"
	"github.com/maseology/glbopt"
	"github.com/maseology/goHydro/solirrad"
	"github.com/maseology/mmio"
	"github.com/maseology/objfunc"
	mrg63k3a "github.com/maseology/pnrg/MRG63k3a"
	io "github.com/maseology/rainrun/inout"
	rr "github.com/maseology/rainrun/models"
)

// CCFHBV a single or set of rainrun models
func CCFHBV(fp, logfp string) {
	logger := mmio.GetInstance(logfp)
	io.LoadMET(fp, true)

	lat, _, err := UTM.ToLatLon(io.Loc[1], io.Loc[2], 17, "", true)
	if err != nil {
		log.Fatalf("%v", err)
	}
	si := solirrad.New(lat, math.Tan(io.Loc[4]), io.Loc[5])

	obs := make([]float64, io.Nfrc)
	for i, v := range io.FRC {
		obs[i] = v[4] // [m/d]
	}

	rng := rand.New(mrg63k3a.New())
	rng.Seed(time.Now().UnixNano())

	genCCFHBV := func(u []float64) float64 {
		var m rr.CCFHBV
		m.New(sampleCCFHBV(u, io.TS)...)
		m.SI = &si

		f := func(obs []float64) float64 {
			sim := make([]float64, io.Nfrc)
			for i, v := range io.FRC {
				_, _, r, _ := m.Update(v, io.DOY[i])
				sim[i] = r
			}
			return fitness(obs[365:], sim[365:])
		}(obs)
		if math.IsNaN(f) {
			log.Fatalf("Objective function error, u: %v\n", u)
		}
		return f
	}

	uFinal, _ := glbopt.SCE(ncmplx, 13, rng, genCCFHBV, true)
	// uFinal, _ := glbopt.SurrogateRBF(nrbf, 13, rng, genCCFHBV)

	func() {

		// uFinal := []float64{0.36, 0.86, 0.20, 0.99, 0.74, 0.71, 0.28, 0.78, 0.37, 0.63, 0.3, 0.92, 0.52}
		par := []string{"fc", "lp", "beta", "uzl", "k0", "k1", "k2", "perc", "maxbas", "tindex", "ddfc", "baseT", "tsf"}
		pFinal := sampleCCFHBV(uFinal, io.TS)
		fmt.Println("Optimum:")
		for i, v := range par {
			fmt.Printf(" %s:\t\t%.4f\t[%.4e]\n", v, pFinal[i], uFinal[i])
		}

		// fmt.Println("\nparameter names:\t[fc lp beta uzl k0 k1 k2 perc maxbas tindex ddfc baseT tsf]")
		// fmt.Printf("final parameters:\t%.3e\n", pFinal)
		// fmt.Printf("sample space:\t\t%f\n", uFinal)

		var m rr.CCFHBV
		m.SI = &si
		m.New(pFinal...)
		sim, aet, bf := make([]float64, io.Nfrc), make([]float64, io.Nfrc), make([]float64, io.Nfrc)
		y := make([]float64, io.Nfrc)
		for i, v := range io.FRC {
			yy, a, r, g := m.Update(v, io.DOY[i])
			y[i] = yy
			aet[i] = a
			sim[i] = r
			bf[i] = g
		}
		kge, nse, mwr2, bias := objfunc.KGE(obs[365:], sim[365:]), objfunc.NSE(obs[365:], sim[365:]), objfunc.Krause(obs[365:], sim[365:]), objfunc.Bias(obs[365:], sim[365:])
		fmt.Printf(" KGE: %.3f\tNSE: %.3f\tmon-wr2: %.3f\tBias: %.3f\n", kge, nse, mwr2, bias)
		func() {
			idt, iy, ia, iob, is, ig := make([]interface{}, io.Nfrc), make([]interface{}, io.Nfrc), make([]interface{}, io.Nfrc), make([]interface{}, io.Nfrc), make([]interface{}, io.Nfrc), make([]interface{}, io.Nfrc)
			for i := range obs {
				idt[i] = io.DT[i]
				iy[i] = y[i]
				ia[i] = aet[i]
				iob[i] = obs[i]
				is[i] = sim[i]
				ig[i] = bf[i]
			}
			mmio.WriteCSV(mmio.RemoveExtension(fp)+".hydrograph.csv", "date,y,aet,obs,sim,bf", idt, iy, ia, iob, is, ig)
			logger.Println(fmt.Sprintf("\nnam\t%v\nU\t%v\nP\t%v\nKGE\t%f\nNSE\t%f\nmwr2\t%f\nbias\t%f\n", par, uFinal, pFinal, kge, nse, mwr2, bias))
		}()
	}()
}
