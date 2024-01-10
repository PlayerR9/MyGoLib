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

// MediaDownloader downloads a file from the given URL and saves it
// to the specified destination.
// The name of the file is extracted from the URL.
// If the file already exists at the destination, the function returns
// the path to the existing file.
// If the file does not exist, the function checks if the file exists
// on the server.
// If the file exists on the server, the function downloads the file and
// saves it to the destination.
// The function returns the path to the downloaded file, or an error if
// the download fails.
//
// dest is the directory where the file will be saved.
// url is the URL of the file to download.
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
	file_path := path.Join(dest, fields[len(fields)-1])

	// Check if the file already exists
	_, err := os.Stat(file_path)
	if err == nil {
		return file_path, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	// Check if the file exists on the server
	resp, err := http.Head(url)
	if err != nil {
		return "", err
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download the file: %s", resp.Status)
	}

	// Download the file
	resp, err = http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	file, err := os.Create(file_path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Copy the data to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return file_path, nil
}

// ReadWholeFileLineByLine reads a file from the provided path line by line and returns a slice of strings,
// where each string represents a line from the file. If an error occurs while opening the file, the function
// returns the error and a nil slice.
//
// Parameters:
//
//   - path: The path to the file to be read.
//
// Returns:
//
//   - A slice of strings where each string is a line from the file.
//   - An error if one occurred while opening or scanning the file.
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

// AppendToFile is a function that appends content to a file.
// It takes two parameters: filePath, which is a string representing the path
// to the file, and content, which is a variadic parameter of strings
// representing the content to be appended.
// The function returns an error if it fails to open the file or write to it.
//
// The function first attempts to open the file at the given filePath in append
// and write-only mode. If it fails to open the file, it returns an error wrapped
// with a message indicating that it could not open the file.
//
// If the file is opened successfully, the function defers the closing of the
// file to ensure that it is closed when the function returns, either normally or
// due to a panic.
//
// The function then attempts to write the content to the file. The content
// strings are joined with a newline character ("\n") in between each string.
// If it fails to write to the file, it returns an error wrapped with a message
// indicating that it could not write to the file.
//
// If the function succeeds in writing to the file, it returns nil, indicating
// that no errors occurred.
func AppendToFile(filePath string, content ...string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(content, "\n"))
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

// GetAllFileNamesInDirectory is a function that retrieves all file names
// in a given directory and their extensions.
// It takes a parameter, directoryPath, which is a string representing the
// path to the directory.
// The function returns a map where the keys are the file paths and the
// values are the file extensions, and an error if it fails to read the
// directory or any of its files.
//
// The function first initializes an empty map, fileExtensionMap, to store
// the file paths and their extensions.
//
// The function then calls the filepath.Walk function to traverse the
// directory tree. filepath.Walk takes a root path and a function to be
// called for each file or directory in the tree.
// The function passed to filepath.Walk takes the current path, the file
// or directory info, and an error as parameters.
// If an error occurs while retrieving the info for a file or directory,
// the function returns the error to filepath.Walk, which stops the walk
// and returns the error.
// If the info indicates that the current path is a file (not a directory),
// the function adds the file path and its extension to fileExtensionMap.
//
// If filepath.Walk returns an error, the function wraps the error with a
// message indicating that it could not read the directory and returns the
// wrapped error.
//
// If filepath.Walk does not return an error, the function returns
// fileExtensionMap and nil, indicating that no errors occurred.
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
		return nil, fmt.Errorf("could not read directory: %v", walkError)
	}

	return fileExtensionMap, nil
}

// GetFilesEndingIn is a function that retrieves all files in a given directory
// that have a specified extension.
// It takes two parameters: directoryPath, which is a string representing the path
// to the directory, and extensions, which is a variadic parameter of strings
// representing the file extensions to match.
// The function returns a slice of strings representing the paths to the matching
// files, and an error if it fails to read the directory.
//
// The function first initializes an empty slice, matchingFiles, to store the
// paths to the matching files.
//
// The function then attempts to read the directory at the given directoryPath.
// If it fails to read the directory, it returns nil and an error wrapped with
// a message indicating that it could not read the directory.
//
// If the directory is read successfully, the function iterates over the entries
// in the directory.
// For each entry, the function checks if the entry is not a directory and if
// its extension is in the list of extensions to match.
// If both conditions are true, the function appends the path to the file to
// matchingFiles. The path is constructed by joining the directoryPath and the
// name of the entry.
//
// After all entries have been processed, the function returns matchingFiles and
// nil, indicating that no errors occurred.
func GetFilesEndingIn(directoryPath string, extensions ...string) ([]string, error) {
	matchingFiles := make([]string, 0)

	entries, readError := os.ReadDir(directoryPath)
	if readError != nil {
		return nil, fmt.Errorf("could not read directory: %v", readError)
	}

	for _, entry := range entries {
		if !entry.IsDir() && slices.Contains(extensions, filepath.Ext(entry.Name())) {
			matchingFiles = append(matchingFiles, filepath.Join(directoryPath, entry.Name()))
		}
	}

	return matchingFiles, nil
}

// SplitPath splits a file path into its components.
// It takes a file path as a string and returns a slice of strings.
// The function iterates over the file path, splitting it into parts.
// Each part is added to the parts slice.
// If a part is not empty, it is appended to the parts slice.
// After all parts have been added, the parts slice is reversed to ensure
// the parts are in the correct order.
// The function then returns the parts slice.
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
