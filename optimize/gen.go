package main

import (
	"log"
	"math"

	. "github.com/maseology/rainrun"
)

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
	m.New(sampleAtkinson(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genDawdyODonnell(u []float64) float64 {
	var m Lumper = &DawdyODonnell{}
	m.New(sampleDawdyODonnell(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genGR4J(u []float64) float64 {
	var m Lumper = &GR4J{}
	m.New(sampleGR4J(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genHBV(u []float64) float64 {
	var m Lumper = &HBV{}
	m.New(sampleHBV(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genManabeGW(u []float64) float64 {
	var m Lumper = &ManabeGW{}
	m.New(sampleManabeGW(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genMultiLayerCapacitance(u []float64) float64 {
	var m Lumper = &MultiLayerCapacitance{}
	m.New(sampleMultiLayerCapacitance(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genQuinn(u []float64) float64 {
	var m Lumper = &Quinn{}
	m.New(sampleQuinn(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genSIXPAR(u []float64) float64 {
	var m Lumper = &SIXPAR{}
	m.New(sampleSIXPAR(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}

func genSPLR(u []float64) float64 {
	var m Lumper = &SPLR{}
	m.New(sampleSPLR(u)...)
	f := eval(m)
	if math.IsNaN(f) {
		log.Fatalf("Objective function error, u: %v\n", u)
	}
	return f
}
