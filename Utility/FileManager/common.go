package FileManager

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
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

	ok := errors.Is(err, os.ErrExist)
	if ok {
		return false, nil
	} else {
		return false, err
	}
}

// Create creates the file at the location.
//
// Parameters:
//   - loc: A string representing the path to the file to be created.
//   - dirPerm: The permission to set for the directory containing the file.
//   - filePerm: The permission to set for the file.
//
// Returns:
//   - error: An error if one occurred while creating the file.
//
// Behaviors:
//   - dirPerm and filePerm are optional. If 0 is provided, the default permissions
//     are used: DP_OwnerRestrictOthers for directories and FP_OwnerRestrictOthers for files.
//   - If the file already exists, it closes the previous file and creates a new one.
//   - Once the file is opened, it is kept open until the FileManager is closed.
func Create(loc string, dirPerm, filePerm os.FileMode) (*os.File, error) {
	dir := filepath.Dir(loc)

	if dirPerm == 0 {
		dirPerm = DP_OwnerRestrictOthers
	}

	if filePerm == 0 {
		filePerm = FP_OwnerRestrictOthers
	}

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(loc, os.O_CREATE|os.O_RDWR, filePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// ReadFile reads the content of a file from the provided path and returns
// the content as a string.
//
// Parameters:
//   - filePath: A string representing the path to the file to be read.
//   - create: A boolean indicating whether to create the file if it does not exist.
//
// Returns:
//   - string: A string representing the content of the file.
//   - error: An error of type *os.PathError if the file could not be opened
//     or read.
//
// Behaviors:
//   - The function reads the entire file into memory.
//   - If an error occurs, the function returns the error and an empty string.
func Read(loc string, create bool) (string, error) {
	exists, err := FileExists(loc)
	if err != nil {
		return "", err
	}

	var file *os.File

	if !exists {
		if !create {
			return "", os.ErrNotExist
		}

		file, err = Create(loc, 0, 0)
		if err != nil {
			return "", err
		}
	} else {
		file, err = os.Open(loc)
		if err != nil {
			return "", err
		}
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("could not get file information: %w", err)
	}

	fileSize := info.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}

	return string(buffer), nil
}

// Lines reads the file line by line and returns a slice of strings, where each
// string represents a line from the file.
//
// Parameters:
//   - loc: The location of the file.
//   - create: A boolean indicating whether to create the file if it does not exist.
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
func Lines(loc string, create bool) ([]string, error) {
	exists, err := FileExists(loc)
	if err != nil {
		return nil, err
	}

	var file *os.File

	if !exists {
		file, err = Create(loc, 0, 0)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.Open(loc)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// PerLine reads the file line by line and applies a function to each line.
//
// Parameters:
//   - loc: The location of the file.
//   - create: A boolean indicating whether to create the file if it does not exist.
//   - f: A function that processes a line of text and returns a value.
//
// Returns:
//   - []T: A slice of values returned by the function for each line.
//   - error: An error if one occurred while opening or scanning the file.
//
// Behaviors:
//   - The function reads the file line by line and applies the function f to each line.
//   - If an error occurs, the function returns the error and the values processed up to that point.
func PerLine[T any](loc string, create bool, f uc.EvalOneFunc[string, T]) ([]T, error) {
	exists, err := FileExists(loc)
	if err != nil {
		return nil, err
	}

	var file *os.File

	if !exists {
		file, err = Create(loc, 0, 0)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.Open(loc)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	var lines []T
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()

		res, err := f(text)
		if err != nil {
			return nil, fmt.Errorf("could not process line: %w", err)
		}

		lines = append(lines, res)
	}

	return lines, scanner.Err()
}

// CheckPath checks if the path is a directory or file as expected.
//
// Parameters:
//   - loc: A string representing the path to check.
//   - isDir: A boolean indicating whether the path should be a directory.
//
// Returns:
//   - bool: True if the path is as expected, false otherwise.
//   - error: An error if the path is not as expected or if an error occurred
//     while checking the path.
func CheckPath(loc string, isDir bool) (bool, error) {
	if loc == "" {
		return false, ue.NewErrEmpty(loc)
	}

	stat, err := os.Stat(loc)
	if err != nil {
		return false, err
	}

	ok := isDir == stat.IsDir()
	return ok, nil
}
