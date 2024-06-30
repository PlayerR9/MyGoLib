package Debugging

// Verbose is a struct that can be used to print verbose output.
type Verbose struct {
	// isActive is a boolean that determines if the verbose output is active.
	isActive bool
}

// NewVerbose creates a new Verbose struct.
//
// Returns:
//   - *Verbose: The new Verbose struct.
func NewVerbose() *Verbose {
	v := &Verbose{}

	return v
}

// Printf prints a formatted string if the verbose output is active.
//
// Parameters:
//   - format: The format string.
func (v *Verbose) Printf(format string, args ...interface{}) {
	if !v.isActive {
		return
	}

	// Do nothing
}

// Println prints a string if the verbose output is active.
//
// Parameters:
//   - args: The arguments to print.
func (v *Verbose) Println(args ...interface{}) {
	if !v.isActive {
		return
	}

	// Do nothing
}

// Activate activates or deactivates the verbose output.
//
// Parameters:
//   - active: The boolean that determines if the verbose output is active.
func (v *Verbose) Activate(active bool) {
	v.isActive = active
}

// IsActive returns true if the verbose output is active.
//
// Returns:
//   - bool: True if the verbose output is active, false otherwise.
func (v *Verbose) IsActive() bool {
	return v.isActive
}
