package FileManager

import (
	"encoding/json"
	"os"
)

// JSONEncoder is an interface for encoding and decoding JSON data.
type JSONEncoder interface {
	json.Marshaler
	json.Unmarshaler

	// Empty returns the default values of the data.
	//
	// Returns:
	//   - JSONEncoder: The default values of the data.
	Empty() JSONEncoder
}

// JSONManager is a manager for saving and loading data to and from a JSON file.
type JSONManager[T JSONEncoder] struct {
	// Data is the data to save and load.
	Data T

	// loc is the location to the JSON file.
	loc string
}

// GetLocation returns the path to the JSON file.
//
// Returns:
//   - string: The path to the JSON file.
func (m *JSONManager[T]) GetLocation() string {
	return m.loc
}

// Exists checks if the JSON file exists.
//
// Returns:
//   - bool: True if the file exists, false otherwise.
//   - error: An error if there was an issue checking if the file exists.
func (m *JSONManager[T]) Exists() (bool, error) {
	return fileExists(m.loc)
}

// NewJSONManager creates a new JSONManager.
//
// Parameters:
//   - loc: The path to the JSON file.
//
// Returns:
//   - JSONManager[T]: The new JSONManager.
func NewJSONManager[T JSONEncoder](loc string) *JSONManager[T] {
	return &JSONManager[T]{
		loc: loc,
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
		return err
	}

	return json.Unmarshal(data, m)
}

// Save saves the data to the JSON file. It will overwrite the file if it already exists.
//
// Returns:
//   - error: An error if there was an issue saving the data.
func (m *JSONManager[T]) Save() error {
	data, err := json.MarshalIndent(m.Data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.loc, data, 0644)
}

// CreateEmpty creates an empty JSON file with the default values of the data.
// If the file already exists, it will overwrite the file.
//
// Returns:
//   - error: An error if there was an issue creating the empty file.
func (m *JSONManager[T]) CreateEmpty() error {
	data, err := json.MarshalIndent(m.Data.Empty(), "", "  ")
	if err != nil {
		return err
	}

	f, err := CreateAll(m.loc, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

// Delete deletes the JSON file. No operation is performed if the file does not exist.
//
// Returns:
//   - error: An error if there was an issue deleting the file.
func (m *JSONManager[T]) Delete() error {
	return os.Remove(m.loc)
}
