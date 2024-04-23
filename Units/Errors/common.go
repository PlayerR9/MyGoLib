package Errors

import (
	"io/fs"
	"log"
	"os"
)

const (
	// DefaultLFPermissions is the default file permissions for creating a log file.
	DefaultLFPermissions fs.FileMode = 0666

	// DefaultLFExtension is the default file extension for creating a log file.
	DefaultLFExtension string = ".log.md"

	// DefaultLFFlags is the default flags for creating a log file.
	DefaultLFFlags int = os.O_CREATE | os.O_WRONLY | os.O_APPEND

	// DefaultLoggerFlags is the default flags for creating a logger.
	DefaultLoggerFlags int = log.LstdFlags | log.Llongfile
)
