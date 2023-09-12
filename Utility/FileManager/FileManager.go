package FileManager

import (
	"bufio"
	"fmt"
	"os"

	sext "github.com/PlayerR9/MyGoLib/Utility/StrExt"
)

// GLOBAL VARIABLES
var (
	// DebugMode is a boolean that is used to enable or disable debug mode. When debug mode is enabled, the package will print debug messages.
	// **Note:** Debug mode is disabled by default.
	DebugMode bool = false

	// debugger *log.Logger = log.New(os.Stdout, "[FileManager] ", log.LstdFlags) // Debugger
)

// ReadWholeFile reads a file and returns a slice of strings, each string representing a line in the file.
//
// Parameters:
//   - path: The path to the file to read.
//
// Returns:
//   - []string: A slice of strings, each string representing a line in the file.
//   - error: If the file could not be opened, or if there was an error reading the file.
func ReadWholeFile(path string) ([]string, error) {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Initialize variables
	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)

	// Read the file
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // Add the line to the slice of lines
	}

	return lines, scanner.Err()
}

// FileExists checks if a file exists at the given path.
//
// Parameters:
//   - path: The path to the file to check.
//
// Returns:
//   - bool: True if the file exists, false if it does not.
//   - error: If there was an error checking if the file exists.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path) // Get the file info

	// Check if the file exists
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// WriteToFile writes the given content to the given file. If the file does not exist, it will be created. If the file does exist, it will be overwritten.
//
// Parameters:
//   - file_path: The path to the file to write to.
//   - content: The content to write to the file.
//
// Returns:
//   - Error: If the file could not be created or written to.
//
// Information:
//   - Each string in the content slice will be written to the file consecutively. To write a new line, add a newline character to the end of the string.
func WriteToFile(file_path string, content ...string) error {
	// Create the file
	file, err := os.Create(file_path)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	// Write to the file
	for i, line := range content {
		_, err := file.WriteString(line)
		if err != nil {
			// Could not write to file, return error.
			return fmt.Errorf("could not write %s element of content to file: %v", sext.PrintOrdinalNumber(i+1), err)
		}
	}

	return nil
}

// CreateEmptyFile creates an empty file at the given path. If the file already exists, it will be overwritten.
//
// Parameters:
//   - file_path: The path to the file to create.
//
// Returns:
//   - Error: If the file could not be created.
func CreateEmptyFile(file_path string) error {
	// Create the file
	file, err := os.Create(file_path)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	return nil
}

// AppendToFile appends the given content to the given file. If the file does not exist, it will give an error.
//
// Parameters:
//   - file_path: The path to the file to append to.
//   - content: The content to append to the file.
//
// Returns:
//   - Error: If the file could not be created or written to.
//
// Information:
//   - Each string in the content slice will be written to the file consecutively. To write a new line, add a newline character to the end of the string.
func AppendToFile(file_path string, content ...string) error {
	// Open the file
	file, err := os.OpenFile(file_path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Write to the file
	for i, line := range content {
		_, err := file.WriteString(line)
		if err != nil {
			// Could not write to file, return error.
			return fmt.Errorf("could not write %s element of content to file: %v", sext.PrintOrdinalNumber(i+1), err)
		}
	}

	return nil
}
