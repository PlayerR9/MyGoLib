package CmdLineParser

import (
	"strconv"
	"testing"
)

func TestParseCommandInfo(t *testing.T) {
	var (
		args []string = []string{"test", "--v", "5", "--a", "B"}
	)

	cmdline, err := NewCmdLineParser(
		"test",
		[]string{"Test command"},
		NewCmdBuilder().
			SetCmd(
				"test",
				nil,
				nil,
				NewFlagBuilder().
					SetFlag(
						"--v",
						true,
						nil,
						nil,
						NewArgBuilder().
							SetArg(
								"max",
								func(args []string) ([]any, error) {
									num, err := strconv.ParseInt(args[0], 10, 0)
									if err != nil {
										return nil, err
									}

									return []any{int(num)}, nil
								},
							),
					),
			),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %s instead", err.Error())
	}

	parsed, err := cmdline.Parse(args)
	if err != nil {
		t.Fatalf("Expected no error, got %s", err.Error())
	}

	if parsed == nil {
		t.Fatalf("Expected a parsed command, got nil instead")
	}

	t.Fatalf("Only testing the parsing of the command, no further checks implemented")
}
