package FileManager

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

// FileExists checks if a file exists at the given path.
//
// Parameters:
//
//   - filePath: A string representing the path to the file.
//
// Returns:
//
//   - A boolean indicating whether the file exists.
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// MediaDownloader downloads a file from the given URL and saves it
// to the specified destination.
//
// The name of the file is derived from the URL and it does not
// download the file if the name already exists in the destination.
//
// Parameters:
//
//   - dest: A string representing the path to the directory where the file
//     will be saved.
//   - url: A string representing the URL of the file to download.
//
// Returns:
//
//   - string: The path to the downloaded file.
//   - error: An error if the download fails.
//
// Example:
//
//	file_path, err := MediaDownloader("/path/to/destination", "http://example.com/file.mp3")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(file_path) // Output: /path/to/destination/file.mp3
func MediaDownloader(dest, url string) (string, error) {
	// Extract the name of the file from the URL
	fields := strings.Split(url, "/")
	filePath := path.Join(dest, fields[len(fields)-1])

	exists, err := FileExists(filePath)
	if err != nil {
		return "", fmt.Errorf("could not check if file exists: %s", err.Error())
	} else if exists {
		return filePath, nil // Do nothing
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not send GET request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("file does not exist on the server: %s", url)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("could not create file: %s", err.Error())
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not copy file: %s", err.Error())
	}

	return filePath, nil
}

// ReadWholeFileLineByLine reads a file from the provided path line by line
// and returns a slice of strings, where each string represents a line from
// the file. If an error occurs while opening the file, the function
// returns the error and a nil slice.
//
// Parameters:
//
//   - path: The path to the file to be read.
//
// Returns:
//
//   - []string: A slice of strings where each string is a line from the file.
//   - error: An error if one occurred while opening or scanning the file.
func ReadWholeFileLineByLine(path string) ([]string, error) {
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

// AppendToFile is a function that appends content to a file, each content
// on a new line.
//
// Panics if it fails to open the file or write to it.
//
// Parameters:
//
//   - filePath: A string representing the path to the file.
//   - contents: A variadic parameter of strings representing the content
//     to be appended.
func AppendToFile(filePath string, contents ...string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %s", err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(contents, "\n"))
	if err != nil {
		return fmt.Errorf("could not write to file: %s", err.Error())
	}

	return nil
}

// GetAllFileNamesInDirectory is a function that retrieves all file names
// in a given directory and their extensions.
//
// Parameters:
//
//   - directoryPath: A string representing the path to the directory.
//
// Returns:
//
//   - map[string]string: A map where the keys are the file paths and the
//     values are the file extensions.
//   - error: An error if it fails to read the directory or any of its files.
func GetAllFileNamesInDirectory(directoryPath string) (map[string]string, error) {
	fileExtensionMap := make(map[string]string)

	walkError := filepath.Walk(directoryPath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileExtensionMap[currentPath] = filepath.Ext(info.Name())
		}
		return nil
	})

	if walkError != nil {
		return nil, fmt.Errorf("could not read directory: %s", walkError.Error())
	}

	return fileExtensionMap, nil
}

// GetFilesEndingIn is a function that retrieves all files in a given directory
// that have a specified extension.
// This function does not search subdirectories.
//
// Parameters:
//
//   - directoryPath: A string representing the path to the directory.
//   - extensions: A variadic parameter of strings representing the file extensions
//     to match.
//
// Returns:
//
//   - []string: A slice of strings representing the paths to the matching files.
//   - error: An error if it fails to read the directory.
func GetFilesEndingIn(directoryPath string, extensions ...string) ([]string, error) {
	matchingFiles := make([]string, 0)

	entries, readError := os.ReadDir(directoryPath)
	if readError != nil {
		return nil, fmt.Errorf("could not read directory: %s", readError.Error())
	}

	for _, entry := range entries {
		if !entry.IsDir() && slices.Contains(extensions, filepath.Ext(entry.Name())) {
			matchingFiles = append(matchingFiles, filepath.Join(directoryPath, entry.Name()))
		}
	}

	return matchingFiles, nil
}

// SplitPath splits a file path into its components, where each component
// is a part of the file path.
//
// Parameters:
//
//   - filePath: A string representing the path to the file.
//
// Returns:
//
//   - []string: A slice of strings representing the parts of the file path.
func SplitPath(filePath string) []string {
	var parts []string

	for filePath != "" {
		var part string

		filePath, part = path.Split(filePath)
		if part != "" {
			parts = append(parts, part)
		}
	}

	slices.Reverse(parts)

	return parts
}
