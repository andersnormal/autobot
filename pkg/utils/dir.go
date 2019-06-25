package utils

import (
	"os"
)

// CreateDirIfNotExist ...
func CreateDirIfNotExist(dir string, chmod os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, chmod); err != nil {
			return err
		}
	}

	return nil
}
