package cache

import (
	"log"
	"os"
)

// Returns an application specific temporary directory
func GetTempDir() string {
	dir := os.TempDir() + "/JODA-Web"
	err := ensureDir(dir)
	if err != nil {
		log.Printf("Could not create temporary dir '%s': %v", dir, err)
		return ""
	}
	return dir
}

// Returns an application specific user-cache directory.
// If no user specific cache directory can be found, falls back to GetTempDir()
func GetCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return GetTempDir()
	}
	dir = dir + "/JODA-Web"
	err = ensureDir(dir)
	if err != nil {
		log.Printf("Could not create cache dir '%s': %v", dir, err)
		return ""
	}
	return dir
}

func ensureDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(dir, os.ModeDir|os.ModePerm)
			return err
		} else {
			return err
		}
	}
	return nil
}
