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

	// Path to the JSON file.
	path string
}

// NewJSONManager creates a new JSONManager.
//
// Parameters:
//   - path: The path to the JSON file.
//
// Returns:
//   - JSONManager[T]: The new JSONManager.
func NewJSONManager[T JSONEncoder](path string) JSONManager[T] {
	return JSONManager[T]{
		path: path,
	}
}

// ChangePath changes the path of the JSON file.
//
// Parameters:
//   - path: The new path to the JSON file.
func (m *JSONManager[T]) ChangePath(path string) {
	m.path = path
}

// Load loads the data from the JSON file.
//
// Returns:
//   - error: An error if there was an issue loading the data.
func (m *JSONManager[T]) Load() error {
	data, err := os.ReadFile(m.path)
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
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(m.path, data, 0644)
}

// Exists checks if the JSON file exists.
//
// Returns:
//   - bool: True if the file exists, false otherwise.
//   - error: An error if there was an issue checking if the file exists.
func (m *JSONManager[T]) Exists() (bool, error) {
	_, err := os.Stat(m.path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// CreateEmpty creates an empty JSON file with the default values of the data.
// If the file already exists, it will overwrite the file.
//
// Returns:
//   - error: An error if there was an issue creating the empty file.
func (m *JSONManager[T]) CreateEmpty() error {
	data, err := json.Marshal(m.Data.Empty())
	if err != nil {
		return err
	}

	return os.WriteFile(m.path, data, 0644)
}

// Delete deletes the JSON file. No operation is performed if the file does not exist.
//
// Returns:
//   - error: An error if there was an issue deleting the file.
func (m *JSONManager[T]) Delete() error {
	return os.Remove(m.path)
}
