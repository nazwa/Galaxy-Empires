package ge

func (p *PlanetStruct) ToPublic(private bool) map[string]interface{} {
	data := make(map[string]interface{})

	data["Position"] = p.Position
	data["Owner"] = p.Owner.ID
	data["Name"] = p.Name
	// Only show private data if planet is owned by current player
	if private {
		data["Resources"] = p.Resources
		data["ResourcesHourly"] = p.ResourcesHourly
		data["Buildings"] = p.Buildings.SimplifiedLevels()
		data["Research"] = p.Research
		data["BuildingProgress"] = p.BuildingInProgress
	}

	return data
}
