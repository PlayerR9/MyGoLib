package FileManager

import (
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
//
// The name of the file is derived from the URL and it does not
// download the file if the name already exists in the destination.
//
// Parameters:
//   - dest: A string representing the path to the directory where the file
//     will be saved.
//   - url: A string representing the URL of the file to download.
//
// Returns:
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
func MediaDownloader(dest, url string, force bool) (string, error) {
	// Extract the name of the file from the URL
	fields := strings.Split(url, "/")
	filePath := path.Join(dest, fields[len(fields)-1])

	exists, err := FileExists(filePath)
	if err != nil {
		return "", fmt.Errorf("could not check if file exists: %w", err)
	}

	if exists && !force {
		// Do nothing
		return filePath, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not send GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("file does not exist on the server: %s", url)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not copy file: %w", err)
	}

	return filePath, nil
}

// GetAllFileNamesInDirectory is a function that retrieves all file names
// in a given directory and their extensions.
//
// Parameters:
//   - directoryPath: A string representing the path to the directory.
//
// Returns:
//   - map[string]string: A map where the keys are the file paths and the
//     values are the file extensions.
//   - error: An error if it fails to read the directory or any of its files.
//
// Behaviors:
//   - The function does not search subdirectories, nor returns directories.
//   - The file paths are relative to the directory path.
//   - The keys contain the full path to the file (including the extension).
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
		return nil, fmt.Errorf("could not read directory: %w", walkError)
	}

	return fileExtensionMap, nil
}

// GetFilesEndingIn is a function that retrieves all files in a given directory
// that have a specified extension.
// This function does not search subdirectories.
//
// Parameters:
//   - directoryPath: A string representing the path to the directory.
//   - extensions: A variadic parameter of strings representing the file extensions
//     to match.
//
// Returns:
//   - []string: A slice of strings representing the paths to the matching files.
//   - error: An error if it fails to read the directory.
func GetFilesEndingIn(directoryPath string, extensions ...string) ([]string, error) {
	matchingFiles := make([]string, 0)

	entries, readError := os.ReadDir(directoryPath)
	if readError != nil {
		return nil, fmt.Errorf("could not read directory: %w", readError)
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
//   - filePath: A string representing the path to the file.
//
// Returns:
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

// GroupFilesByExtension is a function that groups files in a directory by their
// extension and moves them to a folder with the same name as the extension.
//
// Parameters:
//   - dir: A string representing the path to the directory.
//
// Returns:
//   - error: An error if it fails to read the directory, create folders, or move files.
func GroupFilesByExtension(dir string) error {
	files, err := GetAllFileNamesInDirectory(dir)
	if err != nil {
		return err
	}

	groupedFiles := make(map[string][]string)
	for file, ext := range files {
		groupedFiles[ext] = append(groupedFiles[ext], file)
	}

	// Create a folder for each extension and then move the files
	// to their respective folders
	for ext, files := range groupedFiles {
		folder := path.Join(dir, ext)

		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create folder: %w", err)
		}

		for _, file := range files {
			newPath := path.Join(folder, path.Base(file))
			err := os.Rename(file, newPath)
			if err != nil {
				return fmt.Errorf("could not move file: %w", err)
			}
		}
	}

	return nil
}
