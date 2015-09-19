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
	a float64
	b float64
	c float64
}
