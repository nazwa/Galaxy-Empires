package main

import ()

type BaseDataStore struct {
	Buildings map[string]BuildingStruct
	Research  map[string]ResearchStruct
}

func NewBaseDataStore(buildings, research string) *BaseDataStore {
	store := &BaseDataStore{}

	store.Buildings = make(map[string]BuildingStruct)
	if err := LoadFile(buildings, &store.Buildings); err != nil {
		panic(err)
	}
	store.Research = make(map[string]ResearchStruct)
	if err := LoadFile(buildings, &store.Research); err != nil {
		panic(err)
	}

	return store
}

/*
func (s *BaseDataStore) LoadBuildings(file string) error {
	s.Buildings = make(map[int64]BuildingStruct)
	temp := make([]BuildingStruct, 10)

	if err := LoadFile(file, &temp); err != nil {
		return err
	}
	for _, item := range temp {
		if item.ID == 0 {
			return errors.New("Building ID is zero!")
		}
		s.Buildings[item.ID] = item
	}
	return nil
}

func (s *BaseDataStore) LoadResearch(file string) error {
	s.Research = make(map[int64]ResearchStruct)
	temp := make([]ResearchStruct, 10)

	if err := LoadFile(file, &temp); err != nil {
		return err
	}
	for _, item := range temp {
		if item.ID == 0 {
			return errors.New("Building ID is zero!")
		}
		s.Research[item.ID] = item
	}
	return nil
}
*/
