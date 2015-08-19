package main

type BuildingStruct struct {
	ID             string
	Name           string
	Category       string
	Description    string
	BaseCost       ResourcesStruct
	BaseProduction ResourcesStruct
	Requirements   RequirementsStruct
}
