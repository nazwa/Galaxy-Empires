package ge

import (
	"errors"
	"github.com/gin-gonic/gin"
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
	PrecalculateLevels int64 = 50
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
// Default success response lol
var (
	DefaultSuccessResponse gin.H = gin.H{"Success": true}
)

// <----
// ERRORS
// ---->
var (
	ErrorPlayerDatabaseMissing error = errors.New("Player database is missing")
	ErrorInvalidPlanetID       error = errors.New("Invalid planet ID")
	ErrorInvalidBuildingType   error = errors.New("Invalid building type")
	ErrorInvalidBuildingID     error = errors.New("Invalid building ID")
	ErrorInvalidCredentials    error = errors.New("Invalid credentials")
	ErrorPlayerNotFound        error = errors.New("Player not found")
	ErrorPlayerIDInUse         error = errors.New("Player ID has already been used")
	ErrorEmailInUse            error = errors.New("This email has already been used")
	ErrorUniverseFull          error = errors.New("No spaces left in the universe")
	ErrorMineNotFound          error = errors.New("Not all mines are present")
	ErrorInsufficientResources error = errors.New("Insufficient resources")
	ErrorBuildingInProgress    error = errors.New("Another building in progress")
	ErrorResearchInProgress    error = errors.New("Another research in progress")
)
