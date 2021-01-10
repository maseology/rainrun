package optimize

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/maseology/glbopt"
	"github.com/maseology/mmio"
	"github.com/maseology/montecarlo/smpln"
	"github.com/maseology/objfunc"
	mrg63k3a "github.com/maseology/pnrg/MRG63k3a"
	rr "github.com/maseology/rainrun/models"
	"github.com/maseology/rainrun/sample"
)

const (
	nrbf   = 100
	ncmplx = 200
)

var minimizer = func(o, s []float64) float64 { return 1. - objfunc.NSE(o, s) }

// Optimize a single or set of rainrun models
func Optimize(fp, mdl, logfp string) {
	logger := mmio.GetInstance(logfp)
	rr.LoadMET(fp, true)

	rng := rand.New(mrg63k3a.New())
	rng.Seed(time.Now().UnixNano())

	const nrbf = 100
	ncmplx := 64

	switch mdl {
	case "Atkinson":
		func() {
			uFinal, _ := glbopt.SCE(ncmplx, 7, rng, genAtkinson, true)
			// uFinal, _ := glbopt.SurrogateRBF(nrbf, 7, rng, genAtkinson)

			var m rr.Lumper = &rr.Atkinson{}
			pFinal := sample.Atkinson(uFinal)
			fmt.Printf("\nfinal parameters: %v\n", pFinal)
			fmt.Printf("sample space:\t%f\n", uFinal)
			m.New(pFinal...)
			rr.EvalPNG(m)
		}()
	case "DawdyODonnell":
		func() {
			uFinal, _ := glbopt.SCE(ncmplx, 6, rng, genDawdyODonnell, true)
			// uFinal, _ := glbopt.SurrogateRBF(nrbf, 6, rng, genDawdyODonnell)

			var m rr.Lumper = &rr.DawdyODonnell{}
			pFinal := sample.DawdyODonnell(uFinal, rr.TS)
			fmt.Printf("\nfinal parameters: %v\n", pFinal)
			fmt.Printf("sample space:\t%f\n", uFinal)
			m.New(pFinal...)
			rr.EvalPNG(m)
		}()
	case "GR4J":
		func() {
			uFinal, _ := glbopt.SCE(ncmplx, 4, rng, genGR4J, true)
			// uFinal, _ := glbopt.SurrogateRBF(nrbf, 4, rng, genGR4J)

			var m rr.Lumper = &rr.GR4J{}
			pFinal := sample.GR4J(uFinal)
			sp := fmt.Sprintf("\nfinal parameters:\t%.3e\n", pFinal)
			su := fmt.Sprintf("sample space:\t\t%f\n", uFinal)
			fmt.Print(sp + su)
			m.New(pFinal...)
			logger.Println(mmio.FileName(fp, false))
			logger.Print(sp + su)
			logger.Println("\n" + rr.EvalPNG(m))
			// var m rr.Lumper = &rr.GR4J{}
			// ss := sampler.NewSet(sample.GR4J()) //////////////////////////////////  TO FIX
			// pFinal := ss.Sample(uFinal)
			// fmt.Printf("\nfinal parameters: %v\n", pFinal)
			// m.New(pFinal...)
			// rr.EvalPNG(m)
		}()
	case "HBV":
		func() {
			uFinal, _ := glbopt.SCE(ncmplx, 9, rng, genHBV, true)
			// uFinal, _ := glbopt.SurrogateRBF(nrbf, 9, rng, genHBV)

			var m rr.Lumper = &rr.HBV{}
			pFinal := sample.HBV(uFinal, rr.TS)
			sp := fmt.Sprintf("\nfinal parameters:\t%.3e\n", pFinal)
			su := fmt.Sprintf("sample space:\t\t%f\n", uFinal)
			m.New(pFinal...)
			logger.Println(mmio.FileName(fp, false))
			logger.Print(sp + su)
			logger.Println("\n" + rr.EvalPNG(m))
		}()
	case "ManabeGW":
		func() { // check
			uFinal, _ := glbopt.SCE(ncmplx, 5, rng, genManabeGW, true)
			// uFinal, _ := glbopt.SurrogateRBF(nrbf, 5, rng, genManabeGW)

			var m rr.Lumper = &rr.ManabeGW{}
			pFinal := sample.ManabeGW(uFinal)
			fmt.Printf("\nfinal parameters: %v\n", pFinal)
			fmt.Printf("sample space:\t\t%f\n", uFinal)
			m.New(pFinal...)
			rr.EvalPNG(m)
		}()
	case "MultiLayerCapacitance":
		func() { // check
			uFinal, _ := glbopt.SCE(ncmplx, 9, rng, genMultiLayerCapacitance, true)
			// uFinal, _ := glbopt.SurrogateRBF(nrbf, 9, rng, genMultiLayerCapacitance)

			var m rr.Lumper = &rr.MultiLayerCapacitance{}
			pFinal := sample.MultiLayerCapacitance(uFinal)
			fmt.Printf("\nfinal parameters: %v\n", pFinal)
			fmt.Printf("sample space:\t\t%f\n", uFinal)
			m.New(pFinal...)
			rr.EvalPNG(m)
		}()
	case "Quinn":
		func() { // check
			uFinal, _ := glbopt.SCE(ncmplx, 11, rng, genQuinn, true)
			// uFinal, _ := glbopt.SurrogateRBF(nrbf, 11, rng, genQuinn)

			var m rr.Lumper = &rr.Quinn{}
			pFinal := sample.Quinn(uFinal)
			fmt.Printf("\nfinal parameters: %v\n", pFinal)
			fmt.Printf("sample space:\t\t%f\n", uFinal)
			m.New(pFinal...)
			rr.EvalPNG(m)
		}()
	case "SIXPAR":
		func() { // check
			uFinal, _ := glbopt.SCE(ncmplx, 6, rng, genSIXPAR, true)
			// uFinal, _ := glbopt.SurrogateRBF(nrbf, 6, rng, genSIXPAR)

			var m rr.Lumper = &rr.SIXPAR{}
			pFinal := sample.SIXPAR(uFinal)
			fmt.Printf("\nfinal parameters: %v\n", pFinal)
			fmt.Printf("sample space:\t\t%f\n", uFinal)
			m.New(pFinal...)
			rr.EvalPNG(m)
		}()
	case "SPLR":
		func() { // check (negative AET)
			uFinal, _ := glbopt.SCE(ncmplx, 6, rng, genSPLR, true)
			// uFinal, _ := glbopt.SurrogateRBF(nrbf, 6, rng, genSPLR)

			var m rr.Lumper = &rr.SPLR{}
			pFinal := sample.SPLR(uFinal)
			fmt.Printf("\nfinal parameters: %v\n", pFinal)
			fmt.Printf("sample space:\t\t%f\n", uFinal)
			m.New(pFinal...)
			rr.EvalPNG(m)
		}()
	default:
		fmt.Println("unrecognized model:" + mdl)
	}
}

// permute used to create a complete sample set of
// every possible permutation of p dimensions and w discrete
// values.
func permute(fp string) {
	rr.LoadMET(fp, true)
	var m rr.Lumper = &rr.DawdyODonnell{}
	for i, u := range smpln.Permutations(6, 3) {
		fmt.Println(i, u)
		m.New(sample.DawdyODonnell(u, rr.TS)...)
		if math.IsNaN(eval(m)) {
			panic("NaN")
		}
	}
}
