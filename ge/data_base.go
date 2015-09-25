package ge

import (
	"strings"
)

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

	store.ExtractMineKeys()
	store.PrecalculateCostsAndProduction()

	return store
}

func (b *BaseDataStore) FindMineByType(produces string) string {
	for key, building := range b.Buildings {
		if strings.EqualFold(building.Category, MineCategoryString) &&
			strings.EqualFold(building.Produces, produces) {
			return key
		}
	}
	panic(ErrorMineNotFound)
}

func (b *BaseDataStore) ExtractMineKeys() {
	MetalMineKey = b.FindMineByType(MetalMineCategoryKey)
	SiliconMineKey = b.FindMineByType(SiliconMineCategoryKey)
	UraniumMineKey = b.FindMineByType(UraniumMineCategoryKey)
	PowerPlantKey = b.FindMineByType(PowerPlantCategoryKey)
}

func (b *BaseDataStore) PrecalculateCostsAndProduction() {
	for _, building := range b.Buildings {
		building.PrecalculateCostTable(PrecalculateLevels)
		building.PrecalculateProductionTable(PrecalculateLevels)
	}
}
