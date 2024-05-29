package FileManager

import (
	"encoding/json"
	"fmt"
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
	Empty() json.Marshaler
}

// JSONManager is a manager for saving and loading data to and from a JSON file.
type JSONManager[T JSONEncoder] struct {
	// Data is the data to save and load.
	Data T

	*fileHandler
}

// Create implements the FileManager interface.
func (m *JSONManager[T]) Create() error {
	err := m.fileHandler.Create()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(m.Data.Empty(), "", "  ")
	if err != nil {
		return err
	}

	_, err = m.file.Write(data)
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
//   - filePerm: The permissions for the file.
//
// Returns:
//   - JSONManager[T]: The new JSONManager.
func NewJSONManager[T JSONEncoder](loc string, dirPerm, filePerm os.FileMode) *JSONManager[T] {
	return &JSONManager[T]{
		fileHandler: newFileHandler(loc, dirPerm, filePerm, os.O_RDWR),
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
	if m.file == nil {
		// Open the file if it is not already open
		err := m.fileHandler.Open()
		if err != nil {
			return err
		}

		return json.Unmarshal(nil, m)
	}

	info, err := m.file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %w", err)
	}

	size := info.Size()
	data := make([]byte, size)

	_, err = m.file.Read(data)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	return json.Unmarshal(data, m)
}

// Save saves the data to the JSON file. It will overwrite the file if it already exists.
//
// Returns:
//   - error: An error if there was an issue saving the data.
func (m *JSONManager[T]) Save() error {
	if m.file == nil {
		// Create the file if it does not exist
		err := m.fileHandler.Create()
		if err != nil {
			return err
		}
	}

	data, err := json.MarshalIndent(m.Data, "", "  ")
	if err != nil {
		return err
	}

	err = m.file.Truncate(0)
	if err != nil {
		return fmt.Errorf("could not truncate file: %w", err)
	}

	_, err = m.file.WriteAt(data, 0)
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}
