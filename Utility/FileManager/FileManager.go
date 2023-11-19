// git tag v0.1.22

package FileManager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

func ReadWholeFile(filePath string) ([]string, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(fileContent), "\n"), nil
}

func FileExists(filePath string) (bool, error) {
	_, fileStatError := os.Stat(filePath)
	if fileStatError == nil {
		return true, nil
	}
	if os.IsNotExist(fileStatError) {
		return false, nil
	}
	return false, fileStatError
}

func WriteToFile(filePath string, content ...string) error {
	writeError := os.WriteFile(filePath, []byte(strings.Join(content, "\n")), 0644)
	if writeError != nil {
		return fmt.Errorf("could not write to file: %v", writeError)
	}
	return nil
}

func CreateEmptyFile(filePath string) error {
	_, createError := os.Create(filePath)
	if createError != nil {
		return fmt.Errorf("could not create file: %v", createError)
	}
	return nil
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
