package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestUniverseCreation(t *testing.T) {
	// Create universe with one galaxy, 2 systems and 2 planets each
	universe := NewUniverseStruct(1, 2, 2)
		
	// First system should be the emptiest
	assert.EqualValues(t, CoordinatesStruct{0,0,0}, universe.GetEmptiestSystem().coordinates)
	
	// LEt's add a new planet. Second system should now be the emptiest
	universe.AddPlanet(&CoordinatesStruct{0,0,0}, &PlanetStruct{})
	assert.EqualValues(t, CoordinatesStruct{0,1,0}, universe.GetEmptiestSystem().coordinates)
	
	// LEt's add another planet. First system should now be the emptiest
	universe.AddPlanet(&CoordinatesStruct{0,1,0}, &PlanetStruct{})
	assert.EqualValues(t, CoordinatesStruct{0,0,0}, universe.GetEmptiestSystem().coordinates)
	
	// Let's get first available coordinates. No error
	coord, err := universe.GetNewCoordinates()
	assert.Nil(t, err)
	assert.EqualValues(t, &CoordinatesStruct{0,0,1}, coord)
	
	// Add planet at last coords and get a new one. No error
	universe.AddPlanet(coord, &PlanetStruct{})
	coord2, err := universe.GetNewCoordinates()
	assert.Nil(t, err)
	assert.EqualValues(t, &CoordinatesStruct{0,1,1}, coord2)

	// Add planet at last coords and get a new one. 
	// We should now get Universe full error
	universe.AddPlanet(coord2, &PlanetStruct{})
	coord3, err := universe.GetNewCoordinates()
	assert.Error(t, err)
	assert.Nil(t, coord3)
	
}

func TestSystemOperations(t *testing.T) {	
	// Create universe with one galaxy, 2 systems and 2 planets each
	universe := NewUniverseStruct(1, 1, 1)
	
	planet := &PlanetStruct{
		Name: "TestPlanet",
	}
	
	// There should be no planets
	assert.EqualValues(t, 0, universe.galaxies[0][0].count)
	
	// Add one and check the counter
	universe.AddPlanet(&CoordinatesStruct{0,0,0}, planet)
	assert.EqualValues(t, 1, universe.galaxies[0][0].count)
	
	// The planet we get should be exactly the one we put in
	planet2 := universe.GetPlanet(&CoordinatesStruct{0,0,0})
	assert.Equal(t, planet, planet2)
}