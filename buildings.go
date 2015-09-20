package main

import (
	"fmt"
	"math"
	"time"
)

type BuildingStruct struct {
	ID             string
	Name           string
	Category       string
	Produces       string
	Description    string
	BaseCost       ResourcesStruct
	BaseProduction float64
	Requirements   RequirementsStruct

	CostEquations       EquationStruct
	BuildEqations       EquationStruct
	ProductionEquations EquationStruct

	costTable       []ResourcesStruct
	productionTable []int64
}

func (b *BuildingStruct) PrecalculateProductionTable(limit int64) {
	b.productionTable = make([]int64, limit)

	for i := int64(0); i < limit; i++ {
		base := b.BaseProduction * b.ProductionEquations.A * float64(i)
		power := math.Pow(b.ProductionEquations.B, float64(i)*b.ProductionEquations.C)
		b.productionTable[i] = int64(base * power)
	}
}

func (b *BuildingStruct) PrecalculateCostTable(limit int64) {
	b.costTable = make([]ResourcesStruct, limit)

	for i := int64(0); i < limit; i++ {
		factor := b.CostEquations.A * math.Pow(b.CostEquations.B, float64(i-1)*b.CostEquations.C)
		b.costTable[i] = ResourcesStruct{
			Metal:   b.BaseCost.Metal * factor,
			Silicon: b.BaseCost.Silicon * factor,
			Uranium: b.BaseCost.Uranium * factor,
			Energy:  b.BaseCost.Energy * factor,
			Time:    time.Duration(float64(b.BaseCost.Time*time.Second) * factor),
		}
	}
	fmt.Println(b.Name)
	fmt.Println(b.costTable[1])
	fmt.Println(b.costTable[2])
	fmt.Println("--")
	time.Now()
}

func (b *BuildingStruct) GetCost(level int64) ResourcesStruct {
	return b.costTable[level]
}

func (b *BuildingStruct) GetProduction(level int64) int64 {
	return b.productionTable[level]
}
