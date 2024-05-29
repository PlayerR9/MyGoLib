package Debugger

import (
	"fmt"
	"log"
	"sync"

	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
)

// Debugger is a struct that provides a way to print debug messages.
type Debugger struct {
	// logger is the logger to use.
	logger *log.Logger

	// msgBuffer is the buffer for messages.
	msgBuffer *sfb.Buffer[PrintMessager]

	// receiveChan is the channel to receive messages.
	receiveChan <-chan PrintMessager

	// sendChan is the channel to send messages.
	sendChan chan<- PrintMessager

	// DebugMode is a flag that determines whether or not to print debug messages.
	debugMode bool

	// wg is the wait group for the goroutines.
	wg sync.WaitGroup

	// once is a flag to ensure that the debugger is only started once.
	once sync.Once
}

// NewDebugger is a function that creates a new debugger.
//
// Parameters:
//   - logger: The logger to use.
//
// Returns:
//   - *Debugger: The new debugger.
func NewDebugger(logger *log.Logger) *Debugger {
	dbg := &Debugger{
		logger: logger,
	}

	buffer, err := sfb.NewBuffer[PrintMessager](1)
	if err != nil {
		panic(err)
	}

	dbg.msgBuffer = buffer
	dbg.receiveChan = buffer.GetReceiveChannel()
	dbg.sendChan = buffer.GetSendChannel()

	return dbg
}

// ToggleDebugMode is a function that toggles the debug mode.
//
// Parameters:
//   - active: The flag to set the debug mode.
func (d *Debugger) ToggleDebugMode(active bool) {
	d.debugMode = active
}

// loggerListener is a function that listens for messages and logs them.
func (d *Debugger) loggerListener() {
	defer d.wg.Done()

	for msg := range d.receiveChan {
		switch msg := msg.(type) {
		case *PrintMessage:
			d.logger.Print(msg.v...)
		case *PrintfMessage:
			d.logger.Printf(msg.format, msg.v...)
		case *PrintlnMessage:
			d.logger.Println(msg.v...)
		default:
			d.logger.Printf("unknown message type: %T", msg)
		}
	}
}

// stdoutListener is a function that listens for messages and prints them to stdout.
func (d *Debugger) stdoutListener() {
	defer d.wg.Done()

	for msg := range d.receiveChan {
		switch msg := msg.(type) {
		case *PrintMessage:
			fmt.Print(msg.v...)
		case *PrintfMessage:
			fmt.Printf(msg.format, msg.v...)
		case *PrintlnMessage:
			fmt.Println(msg.v...)
		default:
			fmt.Printf("unknown message type: %T", msg)
		}
	}
}

// Start is a function that starts the debugger.
func (d *Debugger) Start() {
	d.once.Do(func() {
		d.msgBuffer.Start()

		d.wg.Add(1)

		if d.logger == nil {
			go d.stdoutListener()
		} else {
			go d.loggerListener()
		}
	})
}

// Close is a function that closes the debugger.
func (d *Debugger) Close() {
	if d.msgBuffer == nil {
		// Already closed
		return
	}

	// Close the send channel to signal the goroutine to stop
	close(d.sendChan)

	d.msgBuffer.Wait()

	d.wg.Wait()

	// Clean up
	d.msgBuffer = nil
	d.receiveChan = nil
	d.sendChan = nil

	d.logger = nil
}

// Wait is a function that waits for the debugger to finish.
func (d *Debugger) Wait() {
	d.wg.Wait()
}

// Println is a function that prints a line.
//
// Parameters:
//   - v: The values to print.
func (d *Debugger) Println(v ...interface{}) {
	if !d.debugMode || d.sendChan == nil {
		return
	}

	d.sendChan <- NewPrintlnMessage(v)
}

// Printf is a function that prints formatted text.
//
// Parameters:
//   - format: The format string.
//   - v: The values to print.
func (d *Debugger) Printf(format string, v ...interface{}) {
	if !d.debugMode || d.sendChan == nil {
		return
	}

	d.sendChan <- NewPrintfMessage(format, v)
}

// Print is a function that prints text.
//
// Parameters:
//   - v: The values to print.
func (d *Debugger) Print(v ...interface{}) {
	if !d.debugMode || d.sendChan == nil {
		return
	}

	d.sendChan <- NewPrintMessage(v)
}

// Write is a function that writes to the debugger.
//
// Parameters:
//   - p: The bytes to write.
//
// Returns:
//   - int: Always the length of the bytes.
//   - error: Always nil.
func (d *Debugger) Write(p []byte) (n int, err error) {
	if !d.debugMode || d.sendChan == nil {
		return 0, nil
	}

	d.sendChan <- NewPrintMessage([]any{p})

	return len(p), nil
}
