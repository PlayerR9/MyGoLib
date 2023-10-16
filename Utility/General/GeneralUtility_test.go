// git tag v0.1.17

package General

import (
	"testing"
)

func TestConsoleFlagInfo_ToString(t *testing.T) {
	type args struct {
		indent_level int
	}
	tests := []struct {
		name    string
		cfi     ConsoleFlagInfo
		args    args
		wantStr string
	}{
		{
			name: "Test 1",
			cfi: ConsoleFlagInfo{
				Name: "--o",
				Args: []string{"output file", "file name"},
				Description: "The output file of the Test 1 command. It requires two arguments: <output file> and <file name>.\n" +
					"Where <output file> is the name of the output file and <file name> is the name of the file to be outputted.",
				Required: true,
				Callback: func(args ...string) (interface{}, error) {
					return nil, nil
				},
			},
			args: args{
				indent_level: 0,
			},
			wantStr: "Name: --o\nArgs: <output file> <file name>\nDescription: The output file of the Test 1 command. It requires two arguments: <output file> and <file name>.\n" +
				"Where <output file> is the name of the output file and <file name> is the name of the file to be outputted.\nRequired: Yes",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStr := tt.cfi.ToString(tt.args.indent_level); gotStr != tt.wantStr {
				t.Errorf("ConsoleFlagInfo.ToString() = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}

func TestUsageToString(t *testing.T) {
	type args struct {
		executable_name string
		command         string
		flags           []ConsoleFlagInfo
		indent_level    int
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
	}{
		{
			name: "Test 1",
			args: args{
				executable_name: "Test.exe",
				command:         "run",
				flags: []ConsoleFlagInfo{
					{
						Name: "--o",
						Args: []string{"output file", "file name"},
						Description: "The output file of the Test 1 command. It requires two arguments: <output file> and <file name>.\n" +
							"Where <output file> is the name of the output file and <file name> is the name of the file to be outputted.",
						Required: true,
						Callback: func(args ...string) (interface{}, error) {
							return nil, nil
						},
					},
					{
						Name: "--i",
						Args: []string{"input file", "file name"},
						Description: "The input file of the Test 1 command. It requires two arguments: <input file> and <file name>.\n" +
							"Where <input file> is the name of the input file and <file name> is the name of the file to be inputted.",
						Required: false,
						Callback: func(args ...string) (interface{}, error) {
							return nil, nil
						},
					},
				},
				indent_level: 0,
			},
			wantStr: "Usage: Test.exe run --o <output file> <file name> [--i <input file> <file name>]",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStr := UsageToString(tt.args.executable_name, tt.args.command, tt.args.flags, tt.args.indent_level); gotStr != tt.wantStr {
				t.Errorf("UsageToString() = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}

func TestHelpToString(t *testing.T) {
	type args struct {
		executable_name string
		flags           map[string][]ConsoleFlagInfo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1",
			args: args{
				executable_name: "Test.exe",
				flags: map[string][]ConsoleFlagInfo{
					"run": {
						{
							Name: "--o",
							Args: []string{"output file", "file name"},
							Description: "The output file of the Test 1 command. It requires two arguments: <output file> and <file name>." +
								" Where <output file> is the name of the output file and <file name> is the name of the file to be outputted.",
							Required: true,
							Callback: func(args ...string) (interface{}, error) {
								return nil, nil
							},
						},
					},
					"test": {
						{
							Name:        "--t",
							Args:        []string{},
							Description: "The test flag of the Test 2 command. It requires no arguments.",
							Required:    true,
							Callback: func(args ...string) (interface{}, error) {
								return nil, nil
							},
						},
					},
				},
			},
			want: "Usage: Test.exe run --o <output file> <file name>\nFlags:\n\tName: --o\n" +
				"\tArgs: <output file> <file name>\n\tDescription: The output file of the Test 1 command. It requires two arguments: <output file> and <file name>." +
				" Where <output file> is the name of the output file and <file name> is the name of the file to be outputted.\n\tRequired: Yes\n\n" +
				"Usage: Test.exe test --t\nFlags:\n\tName: --t\n\tArgs: None\n\tDescription: The test flag of the Test 2 command. It requires no arguments.",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HelpToString(tt.args.executable_name, tt.args.flags); got != tt.want {
				t.Errorf("HelpToString() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
