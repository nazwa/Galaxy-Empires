package ge

import (
	"math"
	"time"
)

// Constant IDs for different building types
const (
	BuildingTypeMine geBuildingType = iota + 1
	BuildingTypePower
	BuildingTypeStorage
	BuildingTypeGeneric
)

// List of basic building IDs
const (
	BuildingIdMineMetal geBuildingID = iota + 1
	BuildingIdMineSilicon
	BuildingIdMineUranium

	BuildingIdPowerPlant

	BuildingIdGenericFactory
	BuildingIdGenericShipyard
	BuildingIdGenericLab
	BuildingIdGenericCommand
)

type BuildingInterface interface {
	GetCost(level int64) ResourcesStruct
	PrecalculateCostTable(levels int64)
}

type BuildingStruct struct {
	ID   geBuildingID
	Type geBuildingType

	BaseCost     ResourcesStruct
	Requirements RequirementsStruct

	CostEquations EquationStruct

	costTable []ResourcesStruct
}

type BuildingLevelStruct struct {
	Building BuildingInterface
	Level    int64
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
}

func (b *BuildingStruct) GetCost(level int64) ResourcesStruct {
	return b.costTable[level]
}
