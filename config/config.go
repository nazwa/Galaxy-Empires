package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config represents the configuration information.
type ConfigStruct struct {
	Debug    bool
	Port     string
	Key      string
	Services struct {
		Rollbar struct {
			Token       string
			Environment string
		}
	}
}

var Config ConfigStruct

func LoadConfig(filePath string) error {
	// Get the config file
	config_file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(config_file, &Config)
}
