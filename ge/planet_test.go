package ge

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func MockBuildingMine(id geBuildingID) *BuildingMineStruct {
	return &BuildingMineStruct{
		BuildingStruct: BuildingStruct{
			ID:            id,
			BaseCost:      ResourcesStruct{30, 20, 0, 0, 10 * time.Second},
			CostEquations: EquationStruct{1, 1.5, 1},
		},
		productionEquations: EquationStruct{1, 1.1, 1},
		baseProduction:      100,
	}
}

func TestPlanetResources(t *testing.T) {

	// Create an empty planet
	planet := &PlanetStruct{}
	planet.buildings = make(geBuildingsLevelMap)

	// Check basic resources maths
	// We should start with empty resources
	assert.EqualValues(t, ResourcesStruct{0, 0, 0, 0, 0}, planet.hourlyResources)
	toAdd := ResourcesStruct{100, 100, 100, 0, 0}
	// Let's add some small values
	planet.AddResources(toAdd)
	assert.EqualValues(t, toAdd, planet.resources)
	// Then remove the same amount and we should be back to zero
	err := planet.SubtractResources(toAdd)
	assert.NoError(t, err)
	assert.EqualValues(t, ResourcesStruct{0, 0, 0, 0, 0}, planet.hourlyResources)
	// LEt's try to subtract again. should error out and stay on zero
	err = planet.SubtractResources(toAdd)
	assert.Error(t, err)
	assert.EqualValues(t, ResourcesStruct{0, 0, 0, 0, 0}, planet.hourlyResources)

	// Let's add a metal mine
	mine := MockBuildingMine(BuildingIdMineMetal)
	mine.PrecalculateProductionTable(10)
	// Just make sure we have some values
	assert.NotEmpty(t, mine.productionTable)

	level := &BuildingLevelStruct{
		Building: mine,
		Level:    1,
	}
	planet.buildings[BuildingIdMineMetal] = level
	// We have a first level mine, let's updte productions and see if we got the value
	planet.UpdateHourlyProduction()
	assert.EqualValues(t, mine.GetProduction(level.Level), planet.hourlyResources.Metal)

	// Lets see if our hourly production calculations are fine
	assert.EqualValues(t, 100, planet.CalculateProductionSince(100, 1*time.Hour))
	assert.EqualValues(t, 50, planet.CalculateProductionSince(100, 30*time.Minute))
	assert.EqualValues(t, 150, planet.CalculateProductionSince(100, 90*time.Minute))

	planet.lastResourcesUpdateTime = time.Now()
	// First try with 'instant' update. Nothing should change
	planet.RecalculateResources(planet.lastResourcesUpdateTime)
	assert.Equal(t, ResourcesStruct{0, 0, 0, 0, 0}, planet.resources)
	// Now add one hour
	planet.RecalculateResources(planet.lastResourcesUpdateTime.Add(60 * time.Minute))
	assert.Equal(t, ResourcesStruct{float64(mine.GetProduction(1)), 0, 0, 0, 0}, planet.resources)

}

func TestBuildings(t *testing.T) {

	// Create an empty planet
	planet := &PlanetStruct{}
	planet.buildings = make(geBuildingsLevelMap)

	mine := MockBuildingMine(BuildingIdMineMetal)
	level := &BuildingLevelStruct{
		Building: mine,
		Level:    1,
	}

	planet.buildings[BuildingIdMineMetal] = level
	assert.NotEmpty(t, planet.buildings)

	assert.Equal(t, level, planet.GetBuilding(BuildingIdMineMetal))

	//planet.BuildBuilding()

}

func TestPlanetGettersSetters(t *testing.T) {
	planet := &PlanetStruct{}

	assert.Empty(t, planet.name)
	planet.SetName("hello")
	assert.Equal(t, "hello", planet.name)
}
