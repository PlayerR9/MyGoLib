package FileManager

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// FileReader represents a file reader that can be used to read files.
type FileReader struct {
	// loc is the location of the file.
	loc string
}

// GetLocation returns the location of the file.
//
// Returns:
//   - string: The location of the file.
func (fr *FileReader) GetLocation() string {
	return fr.loc
}

// NewFileReader creates a new FileReader with the given location.
//
// Parameters:
//   - loc: A string representing the location of the file.
//
// Returns:
//   - *FileReader: A pointer to the newly created FileReader.
func NewFileReader(loc string) *FileReader {
	return &FileReader{
		loc: loc,
	}
}

// Exists checks if the file exists at the location.
//
// Returns:
//   - bool: A boolean indicating whether the file exists.
//   - error: An error if one occurred while checking the file.
func (fr *FileReader) Exists() (bool, error) {
	return fileExists(fr.loc)
}

// ReadFile reads the content of a file from the provided path and returns
// the content as a string.
//
// Parameters:
//   - filePath: A string representing the path to the file to be read.
//
// Returns:
//   - string: A string representing the content of the file.
//   - error: An error of type *os.PathError if the file could not be opened
//     or read.
//
// Behaviors:
//   - The function reads the entire file into memory.
//   - If an error occurs, the function returns the error and an empty string.
func (fr *FileReader) Read() (string, error) {
	// 1. Open the file
	file, err := os.Open(fr.loc)
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// 2. Read the file
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("could not get file information: %w", err)
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		panic(errors.New("buffer size is not equal to file size"))
	}

	return string(buffer), nil
}

// Lines reads the file line by line and returns a slice of strings, where each
// string represents a line from the file.
//
// Returns:
//   - []string: A slice of strings where each string is a line from the file.
//   - error: An error if one occurred while opening or scanning the file.
//
// Behaviors:
//   - The function reads the file line by line. Each line does not include the
//     newline character.
//   - If an error occurs, the function returns the error and the lines read up
//     to that point.
func (fr *FileReader) Lines() ([]string, error) {
	f, err := os.Open(fr.loc)
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
