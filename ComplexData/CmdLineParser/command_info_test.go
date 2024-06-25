package CmdLineParser

import (
	"strconv"
	"testing"

	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

func TestParseCommandInfo(t *testing.T) {
	var (
		TestArg []string = []string{"test", "--v", "5", "--a", "B"}
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

	parsed, err := cmdline.Parse(TestArg)
	if err != nil {
		ok := ue.Is[*ue.ErrIgnorable](err)

		if ok {
			t.Fatalf("As expected, got ignorable error: %s", err.Error())
		} else {
			t.Fatalf("Expected no error, got %s", err.Error())
		}
	}

	if parsed == nil {
		t.Fatalf("Expected a parsed command, got nil instead")
	}

	t.Fatalf("Only testing the parsing of the command, no further checks implemented")
}

func TestVariadicArgument(t *testing.T) {
	var (
		TestArg []string = []string{"test", "--s", "1", "2", "3"}
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
						"--s",
						false,
						nil,
						nil,
						NewArgBuilder().
							SetArg(
								"max 1-3",
								func(args []string) ([]any, error) {
									numbers := make([]any, 0)

									for _, arg := range args {
										num, err := strconv.ParseInt(arg, 10, 0)
										if err != nil {
											return nil, err
										}

										numbers = append(numbers, int(num))
									}

									return numbers, nil
								},
							),
					),
			),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %s instead", err.Error())
	}

	parsed, err := cmdline.Parse(TestArg)
	if err != nil {
		ok := ue.Is[*ue.ErrIgnorable](err)

		if ok {
			t.Fatalf("As expected, got ignorable error: %s", err.Error())
		} else {
			t.Fatalf("Expected no error, got %s", err.Error())
		}
	}

	if parsed == nil {
		t.Fatalf("Expected a parsed command, got nil instead")
	}

	t.Fatalf("Only testing the parsing of the command, no further checks implemented")
}
