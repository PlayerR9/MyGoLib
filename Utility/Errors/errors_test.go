package Errors

import (
	"fmt"
	"testing"
)

func TestNewErrCallFailed(t *testing.T) {
	err := NewErrCallFailed("NewErrCallFailed", NewErrCallFailed)
	if err != nil {
		fmt.Println(err.functionSignature)

		t.Errorf("NewErrCallFailed returned non-nil")
	}
}
