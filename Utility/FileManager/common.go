package FileManager

import (
	"bufio"
	"os"
)

// FileExists checks if a file exists at the specified location.
//
// Parameters:
//   - loc: The location of the file.
//
// Returns:
//   - bool: True if the file exists, false otherwise.
func FileExists(loc string) (bool, error) {
	_, err := os.Stat(loc)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// ReadFileByLines reads a file line by line and returns the lines.
//
// Parameters:
//   - loc: The location of the file.
//
// Returns:
//   - []string: The lines of the file.
//   - error: An error if one occurred while reading the file.
func ReadFileByLines(loc string) ([]string, error) {
	file, err := os.Open(loc)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
