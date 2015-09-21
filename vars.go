package main

import (
	"errors"
	"time"
)

const (
	// <----
	// Context keys
	// ---->
	PlanetObjectKey string = "PlanetObject"
	PlayerObjectKey string = "PlayerObject"

	// <----
	// Planets
	// ---->
	DefaultPlanetName   string        = "Super planeta"
	MinimumResourceTime time.Duration = 1 * time.Second

	// <----
	// Buildings
	// ---->
	MineCategoryString     string = "Mine"
	MetalMineCategoryKey   string = "Metal"
	SiliconMineCategoryKey string = "Silicon"
	UraniumMineCategoryKey string = "Uranium"
	PowerPlantCategoryKey  string = "Energy"
	PrecalculateLevels     int64  = 50
)

// <----
// Dynamic mine IDs (set on load)
// ---->
var (
	MetalMineKey   string
	SiliconMineKey string
	UraniumMineKey string
	PowerPlantKey  string
)

// <----
// ERRORS
// ---->
var (
	ErrorPlayerDatabaseMissing error = errors.New("Player database is missing")
	ErrorInvalidPlanetID       error = errors.New("Invalid planet ID")
	ErrorInvalidBuildingID     error = errors.New("Invalid building ID")
	ErrorInvalidCredentials    error = errors.New("Invalid credentials")
	ErrorPlayerNotFound        error = errors.New("Player not found")
	ErrorEmailInUse            error = errors.New("This email has already been used")
	ErrorUniverseFull          error = errors.New("No spaces left in the universe")
	ErrorMineNotFound          error = errors.New("Not all mines are present")
	ErrorInsufficientResources error = errors.New("Insufficient resources")
	ErrorBuildingInProgress    error = errors.New("Another building in progress")
)
