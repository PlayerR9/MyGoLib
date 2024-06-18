package Strings

import (
	"testing"
)

func TestTabStop(t *testing.T) {
	const (
		TestString      string = "\tThis\tis\t\ta\ttest  \tstring"
		TestTabStopSize int    = 3
		Expected        string = "   This  is    a  test     string"
	)

	res, _ := FixTabStop(0, TestTabStopSize, " ", TestString)
	if res != Expected {
		t.Errorf("Expected %s, but got %s", Expected, res)
	}
}
