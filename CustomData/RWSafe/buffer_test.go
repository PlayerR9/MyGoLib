package RWSafe

import (
	"testing"
)

func TestInit(t *testing.T) {
	buffer := new(Buffer[int])
	sendTo, receiveFrom := buffer.Init(1)
	defer buffer.Wait()

	sendTo <- 1
	sendTo <- 2
	sendTo <- 3

	close(sendTo)

	for i := 1; i <= 3; i++ {
		x := <-receiveFrom

		if x != i {
			t.Errorf("Expected %v, got %v", i, x)
		}
	}
}

func TestTrimFrom(t *testing.T) {
	buffer := new(Buffer[int])
	sendTo, receiveFrom := buffer.Init(0)
	defer buffer.Wait()

	sendTo <- 1
	sendTo <- 2
	sendTo <- 3

	close(sendTo)

	buffer.CleanBuffer()

	x, ok := <-receiveFrom

	if !ok {
		t.Errorf("Expected %v, got %v", true, ok)
	}

	if x != 1 {
		t.Errorf("Expected %v, got %v", 1, x)
	}

	_, ok = <-receiveFrom

	if ok {
		t.Errorf("Expected %v, got %v", false, ok)
	}
}
