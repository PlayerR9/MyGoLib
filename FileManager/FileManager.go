package FileManager

import (
	"bufio"
	"fmt"
	"os"
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
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
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
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
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
	file, err := os.Create(file_path)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	for _, line := range content {
		file.WriteString(line)
	}

	return nil
}
