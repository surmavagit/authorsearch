package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// fileNotExist returns true and error if the cache file with the specified
// name doesn't exist. It returns false if either there is no error or
// if there is a different error
func fileNotExist(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	return errors.Is(err, os.ErrNotExist), err
}

// loadFileJSON reads the cache file and unmarshals the json into the target
func loadFileJSON(fileName string, target interface{}) error {
	stream, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(stream, target)
	return err
}

// writeFileJSON turns source into indented JSON and writes it into file
func writeFileJSON(fileName string, source any) error {
	stream, err := json.MarshalIndent(source, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, []byte(stream), 0644)
}

func createDirIfNotExist(directory string) error {
	info, err := os.Stat(directory)

	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(directory, 0755)
		if err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", directory)
	}

	return nil
}
