package ge

import ()

type GalaxyEmpires struct {
	BaseData   *BaseDataStore
	PlayerData *PlayerDataStore
	Universe   *UniverseStruct
}

func NewGalaxyEmpires(dataDir string, universeSize CoordinatesStruct) *GalaxyEmpires {
	g := &GalaxyEmpires{}

	g.PlayerData = NewPlayerDataStore(
		BuildFullPath(dataDir, "players.json"),
	)
	g.BaseData = NewBaseDataStore(
		BuildFullPath(dataDir, "buildings.json"),
		BuildFullPath(dataDir, "research.json"),
	)
	g.Universe = NewUniverseStruct(universeSize)

	return g

}
