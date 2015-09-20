package main

func (p *PlayerStruct) ToPublic(private bool) map[string]interface{} {
	data := make(map[string]interface{})

	data["Name"] = p.Name

	if private {
		data["Email"] = p.Email
		p.planetMutex.Lock()
		defer p.planetMutex.Unlock()
		planets := make([]map[string]interface{}, len(p.Planets))
		for i := range p.Planets {
			planets[i] = p.Planets[i].ToPublic(private)
		}
		data["Planets"] = planets
	}

	return data
}
