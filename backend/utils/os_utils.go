package utils

import "os"

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || (!os.IsNotExist(err))
}

func SymlinkExists(name string) bool {
	_, err := os.Lstat(name)
	return err == nil || (!os.IsNotExist(err))
}
