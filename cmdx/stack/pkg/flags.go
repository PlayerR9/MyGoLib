package pkg

import (
	"errors"
	"flag"
	"fmt"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	utgo "github.com/PlayerR9/MyGoLib/Utility/Go"
	utse "github.com/PlayerR9/MyGoLib/Utility/StringExt"
)

var (
	// DataType is the data type of the linked stack.
	DataType *string

	// TypeName is the name of the linked stack.
	TypeName *string

	// OutputDir is the output directory of the generated code.
	OutputDir *string
)

func init() {
	DataType = flag.String("type", "", "the data type of the linked stack. This must be set and it is "+
		"the data type of the linked stack.")

	TypeName = flag.String("name", "", "the name of the linked stack. Must be a valid Go identifier. If not set, "+
		"the default name of 'Linked<DataType>Stack' will be used instead.")

	OutputDir = flag.String("o", "stack.go", "the output directory of the generated code. If not set, the default "+
		"output directory of '<type>_stack.go' will be used instead.")
}

func ParseFlags() error {
	flag.Parse()

	if *DataType == "" {
		return errors.New("type must be set")
	}

	if *TypeName == "" {
		dt, err := utse.Title(*DataType)
		uc.AssertErr(err, "Utility.Title(%s)", *DataType)
		if err != nil {
			return err
		}

		*TypeName = "Linked" + dt + "Stack"
	} else {
		err := utgo.IsValidName(*TypeName, nil)
		if err != nil {
			err := fmt.Errorf("name of the type is invalid: %w", err)
			return err
		}
	}

	dest_loc, err := utgo.FixImportDir(*OutputDir)
	if err != nil {
		return err
	}

	*OutputDir = dest_loc

	return nil
}
