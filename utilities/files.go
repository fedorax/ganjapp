package utilities

import "os"

// FileExists checks if a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// FileDelete attempts to remove a file from disk
func FileDelete(filename string) error {
	return os.Remove(filename)
}
