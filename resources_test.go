package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func GetBaseResources() ResourcesStruct {
	return ResourcesStruct{
		Metal:   100,
		Silicon: 100,
		Uranium: 100,
		Energy:  100,
		Time:    10 * time.Second,
	}
}

func TestResourcesComparison(t *testing.T) {
	base := GetBaseResources()

	assert.True(t, base.HasEnoughBasic(ResourcesStruct{Metal: 50, Silicon: 50, Uranium: 50}))
	assert.False(t, base.HasEnoughBasic(ResourcesStruct{Metal: 150, Silicon: 50, Uranium: 50}))
	assert.True(t, base.HasEnoughBasic(ResourcesStruct{Metal: 50, Silicon: 50, Uranium: 50, Energy: 200, Time: 20 * time.Second}))

	assert.True(t, base.HasEnoughFull(ResourcesStruct{Metal: 50, Silicon: 50, Uranium: 50, Energy: 50, Time: 10 * time.Second}))
	assert.False(t, base.HasEnoughFull(ResourcesStruct{Metal: 50, Silicon: 50, Uranium: 50, Energy: 200, Time: 20 * time.Second}))
}

func TestResourcesMath(t *testing.T) {
	base := GetBaseResources()

	other := ResourcesStruct{
		Metal:   50,
		Silicon: 50,
		Uranium: 50,
		Energy:  50,
		Time:    10 * time.Second,
	}

	base.SubtractBasic(other)
	assert.EqualValues(t, ResourcesStruct{Metal: 50, Silicon: 50, Uranium: 50, Energy: 100, Time: 10 * time.Second}, base)

	base = GetBaseResources()
	base.SubtractFull(other)
	assert.EqualValues(t, ResourcesStruct{Metal: 50, Silicon: 50, Uranium: 50, Energy: 50, Time: 0 * time.Second}, base)

}
