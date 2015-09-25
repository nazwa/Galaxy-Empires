package ge

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func LoadFile(file string, target interface{}) error {
	// Get the config file
	data_file, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data_file, target)
	if err != nil {
		return err
	}
	return nil
}

func BuildFullPath(root, file string) string {
	return filepath.Join(root, file)
}
