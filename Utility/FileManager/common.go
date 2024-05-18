package FileManager

import (
	"os"
)

// fileExists checks if a file exists at the given path.
//
// Parameters:
//   - loc: The path to the file to check for existence.
//
// Returns:
//   - bool: A boolean indicating whether the file exists.
//   - error: An error if one occurred while checking the file.
func fileExists(loc string) (bool, error) {
	_, err := os.Stat(loc)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// FileManager represents a file manager that can be used to read and write files.
type FileManager interface {
	// GetLocation returns the location of the file.
	//
	// Returns:
	//   - string: The location of the file.
	GetLocation() string

	// Exists checks if the file exists at the location.
	//
	// Returns:
	//   - bool: A boolean indicating whether the file exists.
	//   - error: An error if one occurred while checking the file.
	Exists() (bool, error)
}
