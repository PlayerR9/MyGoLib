package go_generator

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// InitLogger initializes the logger with the given prefix.
//
// Parameters:
//   - prefix: The prefix to use for the logger.
//
// Returns:
//   - *log.Logger: The initialized logger. Never nil.
//
// If the prefix is empty, it defaults to "go_generator".
func InitLogger(prefix string) *log.Logger {
	if prefix == "" {
		prefix = "go_generator"
	}

	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(prefix)
	builder.WriteString("]: ")

	logger_prefix := builder.String()

	logger := log.New(os.Stdout, logger_prefix, log.Lshortfile)
	return logger
}

// Generate generates code using the given generator and writes it to the given destination file.
//
// WARNING: Remember to call this function iff the function go_generator.SetOutputFlag() was called
// and only after the function flag.Parse() was called.
//
// Parameters:
//   - actual_loc: The actual location of the generated code.
//   - data: The data to use for the generated code.
//   - t: The template to use for the generated code.
//
// Returns:
//   - error: An error if occurred.
//
// Errors:
//   - *common.ErrInvalidParameter: If any of the parameters is nil or if the actual_loc is an empty string when the
//     IsOutputLocRequired flag was set and the output location was not defined.
//   - error: Any other type of error that may have occurred.
//
// The parameter actual_loc is only used if the OutputLoc flag was not specified and the IsOutputLocRequired flag
// was not set. Otherwise, it is ignored.
func Generate(actual_loc string, data any, t *template.Template) error {
	if data == nil {
		return uc.NewErrNilParameter("data")
	} else if t == nil {
		return uc.NewErrNilParameter("t")
	}

	var output_loc string

	if OutputLoc == nil {
		return uc.NewErrInvalidUsage(
			errors.New("output location was not defined"),
			"Please call the go_generator.SetOutputFlag() function before calling this function.",
		)
	} else {
		output_loc = *OutputLoc
	}

	if output_loc == "" {
		if IsOutputLocRequired {
			return errors.New("flag must be set")
		} else if actual_loc == "" {
			return uc.NewErrInvalidParameter("actual_loc", uc.NewErrEmpty(actual_loc))
		}

		output_loc = actual_loc
	}

	ext := filepath.Ext(output_loc)
	if ext == "" {
		return errors.New("location cannot be a directory")
	} else if ext != ".go" {
		return errors.New("location must be a .go file")
	}

	var buff bytes.Buffer

	err := t.Execute(&buff, data)
	if err != nil {
		return err
	}

	res := buff.Bytes()

	err = os.WriteFile(output_loc, res, 0644)
	if err != nil {
		return err
	}

	return nil
}
