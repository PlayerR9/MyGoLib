// git tag v0.1.35

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

func AppendToFile(filePath string, content ...string) error {
	file, openError := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if openError != nil {
		return fmt.Errorf("could not open file: %v", openError)
	}
	defer file.Close()

	_, writeError := file.WriteString(strings.Join(content, "\n"))
	if writeError != nil {
		return fmt.Errorf("could not write to file: %v", writeError)
	}

	return nil
}

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
