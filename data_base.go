package main

import (
	"errors"
	"strings"
)

type BaseDataStore struct {
	Buildings map[string]*BuildingStruct
	Research  map[string]*ResearchStruct
}

var (
	MetalMineKey   string
	SiliconMineKey string
	UraniumMineKey string
	PowerPlantKey  string

	ErrorMineNotFound error = errors.New("Not all mines are present")
)

const (
	MineCategoryString     string = "Mine"
	MetalMineCategoryKey   string = "Metal"
	SiliconMineCategoryKey string = "Silicon"
	UraniumMineCategoryKey string = "Uranium"
	PowerPlantCategoryKey  string = "Energy"
)

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
