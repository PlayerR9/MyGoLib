package Go

import (
	"testing"
)

func TestMakeVariableName(t *testing.T) {
	res, err := MakeVariableName("TreeNode")
	if err != nil {
		t.Errorf("MakeVariableName failed: %s", err.Error())
	}

	if res != "tn" {
		t.Errorf("MakeVariableName failed: expected %s, got %s", "tn", res)
	}
}

func TestIsValidName(t *testing.T) {
	ok := IsValidName("tn", []string{"child"})
	if !ok {
		t.Errorf("IsValidName failed")
	}
}

func TestFixVariableName(t *testing.T) {
	res, ok := FixVariableName("tn", []string{"child"}, 1, "_")
	if !ok {
		t.Errorf("FixVariableName failed")
	}

	if res != "tn" {
		t.Errorf("FixVariableName failed: expected %s, got %s", "tn", res)
	}
}

func TestFixVarNameIncremental(t *testing.T) {
	res, ok := FixVarNameIncremental("tn", []string{"child"}, 1, 1)
	if !ok {
		t.Errorf("FixVarNameIncremental failed")
	}

	if res != "tn1" {
		t.Errorf("FixVarNameIncremental failed: expected %s, got %s", "tn1", res)
	}
}
