package ge

import (
	"sync"
	"time"
)

type CoordinatesStruct struct {
	Galaxy int64
	System int64
	Planet int64
}

type PlanetStruct struct {
	position CoordinatesStruct
	owner    *PlayerStruct

	resourcesMutex          sync.Mutex
	resources               ResourcesStruct
	hourlyResources         ResourcesStruct
	lastResourcesUpdateTime time.Time

	name string

	research  map[string]int64
	buildings geBuildingsLevelMap

	buildingInProgress      *BuildingProgressStruct
	buildingInProgressMutex sync.Mutex

	researchInProgress      *ResearchProgressStruct
	researchInProgressMutex sync.Mutex
}

func GenerateNewPlanet(universe *UniverseStruct, baseData *BaseDataStore) (*PlanetStruct, error) {
	position, err := universe.GetNewCoordinates()
	if err != nil {
		return nil, err
	}

	planet := &PlanetStruct{}
	planet.position = *position
	planet.name = DefaultPlanetName
	planet.buildings = baseData.GetStartingBuildings()

	// Set basic mine levels
	planet.research = make(map[string]int64)
	planet.resources = ResourcesStruct{
		Metal:   1000,
		Silicon: 1000,
		Uranium: 0,
		Energy:  0,
	}
	planet.lastResourcesUpdateTime = time.Now()
	planet.UpdateHourlyProduction()

	universe.AddPlanet(position, planet)

	return planet, nil
}

func (p *PlanetStruct) GetBuilding(id geBuildingID) *BuildingLevelStruct {
	if b, ok := p.buildings[id]; ok {
		return b
	}
	return nil

}

func (p *PlanetStruct) SetName(name string) {
	p.name = name
}

func (p *PlanetStruct) UpdatePlanet(now time.Time) {
	// Make sure no mines have been built while we were away
	// This will recalculate resources up to the time of the construction end
	p.UpdateConstruction(now)

	// Recalculate resources for real now
	p.RecalculateResources(now)
}

func (p *PlanetStruct) RecalculateResources(now time.Time) {
	p.resourcesMutex.Lock()
	defer p.resourcesMutex.Unlock()

	timeDiff := now.Sub(p.lastResourcesUpdateTime)
	// We dont want to udpate resources more often than this time
	if timeDiff < MinimumResourceTime {
		return
	}

	p.resources.Metal += p.CalculateProductionSince(p.hourlyResources.Metal, timeDiff)
	p.resources.Silicon += p.CalculateProductionSince(p.hourlyResources.Silicon, timeDiff)
	p.resources.Uranium += p.CalculateProductionSince(p.hourlyResources.Uranium, timeDiff)

	p.lastResourcesUpdateTime = now
}

// Calculates production of the resource during given time duration
func (p *PlanetStruct) CalculateProductionSince(production float64, timeDiff time.Duration) float64 {
	return production / 3600 * timeDiff.Seconds()
}

// Updates cached production values for current mine levels
func (p *PlanetStruct) UpdateHourlyProduction() {
	p.hourlyResources = ResourcesStruct{
		Metal:   p.GetBuilding(BuildingIdMineMetal).GetBuildingProduction(),
		Silicon: p.GetBuilding(BuildingIdMineSilicon).GetBuildingProduction(),
		Uranium: p.GetBuilding(BuildingIdMineUranium).GetBuildingProduction(),
	}
}

func (p *PlanetStruct) SubtractResources(resources ResourcesStruct) error {
	p.resourcesMutex.Lock()
	defer p.resourcesMutex.Unlock()

	if !p.resources.HasEnoughBasic(resources) {
		return ErrorInsufficientResources
	}
	p.resources.SubtractBasic(resources)

	return nil
}

//Adding should always be successful, so no need for errors! yay
func (p *PlanetStruct) AddResources(resources ResourcesStruct) {
	p.resourcesMutex.Lock()
	defer p.resourcesMutex.Unlock()

	p.resources.AddBasic(resources)
}

// Checks if there are any buildings that have finished
// Recalculates resources before completition
// @todo: FIRE EVENT TO NOTIFY BUILDING COMPLETITION
func (p *PlanetStruct) UpdateConstruction(now time.Time) {
	p.buildingInProgressMutex.Lock()
	defer p.buildingInProgressMutex.Unlock()

	if p.buildingInProgress == nil {
		return
	}

	// Triggers if building end time has already passed
	if p.buildingInProgress.EndTime.Sub(now) < 0 {
		// Recalculate the resources with old levels
		p.RecalculateResources(p.buildingInProgress.EndTime)

		// Finish the actual build
		// This should NEVER EVER panic, but including it as a failsafe
		building := p.GetBuilding(p.buildingInProgress.Building.(BuildingInterface).GetId())
		if building == nil {
			panic(ErrorInvalidBuildingID)
		}
		building.Level++

		p.buildingInProgress = nil
		p.UpdateHourlyProduction()
	}
}

func (p *PlanetStruct) BuildBuilding(id geBuildingID) error {
	p.buildingInProgressMutex.Lock()
	defer p.buildingInProgressMutex.Unlock()

	if p.buildingInProgress != nil {
		return ErrorBuildingInProgress
	}

	current := p.GetBuilding(id)
	if current == nil {
		return ErrorInvalidBuildingID
	}

	toLevel := current.Level + 1

	cost := current.Building.(BuildingInterface).GetCost(toLevel)

	if err := p.SubtractResources(cost); err != nil {
		return err
	}

	p.buildingInProgress = &BuildingProgressStruct{
		Building:  current.Building,
		Cost:      cost,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(cost.Time),
	}

	return nil
}

func (p *PlanetStruct) CancelBuilding() {
	p.buildingInProgressMutex.Lock()
	defer p.buildingInProgressMutex.Unlock()

	if p.buildingInProgress == nil {
		return
	}

	p.AddResources(p.buildingInProgress.Cost)
	p.buildingInProgress = nil

	return
}

/*
func (p *PlanetStruct) BuildResearch(baseData *BaseDataStore, id string) error {
	p.ResearchInProgressMutex.Lock()
	defer p.ResearchInProgressMutex.Unlock()

	if p.ResearchInProgress != nil {
		return ErrorResearchInProgress
	}

	var research *ResearchStruct
	var toLevel int64
	var ok bool

	if research, ok = baseData.Research[id]; !ok {
		return ErrorInvalidBuildingID
	}

	if toLevel, ok = p.Research[id]; !ok {
		toLevel = 1
	} else {
		toLevel++
	}

	cost := research.GetCost(toLevel)

	if err := p.SubtractResources(cost); err != nil {
		return err
	}

	p.ResearchInProgress = &ResearchProgressStruct{
		Research:  research,
		Cost:      cost,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(cost.Time),
	}

	return nil
}

func (p *PlanetStruct) CancelResearch() {
	p.ResearchInProgressMutex.Lock()
	defer p.ResearchInProgressMutex.Unlock()

	if p.ResearchInProgress == nil {
		return
	}

	p.AddResources(p.ResearchInProgress.Cost)
	p.ResearchInProgress = nil

	return
}*/
