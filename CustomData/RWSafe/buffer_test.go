package RWSafe

import (
	"testing"
)

func TestInit(t *testing.T) {
	buffer := new(Buffer[int])

	if err := buffer.Init(1); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer buffer.Wait()

	sendTo, receiveFrom := buffer.GetSendChannel(), buffer.GetReceiveChannel()

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

	if err := buffer.Init(1); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer buffer.Wait()

	sendTo, receiveFrom := buffer.GetSendChannel(), buffer.GetReceiveChannel()

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
