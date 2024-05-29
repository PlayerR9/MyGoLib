package FileManager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// JSONEncoder is an interface for encoding and decoding JSON data.
type JSONEncoder interface {
	json.Marshaler
	json.Unmarshaler

	// Empty returns the default values of the data.
	//
	// Returns:
	//   - JSONEncoder: The default values of the data.
	Empty() json.Marshaler
}

// JSONManager is a manager for saving and loading data to and from a JSON file.
type JSONManager[T JSONEncoder] struct {
	// Data is the data to save and load.
	Data T

	// loc is the location of the JSON file.
	loc string

	// dirPerm is the permissions of the directory.
	dirPerm os.FileMode
}

// Create implements the FileManager interface.
func (m *JSONManager[T]) Create() error {
	dir := filepath.Dir(m.loc)

	err := os.MkdirAll(dir, m.dirPerm)
	if err != nil {
		return err
	}

	file, err := os.Create(m.loc)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(m.Data.Empty(), "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// NewJSONManager creates a new JSONManager.
//
// Parameters:
//   - loc: The path to the JSON file.
//   - dirPerm: The permissions for the directory.
//
// Returns:
//   - *JSONManager[T]: The new JSONManager.
func NewJSONManager[T JSONEncoder](loc string, dirPerm os.FileMode) *JSONManager[T] {
	return &JSONManager[T]{
		loc:     loc,
		dirPerm: dirPerm,
	}
}

// ChangePath changes the path of the JSON file.
//
// Parameters:
//   - path: The new path to the JSON file.
func (m *JSONManager[T]) ChangePath(path string) {
	m.loc = path
}

// Load loads the data from the JSON file.
//
// Returns:
//   - error: An error if there was an issue loading the data.
func (m *JSONManager[T]) Load() error {
	data, err := os.ReadFile(m.loc)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	err = json.Unmarshal(data, m.Data)
	if err != nil {
		return fmt.Errorf("could not unmarshal data: %w", err)
	}

	return nil
}

// Save saves the data to the JSON file. It will overwrite the file if it already exists.
//
// Returns:
//   - error: An error if there was an issue saving the data.
func (m *JSONManager[T]) Save() error {
	data, err := json.MarshalIndent(m.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal data: %w", err)
	}

	err = os.WriteFile(m.loc, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

// GetLocation returns the location of the file.
//
// Returns:
//   - string: The location of the file.
func (m *JSONManager[T]) GetLocation() string {
	return m.loc
}

// GetPermissions returns the permissions of the file.
//
// Returns:
//   - [2]os.FileMode: An array of os.FileMode representing the permissions of the directory and file.
func (m *JSONManager[T]) GetPermissions() [2]os.FileMode {
	return [2]os.FileMode{m.dirPerm, 0644}
}

// Exists checks if a file exists at the location of the JSONManager[T].
//
// Returns:
//   - bool: A boolean indicating whether the file exists.
//   - error: An error if one occurred while checking the file.
func (m *JSONManager[T]) Exists() (bool, error) {
	_, err := os.Stat(m.loc)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// Delete deletes the file at the location.
//
// Returns:
//   - error: An error if one occurred while deleting the file.
func (m *JSONManager[T]) Delete() error {
	return os.Remove(m.loc)
}
