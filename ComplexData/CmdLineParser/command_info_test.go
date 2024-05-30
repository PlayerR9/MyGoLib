package CmdLineParser

import (
	"strconv"
	"testing"

	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
)

func TestParseCommandInfo(t *testing.T) {
	var (
		args []string = []string{"test", "--v", "5", "--a", "B"}
	)

	command := NewCmdBuilder().
		SetCmd(
			"test",
			nil,
			nil,
			NewFlagBuilder().
				SetFlag(
					"--v",
					false,
					nil,
					nil,
					NewArgBuilder().
						SetArg(
							"max",
							func(s string) (any, error) {
								num, err := strconv.ParseInt(s, 10, 0)
								if err != nil {
									return nil, err
								}

								return int(num), nil
							},
						),
				),
		)

	commands, err := command.Build()
	if err != nil {
		t.Fatalf("Expected no error, got %s", err.Error())
	} else if len(commands) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(commands))
	}

	cmd := commands[0]

	_, err = cmd.Parse(args)
	if err != nil {
		if !ue.As[*ue.ErrIgnorable](err) {
			t.Fatalf("Expected no error, got %s", err.Error())
		} else {
			t.Fatalf("Ignorable error: %s", err.Error())
		}
	}
}
