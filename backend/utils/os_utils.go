package utils

import (
	"encoding/json"
	"os"
)

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || (!os.IsNotExist(err))
}

func SymlinkExists(name string) bool {
	_, err := os.Lstat(name)
	return err == nil || (!os.IsNotExist(err))
}

func ReadJsonFile[T any](filename string, target T) error {
	if content, err := os.ReadFile(filename); err != nil {
		return err
	} else {
		return json.Unmarshal(content, target)
	}
}
