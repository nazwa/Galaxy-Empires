package main

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
	Buildings map[string]int64

	BuildingInProgress      *BuildingProgressStruct
	BuildingInProgressMutex sync.Mutex
}

func GenerateNewPlanet(universe *UniverseStruct, baseData *BaseDataStore) (*PlanetStruct, error) {
	position, err := universe.GetNewCoordinates()
	if err != nil {
		return nil, err
	}

	planet := &PlanetStruct{}
	planet.Position = *position
	planet.Name = DefaultPlanetName
	planet.Buildings = make(map[string]int64)

	// Set basic mine levels
	planet.Buildings[MetalMineKey] = 2
	planet.Buildings[SiliconMineKey] = 1
	planet.Buildings[UraniumMineKey] = 0
	planet.Buildings[PowerPlantKey] = 0

	// Set basic mine levels
	planet.Research = make(map[string]int64)
	planet.Resources = ResourcesStruct{
		Metal:   1000,
		Silicon: 1000,
		Uranium: 0,
		Energy:  0,
	}
	planet.ResourcesUpdateTime = time.Now()
	planet.UpdateHourlyProduction(baseData)

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

func (p *PlanetStruct) UpdateHourlyProduction(baseData *BaseDataStore) {
	p.ResourcesHourly = ResourcesStruct{
		Metal:   float64(baseData.Buildings[MetalMineKey].GetProduction(p.Buildings[MetalMineKey])),
		Silicon: float64(baseData.Buildings[SiliconMineKey].GetProduction(p.Buildings[SiliconMineKey])),
		Uranium: float64(baseData.Buildings[UraniumMineKey].GetProduction(p.Buildings[UraniumMineKey])),
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
		p.Buildings[p.BuildingInProgress.Building.ID]++
		p.BuildingInProgress = nil
		p.UpdateHourlyProduction(baseData)
	}
}

func (p *PlanetStruct) BuildBuilding(baseData *BaseDataStore, id string) error {
	p.BuildingInProgressMutex.Lock()
	defer p.BuildingInProgressMutex.Unlock()

	if p.BuildingInProgress != nil {
		return ErrorBuildingInProgress
	}

	var building *BuildingStruct
	var toLevel int64
	var ok bool

	if building, ok = baseData.Buildings[id]; !ok {
		return ErrorInvalidBuildingID
	}

	if toLevel, ok = p.Buildings[id]; !ok {
		toLevel = 1
	} else {
		toLevel++
	}

	cost := building.GetCost(toLevel)

	if err := p.SubtractResources(cost); err != nil {
		return err
	}

	p.BuildingInProgress = &BuildingProgressStruct{
		Building:  building,
		Cost:      cost,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(cost.Time),
	}

	return nil

}
