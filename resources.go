package main

import (
	"time"
)

type ResourcesStruct struct {
	Metal   float64
	Silicon float64
	Uranium float64
	Energy  float64
	Time    time.Duration
}

type EquationStruct struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	C float64 `json:"c"`
}

func (r *ResourcesStruct) HasEnoughBasic(other ResourcesStruct) bool {
	return r.Metal >= other.Metal &&
		r.Silicon >= other.Silicon &&
		r.Uranium >= other.Uranium
}

func (r *ResourcesStruct) SubtractBasic(other ResourcesStruct) {
	r.Metal -= other.Metal
	r.Silicon -= other.Silicon
	r.Uranium -= other.Uranium
}

func (r *ResourcesStruct) HasEnoughFull(other ResourcesStruct) bool {
	return r.Metal >= other.Metal &&
		r.Silicon >= other.Silicon &&
		r.Uranium >= other.Uranium &&
		r.Energy >= other.Energy &&
		r.Time >= other.Time
}
func (r *ResourcesStruct) SubtractFull(other ResourcesStruct) {
	r.Metal -= other.Metal
	r.Silicon -= other.Silicon
	r.Uranium -= other.Uranium
	r.Energy -= other.Energy
	r.Time -= other.Time
}
