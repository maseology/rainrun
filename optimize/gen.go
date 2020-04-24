package optimize

import (
	"log"
	"math"

	"github.com/maseology/rainrun/inout"
	rr "github.com/maseology/rainrun/models"
	"github.com/maseology/rainrun/sample"
)

func eval(m rr.Lumper) float64 { // evaluate model
	o := make([]float64, inout.Nfrc)
	s := make([]float64, inout.Nfrc)
	for i, v := range inout.FRC {
		_, r, _ := m.Update(v[0], v[1])
		o[i] = v[2]
		s[i] = r
	}
	return fitness(o[365:], s[365:])
}

func genAtkinson(u []float64) float64 {
	var m rr.Lumper = &rr.Atkinson{}
	m.New(sample.Atkinson(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genDawdyODonnell(u []float64) float64 {
	var m rr.Lumper = &rr.DawdyODonnell{}
	m.New(sample.DawdyODonnell(u, inout.TS)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genGR4J(u []float64) float64 {
	var m rr.Lumper = &rr.GR4J{}
	m.New(sample.GR4J(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genHBV(u []float64) float64 {
	var m rr.Lumper = &rr.HBV{}
	m.New(sample.HBV(u, inout.TS)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genManabeGW(u []float64) float64 {
	var m rr.Lumper = &rr.ManabeGW{}
	m.New(sample.ManabeGW(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genMultiLayerCapacitance(u []float64) float64 {
	var m rr.Lumper = &rr.MultiLayerCapacitance{}
	m.New(sample.MultiLayerCapacitance(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genQuinn(u []float64) float64 {
	var m rr.Lumper = &rr.Quinn{}
	m.New(sample.Quinn(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genSIXPAR(u []float64) float64 {
	var m rr.Lumper = &rr.SIXPAR{}
	m.New(sample.SIXPAR(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genSPLR(u []float64) float64 {
	var m rr.Lumper = &rr.SPLR{}
	m.New(sample.SPLR(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}
