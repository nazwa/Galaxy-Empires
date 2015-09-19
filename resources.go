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
