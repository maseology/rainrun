package sample

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/im7mortal/UTM"
	"github.com/maseology/goHydro/solirrad"
	"github.com/maseology/montecarlo"
	"github.com/maseology/montecarlo/sampler"
	mrg63k3a "github.com/maseology/pnrg/MRG63k3a"
	io "github.com/maseology/rainrun/inout"
	rr "github.com/maseology/rainrun/models"
)

// Sample samples a rainrun model
func Sample(metfp string, nsmpl int, fitness func(o, s []float64) float64) ([][]float64, []float64) {
	io.LoadMET(metfp, false)

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

	ndim := 10
	ss := sampler.NewSet(MakkinkCCFGR4J())
	gen := func(u []float64) float64 {
		var m rr.MakkinkCCFGR4J
		m.New(ss.Sample(u)...)
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
			// log.Fatalf("Objective function error, u: %v\n", u)
			return -9999.
		}
		return f
	}

	return montecarlo.GenerateSamples(gen, ndim, nsmpl)
}
