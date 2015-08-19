package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func LoadFile(file string, target interface{}) error {
	path := filepath.Join(ROOT_DIR, file)

	// Get the config file
	config_file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(config_file, target)
	if err != nil {
		return err
	}
	return nil
}
