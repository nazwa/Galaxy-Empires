package ge

import (
	"math"
)

type BuildingMineInterface interface {
	PrecalculateProductionTable(limit int64)
	GetProduction(level int64) int64
}

type BuildingMineStruct struct {
	BuildingStruct

	baseProduction float64

	productionTable     []int64
	productionEquations EquationStruct
}

func (b *BuildingMineStruct) PrecalculateProductionTable(limit int64) {
	b.productionTable = make([]int64, limit)

	for i := int64(0); i < limit; i++ {
		base := b.baseProduction + b.productionEquations.A*float64(i)
		power := math.Pow(b.productionEquations.B, float64(i)*b.productionEquations.C)
		b.productionTable[i] = int64(base * power)
	}
}

func (b *BuildingMineStruct) GetProduction(level int64) int64 {
	if int64(len(b.productionTable)) <= level {
		return 0
	}
	return b.productionTable[level]
}
