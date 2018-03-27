package rainrun

import "math"

const mingtzero = 1e-10
const minfrac = 1e-4

// Lumper : interface to rainfall-runfall lumped models
type Lumper interface {
	New(p ...float64)
	Update(v ...float64) (float64, float64, float64)
}

// res : simple linear reservoir
type res struct {
	sto, cap, k float64
}

// New constructor
func (r *res) New(capacity, recessionCoef float64) {
	r.cap = capacity
	r.k = recessionCoef
}

// StorageFraction property
func (r *res) StorageFraction() float64 {
	return r.sto / r.cap
}

// Update : update state
// same as Overflow, but does not return Excess
func (r *res) Update(p float64) float64 {
	r.sto += p // allows _sto>_cap
	if r.sto < 0.0 {
		sv := r.sto
		r.sto = 0.0
		return sv
	}
	return 0.
}

// Overflow : update state. p is an net addition
// function returns excess. If p<0 and |p|>sto,
// function returns remainder
func (r *res) Overflow(p float64) float64 {
	r.sto += p
	if r.sto < 0. {
		d := r.sto
		r.sto = 0.
		return d
	} else if r.sto > r.cap {
		d := r.cap - r.sto
		r.sto = r.cap
		return d
	} else {
		return 0.
	}
}

// DecayExp : exponential decay of storage
func (r *res) DecayExp() float64 {
	d := r.k * r.sto
	r.sto -= d
	if r.sto < minfrac {
		d += r.sto
		r.sto = 0.
	}
	return d
}

// DecayMin : exponential decay with minimum storage
func (r *res) DecayMin(minsto float64) float64 {
	d := r.k * r.sto
	r.sto -= d
	if r.sto < minsto {
		sv := r.sto
		r.sto = 2.*minsto - sv //as in PRMS --> water added to reservoir = min - _sto
		// check: should dschrg be reduced by _sto-sto_sv???
	}
	return d
}

// DecayExp2 : exponential decay of storage, with better temporal control
func (r *res) DecayExp2(decay, tsec float64) float64 {
	// see ExponentialDecay.xlsx and Exponential_Decay.docx
	// decay rate givin in m/s
	if decay < mingtzero {
		return 0.
	}
	sv := r.sto
	if r.sto/r.cap < minfrac {
		r.sto = 0.
		return sv
	}
	r.sto *= math.Pow(1.-decay/r.sto, tsec)
	return sv - r.sto
}

func fracCheck(v float64) bool {
	return v < 0. || v > 0.
}
