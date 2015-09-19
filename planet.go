package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	DefaultPlanetName   string        = "Planeta matka"
	MinimumResourceTime time.Duration = 1 * time.Second
)

type CoordinatesStruct struct {
	Galaxy int64
	System int64
	Planet int64
}

type PlanetStruct struct {
	Position CoordinatesStruct
	Owner    *PlayerStruct `json:"-"`

	Resources           ResourcesStruct
	ResourcesUpdateTime time.Time
	ResourcesMutex      sync.Mutex

	Name string

	Research  map[string]int64
	Buildings map[string]int64
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
	planet.Buildings[MetalMineKey] = 25
	planet.Buildings[SiliconMineKey] = 0
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

	universe.AddPlanet(position, planet)

	return planet, nil
}

func (p *PlanetStruct) RecalculateResources(baseData *BaseDataStore) {
	p.ResourcesMutex.Lock()
	defer p.ResourcesMutex.Unlock()

	now := time.Now()
	timeDiff := now.Sub(p.ResourcesUpdateTime)
	// We dont want to udpate resources more often than this time
	if timeDiff < MinimumResourceTime {
		return
	}

	p.Resources.Metal += p.CalculateProduction(baseData.Buildings[MetalMineKey], p.Buildings[MetalMineKey], timeDiff)
	p.Resources.Silicon += p.CalculateProduction(baseData.Buildings[SiliconMineKey], p.Buildings[SiliconMineKey], timeDiff)
	p.Resources.Uranium += p.CalculateProduction(baseData.Buildings[UraniumMineKey], p.Buildings[UraniumMineKey], timeDiff)

	p.ResourcesUpdateTime = now
}

func (p *PlanetStruct) CalculateProduction(building *BuildingStruct, level int64, timeDiff time.Duration) float64 {
	baseProduction := building.BaseProduction * building.ProductionEquations.A * float64(level) * building.ProductionEquations.B
	production := baseProduction * math.Pow(building.ProductionEquations.C, float64(level))

	fmt.Println(baseProduction)
	fmt.Println(production)
	fmt.Println((production / 3600) * timeDiff.Seconds())
	return (production / 3600) * timeDiff.Seconds()
}
