package ge

import (
	"math"
	"time"
)

type geBuildingType int64
type geBuildingID int64

// Constant IDs for different building types
const (
	BuildingTypeMine    geBuildingType = iota + 1
	BuildingTypePower   geBuildingType = iota + 1
	BuildingTypeStorage geBuildingType = iota + 1
	BuildingTypeGeneric geBuildingType = iota + 1
)

// List of basic building IDs
const (
	BuildingIdMineMetal   geBuildingID = iota + 1
	BuildingIdMineSilicon geBuildingID = iota + 1
	BuildingIdMineUranium geBuildingID = iota + 1

	BuildingIdPowerPlant geBuildingID = iota + 1

	BuildingIdGenericFactory  geBuildingID = iota + 1
	BuildingIdGenericShipyard geBuildingID = iota + 1
	BuildingIdGenericLab      geBuildingID = iota + 1
	BuildingIdGenericCommand  geBuildingID = iota + 1
)

type BuildingInterface interface {
	GetCost(level int64) ResourcesStruct
	PrecalculateCost(levels int64)
}

type BuildingStruct struct {
	ID   geBuildingID
	Type geBuildingType

	BaseCost       ResourcesStruct
	BaseProduction float64
	Requirements   RequirementsStruct

	CostEquations       EquationStruct
	ProductionEquations EquationStruct

	costTable       []ResourcesStruct
	productionTable []int64
}

type BuildingLevelStruct struct {
	Building BuildingInterface
	Level    int64
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
}

func (b *BuildingStruct) GetCost(level int64) ResourcesStruct {
	return b.costTable[level]
}

func (b *BuildingStruct) GetProduction(level int64) int64 {
	return b.productionTable[level]
}
