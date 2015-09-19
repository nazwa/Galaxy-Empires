package main

type BuildingStruct struct {
	ID                  string
	Name                string
	Category            string
	Produces            string
	Description         string
	BaseCost            ResourcesStruct
	BaseProduction      int64
	Requirements        RequirementsStruct
	CostEquations       EquationStruct
	ProductionEquations EquationStruct
}
