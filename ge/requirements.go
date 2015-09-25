package main

type BuildingRequirement struct {
	ID    string
	Level int64
}

type ResearchRequirement struct {
	ID    string
	Level int64
}

type RequirementsStruct struct {
	Buildings []BuildingRequirement
	Research  []ResearchRequirement
}
