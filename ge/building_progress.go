package ge

import (
	"time"
)

type BuildingProgressStruct struct {
	Building  BuildingInterface
	StartTime time.Time
	EndTime   time.Time
	Cost      ResourcesStruct
}
