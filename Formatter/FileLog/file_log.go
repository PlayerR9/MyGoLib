package FileLog

import (
	"fmt"
	"log"
	"os"
	"time"

	sext "github.com/PlayerR9/MyGoLib/Utility/StringExt"
)

type FileLogger struct {
	fileName string

	file *os.File
	*log.Logger
}

func (fl *FileLogger) Close() {
	fl.file.Close()
	fl.file = nil

	fl.Logger = nil
}

func NewFileLogger(filePath string) (*FileLogger, error) {
	fl := &FileLogger{
		fileName: fmt.Sprintf("%s.log.md", filePath),
		Logger:   log.New(os.Stdout, "", log.LstdFlags|log.Llongfile),
	}

	// Open the file for writing
	var err error

	fl.file, err = os.OpenFile(fl.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Write the current time to the file

	currentTime := time.Now()

	fmt.Fprintf(fl.file, "## Log of %v (%v : %v)\n",
		sext.DateStringer(currentTime), currentTime.Hour(), currentTime.Minute())

	return fl, nil
}

func (fl *FileLogger) GetFileName() string {
	return fl.fileName
}
