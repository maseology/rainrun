package sample

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/maseology/UTM"
	"github.com/maseology/goHydro/solirrad"
	"github.com/maseology/montecarlo"
	mrg63k3a "github.com/maseology/pnrg/MRG63k3a"
	rr "github.com/maseology/rainrun/models"
)

// Sample samples a rainrun model
func Sample(metfp string, nsmpl int, fitness func(o, s []float64) float64) ([][]float64, []float64) {
	rr.LoadMET(metfp, false)

	lat, _, err := UTM.ToLatLon(rr.Loc[1], rr.Loc[2], 17, "", true)
	if err != nil {
		log.Fatalf("%v", err)
	}
	si := solirrad.New(lat, math.Tan(rr.Loc[4]), rr.Loc[5])

	obs := make([]float64, rr.Ndt)
	for i, v := range rr.FRC {
		obs[i] = v[4] // [m/d]
	}

	rng := rand.New(mrg63k3a.New())
	rng.Seed(time.Now().UnixNano())

	ndim := 10
	gen := func(u []float64) float64 {
		var m rr.MakkinkCCFGR4J
		m.New(MakkinkCCFGR4J(u)...)
		m.SI = &si

		f := func(obs []float64) float64 {
			sim := make([]float64, rr.Ndt)
			for i, v := range rr.FRC {
				_, _, r, _ := m.Update(v, rr.DOY[i])
				sim[i] = r
			}
			return fitness(obs[365:], sim[365:])
		}(obs)
		if math.IsNaN(f) {
			// log.Fatalf("Objective function error, u: %v\n", u)
			return -9999.
		}
		return f
	}

	return montecarlo.GenerateSamples(gen, ndim, nsmpl)
}
