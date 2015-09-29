package ge

import ()

type BaseDataStore struct {
	Buildings map[string]*BuildingStruct
	Research  map[string]*ResearchStruct
}

func NewBaseDataStore(buildings, research string) *BaseDataStore {
	store := &BaseDataStore{}

	store.Buildings = make(map[string]*BuildingStruct)
	if err := LoadFile(buildings, &store.Buildings); err != nil {
		panic(err)
	}
	store.Research = make(map[string]*ResearchStruct)
	if err := LoadFile(buildings, &store.Research); err != nil {
		panic(err)
	}

	store.PrecalculateCostsAndProduction()

	return store
}

func (b *BaseDataStore) PrecalculateCostsAndProduction() {
	for _, building := range b.Buildings {
		building.PrecalculateCostTable(PrecalculateLevels)
		building.PrecalculateProductionTable(PrecalculateLevels)
	}
}
