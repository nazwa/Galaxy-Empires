package main

import (
	"time"
)

type ResearchProgressStruct struct {
	Research  *ResearchStruct
	StartTime time.Time
	EndTime   time.Time
	Cost      ResourcesStruct
}
