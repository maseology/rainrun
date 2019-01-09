package rainrun

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"time"

	"github.com/maseology/glbopt"
	"github.com/maseology/goHydro/met"
	"github.com/maseology/mmio"
	"github.com/maseology/objfunc"
	mrg63k3a "github.com/maseology/pnrg/MRG63k3a"
)

var frc [][]float64
var nfrc int
var fitness = objfunc.KGE

// Optimize to a given model structure
func Optimize(fp string) {
	loadMET(fp)

	rng := rand.New(mrg63k3a.New())
	rng.Seed(time.Now().UnixNano())

	nDim := 7

	uFinal, ofFinal := glbopt.SCE(runtime.GOMAXPROCS(0), nDim, rng, genAtkinson, true)
	// uFinal, _ := glbopt.SurrogateRBF(100, nDim, rng, genAtkinson)

	var m Lumper = &Atkinson{}
	pFinal := m.SampleSpace(uFinal)
	fmt.Printf("\nfinal parameters: %v\n", pFinal)
	fmt.Printf("KGE (best): %.5f\n\n", 1.-ofFinal)

	// m.New(pFinal...)
	// evalPNG(m)
}

// // OptimizeAll to a given model structure
// func OptimizeAll(fp string) {
// 	loadMET(fp)

// 	rng := rand.New(mrg63k3a.New())
// 	rng.Seed(time.Now().UnixNano())

// 	nDim := 7
// 	// uFinal, _ := glbopt.SCE(10, nDim, rng, genAtkinson, false)
// 	uFinal, _ := glbopt.SurrogateRBF(100, nDim, rng, genAtkinson)

// }

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

func evalPNG(m Lumper) {
	o := make([]float64, nfrc)
	s := make([]float64, nfrc)
	for i, v := range frc {
		_, r, _ := m.Update(v[0], v[1])
		o[i] = v[2]
		s[i] = r
	}
	// fmt.Println("KGE (final): ", fitness(o[365:], s[365:]))
	mmio.ObsSim("hyd.png", o, s)
}

/////////////////////////////////////////
//// Function that needs optimizing
/////////////////////////////////////////
func eval(m Lumper) float64 { // evaluate model
	o := make([]float64, nfrc)
	s := make([]float64, nfrc)
	for i, v := range frc {
		_, r, _ := m.Update(v[0], v[1])
		o[i] = v[2]
		s[i] = r
	}
	return 1. - fitness(o[365:], s[365:])
}

func genAtkinson(u []float64) float64 {
	var m Lumper = &Atkinson{}
	m.New(m.SampleSpace(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}
