package ge

import (
	"strconv"
)

type geBuildingType int64
type geBuildingID int64
type geBuildingsLevelMap map[geBuildingID]BuildingLevelStruct

// Converts string into geBuildingType
func BuildingTypeFromString(val string) (geBuildingType, error) {
	if buildingType, err := strconv.ParseInt(val, 10, 64); err == nil {
		return geBuildingType(buildingType), nil
	}
	return 0, ErrorInvalidBuildingType
}

// Converts string into geBuildingID
func BuildingIdFromString(val string) (geBuildingID, error) {
	if buildingId, err := strconv.ParseInt(val, 10, 64); err == nil {
		return geBuildingID(buildingId), nil
	}
	return 0, ErrorInvalidBuildingID
}

func (i *geBuildingID) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

// Simplifies the buildings list to only a map of levels for each building.
// Nothing else is needed
func (g *geBuildingsLevelMap) SimplifiedLevels() map[string]int64 {
	levels := make(map[string]int64)
	for id, data := range *g {
		levels[id.String()] = data.Level
	}
	return levels
}
