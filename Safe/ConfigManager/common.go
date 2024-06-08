package ConfigManager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (jm *JSONManager[T]) writeData(elem T) error {
	dir := filepath.Dir(jm.loc)

	err := os.MkdirAll(dir, jm.dirPerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(jm.loc, os.O_WRONLY|os.O_CREATE, jm.filePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(elem, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (jm *JSONManager[T]) openFile() (T, error) {
	dir := filepath.Dir(jm.loc)

	err := os.MkdirAll(dir, jm.dirPerm)
	if err != nil {
		return *new(T), err
	}

	file, err := os.OpenFile(jm.loc, os.O_RDONLY|os.O_CREATE, jm.filePerm)
	if err != nil {
		return *new(T), err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return *new(T), fmt.Errorf("could not get file info: %w", err)
	}

	size := fileInfo.Size()

	if size == 0 {
		return *new(T), fmt.Errorf("file is empty")
	}

	data := make([]byte, size)
	_, err = file.Read(data)
	if err != nil {
		return *new(T), fmt.Errorf("could not read file: %w", err)
	}

	res := *new(T)

	err = json.Unmarshal(data, res)
	if err != nil {
		return *new(T), fmt.Errorf("could not unmarshal data: %w", err)
	}

	return res, nil
}
