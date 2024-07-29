package General

import (
	"os/exec"
	"strings"

	uc "github.com/PlayerR9/lib_units/common"
)

// RunInPowerShell is a function that returns a function that runs a program in
// a new PowerShell process.
//
// Upon calling the returned function, a new PowerShell process is started with
// the specified program and arguments. The function returns an error if the
// process cannot be started.
//
// Parameters:
//   - program: The path to the program to run.
//   - args: The arguments to pass to the program.
//
// Return:
//   - MainFunc: A function that runs the program in a new PowerShell process.
func RunInPowerShell(program string, args ...string) uc.MainFunc {
	var builder strings.Builder

	builder.WriteString("'-NoExit', '")
	builder.WriteString(program)
	builder.WriteString("'")

	for _, arg := range args {
		builder.WriteString(", '")
		builder.WriteString(arg)
		builder.WriteRune('\'')
	}

	cmd := exec.Command(
		"powershell", "-Command", "Start-Process", "powershell", "-ArgumentList",
		builder.String(),
	)

	return cmd.Run
}
