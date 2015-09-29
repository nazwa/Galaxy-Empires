package ge

import (
	"time"
)

type BaseDataStore struct {
	Buildings map[geBuildingID]BuildingInterface
	Research  map[string]*ResearchStruct
}

func NewBaseDataStore(buildings, research string) *BaseDataStore {
	store := &BaseDataStore{}

	store.Buildings = make(map[geBuildingID]BuildingInterface) /*
		if err := LoadFile(buildings, &store.Buildings); err != nil {
			panic(err)
		}*/
	store.LoadStartingBuildings()
	store.Research = make(map[string]*ResearchStruct)
	if err := LoadFile(buildings, &store.Research); err != nil {
		panic(err)
	}

	store.PrecalculateCostsAndProduction()

	return store
}

func (b *BaseDataStore) PrecalculateCostsAndProduction() {
	for _, building := range b.Buildings {
		if build, ok := building.(BuildingInterface); ok {
			build.PrecalculateCostTable(PrecalculateLevels)
		}

		if mine, ok := building.(BuildingMineInterface); ok {
			mine.PrecalculateProductionTable(PrecalculateLevels)
		}
	}
}

func (b *BaseDataStore) GetStartingBuildings() map[geBuildingID]BuildingLevelStruct {
	list := make(map[geBuildingID]BuildingLevelStruct)
	for id, building := range b.Buildings {
		list[id] = BuildingLevelStruct{
			Building: building,
			Level:    0,
		}
	}

	return list
}

func (b *BaseDataStore) LoadStartingBuildings() {
	b.Buildings[BuildingIdMineMetal] = &BuildingMineStruct{
		BuildingStruct: BuildingStruct{
			ID:            BuildingIdMineMetal,
			BaseCost:      ResourcesStruct{30, 20, 0, 0, 10 * time.Second},
			CostEquations: EquationStruct{1, 1.5, 1},
		},
		productionEquations: EquationStruct{1, 1.1, 1},
	}
	b.Buildings[BuildingIdMineSilicon] = &BuildingMineStruct{
		BuildingStruct: BuildingStruct{
			ID:            BuildingIdMineSilicon,
			BaseCost:      ResourcesStruct{30, 20, 0, 0, 10 * time.Second},
			CostEquations: EquationStruct{1, 1.5, 1},
		},
		productionEquations: EquationStruct{1, 1.1, 1},
	}
	b.Buildings[BuildingIdMineUranium] = &BuildingMineStruct{
		BuildingStruct: BuildingStruct{
			ID:            BuildingIdMineUranium,
			BaseCost:      ResourcesStruct{30, 20, 0, 0, 10 * time.Second},
			CostEquations: EquationStruct{1, 1.5, 1},
		},
		productionEquations: EquationStruct{1, 1.1, 1},
	}
	b.Buildings[BuildingIdGenericFactory] = &BuildingStruct{
		ID:            BuildingIdGenericFactory,
		BaseCost:      ResourcesStruct{30, 20, 0, 0, 10 * time.Second},
		CostEquations: EquationStruct{1, 1.5, 1},
	}
}
