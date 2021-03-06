package ge

import (
	"bitbucket.org/tidepayments/gohelpers/maths"
)

type UniverseStruct struct {
	galaxies   [][]GalaxySystemStruct
	size       CoordinatesStruct
	lastGalaxy int64
}

type GalaxySystemStruct struct {
	planets     []*PlanetStruct
	coordinates CoordinatesStruct
	count       int64
}

func NewUniverseStruct(universeSize CoordinatesStruct) *UniverseStruct {
	universe := &UniverseStruct{}
	universe.size = CoordinatesStruct{
		Galaxy: universeSize.Galaxy,
		System: universeSize.System,
		Planet: universeSize.Planet,
	}

	universe.galaxies = make([][]GalaxySystemStruct, universeSize.Galaxy)
	for i := range universe.galaxies {
		universe.galaxies[i] = make([]GalaxySystemStruct, universeSize.System)
		for j := range universe.galaxies[i] {
			universe.galaxies[i][j].planets = make([]*PlanetStruct, universeSize.Planet)
			universe.galaxies[i][j].count = 0
			universe.galaxies[i][j].coordinates = CoordinatesStruct{
				Galaxy: int64(i),
				System: int64(j),
			}
		}
	}

	return universe
}

func (g *UniverseStruct) GetSize() CoordinatesStruct {
	return g.size
}

func (g *UniverseStruct) GetSystem(coord *CoordinatesStruct) *GalaxySystemStruct {
	return &g.galaxies[coord.Galaxy][coord.System]
}

func (g *UniverseStruct) GetPlanet(coord *CoordinatesStruct) *PlanetStruct {
	return g.galaxies[coord.Galaxy][coord.System].planets[coord.Planet]
}

func (g *UniverseStruct) AddPlanet(coord *CoordinatesStruct, newPlanet *PlanetStruct) {
	g.galaxies[coord.Galaxy][coord.System].planets[coord.Planet] = newPlanet
	g.galaxies[coord.Galaxy][coord.System].count++
}

func (g *UniverseStruct) GetEmptiestSystem() *GalaxySystemStruct {
	min := g.galaxies[0][0]

	for i := range g.galaxies {
		for j := range g.galaxies[i] {
			if g.galaxies[i][j].count < min.count {
				min = g.galaxies[i][j]
			}
		}
	}
	return &min
}

func (g *UniverseStruct) GetNewCoordinates() (*CoordinatesStruct, error) {
	system := g.GetEmptiestSystem()

	positions := make([]int64, 0, 0)

	for i, planet := range system.planets {
		if planet == nil {
			positions = append(positions, int64(i))
		}
	}

	if len(positions) == 0 {
		return nil, ErrorUniverseFull
	}
	index := maths.RandomBetweenInt(0, len(positions)-1)

	return &CoordinatesStruct{
		Galaxy: system.coordinates.Galaxy,
		System: system.coordinates.System,
		Planet: positions[index],
	}, nil
}
