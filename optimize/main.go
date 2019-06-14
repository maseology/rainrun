package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/maseology/glbopt"
	"github.com/maseology/mmio"
	"github.com/maseology/montecarlo/smpln"
	"github.com/maseology/objfunc"
	mrg63k3a "github.com/maseology/pnrg/MRG63k3a"
	. "github.com/maseology/rainrun"
)

const (
	ncmplx = 10
	nrbf   = 100
)

var fitness = objfunc.LogKGE

func main() {
	start := time.Now()
	defer func() {
		fmt.Println()
		fmt.Println(time.Now().Sub(start))
		fmt.Printf("n processes: %v\n", runtime.GOMAXPROCS(0))
	}()

	optimize("C:/Users/mason/Desktop/CAMC_5000/02EC002.met")

	// x := []float64{8.56, 2.2588, 65.14, 8.2, 938.51, 3.1, -5., 0.0, 0.0, 0.0, -1.0, 5.}
	// fmt.Println(x)
	// x = mmaths.OnlyPositive(x)
	// fmt.Println(x)
	// sort.Float64s(x)
	// fmt.Println(x)
}

func evalPNG(m Lumper) {
	o := make([]float64, nfrc)
	s := make([]float64, nfrc)
	for i, v := range frc {
		_, r, _ := m.Update(v[0], v[1])
		o[i] = v[2]
		s[i] = r
	}
	fmt.Println("KGE (final): ", fitness(o[365:], s[365:]))
	mmio.ObsSim("hyd.png", o[365:1460], s[365:1460])
	mmio.ObsSimFDC("fdc.png", o[365:], s[365:])
}

func optimize(fp string) {
	loadMET(fp)

	rng := rand.New(mrg63k3a.New())
	rng.Seed(time.Now().UnixNano())

	// func() {
	// 	uFinal, _ := glbopt.SCE(ncmplx, 7, rng, genAtkinson, true)
	// 	// uFinal, _ := glbopt.SurrogateRBF(nrbf, 7, rng, genAtkinson)

	// 	var m Lumper = &Atkinson{}
	// 	pFinal := sampleAtkinson(uFinal)
	// 	fmt.Printf("\nfinal parameters: %v\n", pFinal)
	// 	m.New(pFinal...)
	// 	evalPNG(m)
	// }()

	// func() {
	// 	uFinal, _ := glbopt.SCE(ncmplx, 6, rng, genDawdyODonnell, true)
	// 	// uFinal, _ := glbopt.SurrogateRBF(nrbf, 6, rng, genDawdyODonnell)

	// 	var m Lumper = &DawdyODonnell{}
	// 	pFinal := sampleDawdyODonnell(uFinal)
	// 	fmt.Printf("\nfinal parameters: %v\n", pFinal)
	// 	m.New(pFinal...)
	// 	evalPNG(m)
	// }()

	// func() {
	// 	uFinal, _ := glbopt.SCE(ncmplx, 5, rng, genGR4J, true)
	// 	// uFinal, _ := glbopt.SurrogateRBF(nrbf, 5, rng, genGR4J)

	// 	var m Lumper = &GR4J{}
	// 	pFinal := sampleGR4J(uFinal)
	// 	fmt.Printf("\nfinal parameters: %v\n", pFinal)
	// 	m.New(pFinal...)
	// 	evalPNG(m)
	// }()

	// func() {
	// 	uFinal, _ := glbopt.SCE(ncmplx, 10, rng, genHBV, true)
	// 	// uFinal, _ := glbopt.SurrogateRBF(nrbf, 10, rng, genHBV)

	// 	var m Lumper = &HBV{}
	// 	pFinal := sampleHBV(uFinal)
	// 	fmt.Printf("\nfinal parameters: %v\n", pFinal)
	// 	m.New(pFinal...)
	// 	evalPNG(m)
	// }()

	func() {
		// uFinal, _ := glbopt.SCE(ncmplx, 5, rng, genManabeGW, true)
		uFinal, _ := glbopt.SurrogateRBF(nrbf, 5, rng, genManabeGW)

		var m Lumper = &ManabeGW{}
		pFinal := sampleManabeGW(uFinal) // []float64{0.5, 0.5, 0.5, 0.5, 0.5}) //
		fmt.Printf("\nfinal parameters: %v\n", pFinal)
		m.New(pFinal...)
		evalPNG(m)
	}()

	// func() {
	// 	uFinal, _ := glbopt.SCE(ncmplx, 9, rng, genMultiLayerCapacitance, true)
	// 	// uFinal, _ := glbopt.SurrogateRBF(nrbf, 9, rng, genMultiLayerCapacitance)

	// 	var m Lumper = &MultiLayerCapacitance{}
	// 	pFinal := sampleMultiLayerCapacitance(uFinal)
	// 	fmt.Printf("\nfinal parameters: %v\n", pFinal)
	// 	m.New(pFinal...)
	// 	evalPNG(m)
	// }()

	// func() {
	// 	uFinal, _ := glbopt.SCE(ncmplx, 11, rng, genQuinn, true)
	// 	// uFinal, _ := glbopt.SurrogateRBF(nrbf, 11, rng, genQuinn)

	// 	var m Lumper = &Quinn{}
	// 	pFinal := sampleQuinn(uFinal)
	// 	fmt.Printf("\nfinal parameters: %v\n", pFinal)
	// 	m.New(pFinal...)
	// 	evalPNG(m)
	// }()

	// func() {
	// 	uFinal, _ := glbopt.SCE(ncmplx, 6, rng, genSIXPAR, true)
	// 	// uFinal, _ := glbopt.SurrogateRBF(nrbf, 6, rng, genSIXPAR)

	// 	var m Lumper = &SIXPAR{}
	// 	pFinal := sampleSIXPAR(uFinal)
	// 	fmt.Printf("\nfinal parameters: %v\n", pFinal)
	// 	m.New(pFinal...)
	// 	evalPNG(m)
	// }()

	// func() {
	// 	uFinal, _ := glbopt.SCE(ncmplx, 6, rng, genSPLR, true)
	// 	// uFinal, _ := glbopt.SurrogateRBF(nrbf, 6, rng, genSPLR)

	// 	var m Lumper = &SPLR{}
	// 	pFinal := sampleSPLR(uFinal)
	// 	fmt.Printf("\nfinal parameters: %v\n", pFinal)
	// 	m.New(pFinal...)
	// 	evalPNG(m)
	// }()

}

func permute(fp string) {
	loadMET(fp)
	var m Lumper = &DawdyODonnell{}
	for i, u := range smpln.Permutations(6, 3) {
		fmt.Println(i, u)
		m.New(sampleDawdyODonnell(u)...)
		if math.IsNaN(eval(m)) {
			panic("NaN")
		}
	}
}
