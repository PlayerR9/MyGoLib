package FileManager

import (
	"bufio"
	"fmt"
	"os"
)

// FileReader represents a file reader that can be used to read files.
type FileReader struct {
	*fileHandler
}

// NewFileReader creates a new FileReader with the given location.
//
// Parameters:
//   - loc: A string representing the location of the file.
//
// Returns:
//   - *FileReader: A pointer to the newly created FileReader.
func NewFileReader(loc string, dirPerm, filePerm os.FileMode) *FileReader {
	return &FileReader{
		fileHandler: newFileHandler(loc, dirPerm, filePerm, os.O_RDONLY),
	}
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
	if fr.file == nil {
		// Create the file if it does not exist
		err := fr.fileHandler.Create()
		if err != nil {
			return "", fmt.Errorf("could not create file: %w", err)
		}

		return "", nil
	}

	info, err := fr.file.Stat()
	if err != nil {
		return "", fmt.Errorf("could not get file information: %w", err)
	}

	fileSize := info.Size()
	buffer := make([]byte, fileSize)

	_, err = fr.file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
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
	if fr.file == nil {
		// Create the file if it does not exist
		err := fr.fileHandler.Create()
		if err != nil {
			return nil, fmt.Errorf("could not create file: %w", err)
		}

		return nil, nil
	}

	var lines []string
	scanner := bufio.NewScanner(fr.file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
