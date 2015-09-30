package ge

func (p *PlanetStruct) ToPublic(private bool) map[string]interface{} {
	data := make(map[string]interface{})

	data["Position"] = p.position
	data["Owner"] = p.owner.ID
	data["Name"] = p.name
	// Only show private data if planet is owned by current player
	if private {
		data["Resources"] = p.resources
		data["ResourcesHourly"] = p.hourlyResources
		data["Buildings"] = p.buildings.SimplifiedLevels()
		data["Research"] = p.research
		data["BuildingProgress"] = p.buildingInProgress
	}

	return data
}
