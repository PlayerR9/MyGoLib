package FileManager

import (
	"bufio"
	"os"
)

/*
	ReadWholeFile reads a file and returns a slice of strings, each string representing a line in the file.

Parameters:

	path: The path to the file to read.

Returns:

	A slice of strings, each string representing a line in the file.
	Error: If the file could not be opened, or if there was an error reading the file.
*/
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
