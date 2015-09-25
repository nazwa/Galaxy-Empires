package main

import (
	"time"
)

type BuildingProgressStruct struct {
	Building  *BuildingStruct
	StartTime time.Time
	EndTime   time.Time
	Cost      ResourcesStruct
}
