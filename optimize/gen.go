package optimize

import (
	"log"
	"math"

	"github.com/maseology/rainrun/inout"
	rr "github.com/maseology/rainrun/models"
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
	m.New(sampleAtkinson(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genDawdyODonnell(u []float64) float64 {
	var m rr.Lumper = &rr.DawdyODonnell{}
	m.New(sampleDawdyODonnell(u, inout.TS)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genGR4J(u []float64) float64 {
	var m rr.Lumper = &rr.GR4J{}
	m.New(sampleGR4J(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genHBV(u []float64) float64 {
	var m rr.Lumper = &rr.HBV{}
	m.New(sampleHBV(u, inout.TS)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genManabeGW(u []float64) float64 {
	var m rr.Lumper = &rr.ManabeGW{}
	m.New(sampleManabeGW(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genMultiLayerCapacitance(u []float64) float64 {
	var m rr.Lumper = &rr.MultiLayerCapacitance{}
	m.New(sampleMultiLayerCapacitance(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genQuinn(u []float64) float64 {
	var m rr.Lumper = &rr.Quinn{}
	m.New(sampleQuinn(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genSIXPAR(u []float64) float64 {
	var m rr.Lumper = &rr.SIXPAR{}
	m.New(sampleSIXPAR(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genSPLR(u []float64) float64 {
	var m rr.Lumper = &rr.SPLR{}
	m.New(sampleSPLR(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}
