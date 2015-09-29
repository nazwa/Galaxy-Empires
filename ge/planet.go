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
	Position CoordinatesStruct
	Owner    *PlayerStruct

	Resources           ResourcesStruct
	ResourcesHourly     ResourcesStruct
	ResourcesUpdateTime time.Time
	ResourcesMutex      sync.Mutex

	Name string

	Research  map[string]int64
	Buildings geBuildingsLevelMap

	BuildingInProgress      *BuildingProgressStruct
	BuildingInProgressMutex sync.Mutex

	ResearchInProgress      *ResearchProgressStruct
	ResearchInProgressMutex sync.Mutex
}

func GenerateNewPlanet(universe *UniverseStruct, baseData *BaseDataStore) (*PlanetStruct, error) {
	position, err := universe.GetNewCoordinates()
	if err != nil {
		return nil, err
	}

	planet := &PlanetStruct{}
	planet.Position = *position
	planet.Name = DefaultPlanetName
	planet.Buildings = baseData.GetStartingBuildings()

	// Set basic mine levels
	planet.Research = make(map[string]int64)
	planet.Resources = ResourcesStruct{
		Metal:   1000,
		Silicon: 1000,
		Uranium: 0,
		Energy:  0,
	}
	planet.ResourcesUpdateTime = time.Now()
	planet.UpdateHourlyProduction()

	universe.AddPlanet(position, planet)

	return planet, nil
}

func (p *PlanetStruct) UpdatePlanet(baseData *BaseDataStore, now time.Time) {
	// Make sure no mines have been built while we were away
	// This will recalculate resources up to the time of the construction end
	p.UpdateConstruction(baseData, now)

	// Recalculate resources for real now
	p.RecalculateResources(baseData, now)
}

func (p *PlanetStruct) RecalculateResources(baseData *BaseDataStore, now time.Time) {
	p.ResourcesMutex.Lock()
	defer p.ResourcesMutex.Unlock()

	timeDiff := now.Sub(p.ResourcesUpdateTime)
	// We dont want to udpate resources more often than this time
	if timeDiff < MinimumResourceTime {
		return
	}

	p.Resources.Metal += p.CalculateProduction(p.ResourcesHourly.Metal, timeDiff)
	p.Resources.Silicon += p.CalculateProduction(p.ResourcesHourly.Silicon, timeDiff)
	p.Resources.Uranium += p.CalculateProduction(p.ResourcesHourly.Uranium, timeDiff)

	p.ResourcesUpdateTime = now
}

func (p *PlanetStruct) CalculateProduction(production float64, timeDiff time.Duration) float64 {
	return production / 3600 * timeDiff.Seconds()
}

func (p *PlanetStruct) UpdateHourlyProduction() {
	p.ResourcesHourly = ResourcesStruct{
		Metal:   float64(p.Buildings[BuildingIdMineMetal].Building.(BuildingMineInterface).GetProduction(p.Buildings[BuildingIdMineMetal].Level)),
		Silicon: float64(p.Buildings[BuildingIdMineSilicon].Building.(BuildingMineInterface).GetProduction(p.Buildings[BuildingIdMineSilicon].Level)),
		Uranium: float64(p.Buildings[BuildingIdMineUranium].Building.(BuildingMineInterface).GetProduction(p.Buildings[BuildingIdMineUranium].Level)),
	}
}

func (p *PlanetStruct) SubtractResources(resources ResourcesStruct) error {
	p.ResourcesMutex.Lock()
	defer p.ResourcesMutex.Unlock()

	if !p.Resources.HasEnoughBasic(resources) {
		return ErrorInsufficientResources
	}
	p.Resources.SubtractBasic(resources)

	return nil
}

//Adding should always be successful, so no need for errors! yay
func (p *PlanetStruct) AddResources(resources ResourcesStruct) {
	p.ResourcesMutex.Lock()
	defer p.ResourcesMutex.Unlock()

	p.Resources.AddBasic(resources)
}

// Checks if there are any buildings that have finished
// Recalculates resources before completition
// @todo: FIRE EVENT TO NOTIFY BUILDING COMPLETITION
func (p *PlanetStruct) UpdateConstruction(baseData *BaseDataStore, now time.Time) {
	p.BuildingInProgressMutex.Lock()
	defer p.BuildingInProgressMutex.Unlock()

	if p.BuildingInProgress == nil {
		return
	}

	// Triggers if building end time has already passed
	if p.BuildingInProgress.EndTime.Sub(now) < 0 {
		// Recalculate the resources with old levels
		p.RecalculateResources(baseData, p.BuildingInProgress.EndTime)

		// Finish the actual build
		level := p.Buildings[p.BuildingInProgress.Building.(*BuildingStruct).ID]
		level.Level++

		p.BuildingInProgress = nil
		p.UpdateHourlyProduction()
	}
}

func (p *PlanetStruct) BuildBuilding(id geBuildingID) error {
	p.BuildingInProgressMutex.Lock()
	defer p.BuildingInProgressMutex.Unlock()

	if p.BuildingInProgress != nil {
		return ErrorBuildingInProgress
	}

	var current BuildingLevelStruct
	var ok bool

	if current, ok = p.Buildings[id]; !ok {
		return ErrorInvalidBuildingID
	}

	toLevel := current.Level + 1

	cost := current.Building.(*BuildingStruct).GetCost(toLevel)

	if err := p.SubtractResources(cost); err != nil {
		return err
	}

	p.BuildingInProgress = &BuildingProgressStruct{
		Building:  current.Building,
		Cost:      cost,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(cost.Time),
	}

	return nil
}

func (p *PlanetStruct) CancelBuilding() {
	p.BuildingInProgressMutex.Lock()
	defer p.BuildingInProgressMutex.Unlock()

	if p.BuildingInProgress == nil {
		return
	}

	p.AddResources(p.BuildingInProgress.Cost)
	p.BuildingInProgress = nil

	return
}

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
}
