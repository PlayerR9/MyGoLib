package go_generator

import (
	"testing"
)

func TestIsValidName(t *testing.T) {
	err := IsValidName("tn", []string{"child"}, NotExported)
	if err != nil {
		t.Errorf("IsValidName failed: %s", err.Error())
	}
}
