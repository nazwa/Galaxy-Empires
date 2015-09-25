package ge

import (
	"fmt"
	"math"
	"time"
)

type ResearchStruct struct {
	ID          string
	Name        string
	Category    string
	Description string
	BaseCost    ResourcesStruct

	CostEquations EquationStruct
	//BuildEqations       EquationStruct

	costTable []ResourcesStruct
}

func (r *ResearchStruct) PrecalculateCostTable(limit int64) {
	r.costTable = make([]ResourcesStruct, limit)

	for i := int64(0); i < limit; i++ {
		factor := r.CostEquations.A * math.Pow(r.CostEquations.B, float64(i-1)*r.CostEquations.C)
		r.costTable[i] = ResourcesStruct{
			Metal:   r.BaseCost.Metal * factor,
			Silicon: r.BaseCost.Silicon * factor,
			Uranium: r.BaseCost.Uranium * factor,
			Energy:  r.BaseCost.Energy * factor,
			Time:    time.Duration(float64(r.BaseCost.Time*time.Second) * factor),
		}
	}
	fmt.Println(r.Name)
	fmt.Println(r.costTable[1])
	fmt.Println(r.costTable[2])
	fmt.Println("--")
	time.Now()
}

func (r *ResearchStruct) GetCost(level int64) ResourcesStruct {
	return r.costTable[level]
}
