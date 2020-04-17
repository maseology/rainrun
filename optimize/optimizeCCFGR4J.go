package optimize

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/maseology/glbopt"
	"github.com/maseology/goHydro/solirrad"
	"github.com/maseology/mmio"
	"github.com/maseology/objfunc"
	mrg63k3a "github.com/maseology/pnrg/MRG63k3a"
	io "github.com/maseology/rainrun/inout"
	rr "github.com/maseology/rainrun/models"
)

// CCFGR4J a single or set of rainrun models
func CCFGR4J(fp, logfp string) {
	logger := mmio.GetInstance(logfp)
	io.LoadMET(fp, true)

	si := solirrad.New(43.6, 0., 0.)

	obs := make([]float64, io.Nfrc)
	for i, v := range io.FRC {
		obs[i] = v[4] // [m/d]
	}

	rng := rand.New(mrg63k3a.New())
	rng.Seed(time.Now().UnixNano())

	genCCFGR4J := func(u []float64) float64 {
		var m rr.CCFGR4J
		m.New(sampleCCFGR4J(u)...)
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

	uFinal, _ := glbopt.SCE(ncmplx, 8, rng, genCCFGR4J, true)
	// uFinal, _ := glbopt.SurrogateRBF(nrbf, 8, rng, genCCFGR4J)

	func() {
		par := []string{"x1", "x2", "x3", "x4", "tindex", "ddfc", "baseT", "tsf"}
		pFinal := sampleCCFGR4J(uFinal)
		fmt.Println("Optimum:")
		for i, v := range par {
			fmt.Printf(" %s:\t\t%.4f\t[%.4e]\n", v, pFinal[i], uFinal[i])
		}

		var m rr.CCFGR4J
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
