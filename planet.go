package main

const (
	DefaultPlanetName string = "Planeta matka"
)

type CoordinatesStruct struct {
	Galaxy int64
	System int64
	Planet int64
}

type PlanetStruct struct {
	Position CoordinatesStruct
	Owner *PlayerStruct `json:"-"`
	Name string
	
	Research map[string]int64
	Buildings map[string]int64
}

func GenerateNewPlanet(universe *UniverseStruct, baseData *BaseDataStore) (*PlanetStruct, error){
	position, err := universe.GetNewCoordinates()
	if err != nil {
		return nil, err
	}
	
	planet := &PlanetStruct{}
	planet.Position = *position
	planet.Name = DefaultPlanetName
	planet.Buildings = make(map[string]int64)
	planet.Research = make(map[string]int64)
	
	universe.AddPlanet(position, planet)	
	
	return planet, nil
}
