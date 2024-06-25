package Debugging

type Verbose struct {
	isActive bool
}

func NewVerbose() *Verbose {
	v := &Verbose{}

	return v
}

func (v *Verbose) Printf(format string, args ...interface{}) {
	if !v.isActive {
		return
	}

	// Do nothing
}

func (v *Verbose) Println(args ...interface{}) {
	if !v.isActive {
		return
	}

	// Do nothing
}

func (v *Verbose) Activate(active bool) {
	v.isActive = active
}

func (v *Verbose) IsActive() bool {
	return v.isActive
}
