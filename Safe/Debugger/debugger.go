package Debugger

import (
	"fmt"
	"log"
	"sync"

	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
)

type Debugger struct {
	logger *log.Logger

	msgBuffer   *sfb.Buffer[*PrintMessage]
	receiveChan <-chan *PrintMessage
	sendChan    chan<- *PrintMessage

	// DebugMode is a flag that determines whether or not to print debug messages.
	debugMode bool

	wg sync.WaitGroup

	once sync.Once
}

func NewDebugger(logger *log.Logger) *Debugger {
	dbg := &Debugger{
		logger: logger,
	}

	buffer, err := sfb.NewBuffer[*PrintMessage](1)
	if err != nil {
		panic(err)
	}

	dbg.msgBuffer = buffer
	dbg.receiveChan = buffer.GetReceiveChannel()
	dbg.sendChan = buffer.GetSendChannel()

	return dbg
}

func (d *Debugger) ToggleDebugMode(active bool) {
	d.debugMode = active
}

func (d *Debugger) loggerListener() {
	defer d.wg.Done()

	for msg := range d.receiveChan {
		switch msg.pmType {
		case PM_Println:
			d.logger.Println(msg.v...)
		case PM_Printf:
			d.logger.Printf(msg.format, msg.v...)
		case PM_Print:
			d.logger.Print(msg.v...)
		}
	}
}

func (d *Debugger) stdoutListener() {
	defer d.wg.Done()

	for msg := range d.receiveChan {
		switch msg.pmType {
		case PM_Println:
			fmt.Println(msg.v...)
		case PM_Printf:
			fmt.Printf(msg.format, msg.v...)
		case PM_Print:
			fmt.Print(msg.v...)
		}
	}
}

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

func (d *Debugger) Wait() {
	d.wg.Wait()
}

func (d *Debugger) Println(v ...interface{}) {
	if !d.debugMode || d.sendChan == nil {
		return
	}

	d.sendChan <- NewPrintMessage(PM_Println, "", v...)
}

func (d *Debugger) Printf(format string, v ...interface{}) {
	if !d.debugMode || d.sendChan == nil {
		return
	}

	d.sendChan <- NewPrintMessage(PM_Printf, format, v...)
}

func (d *Debugger) Print(v ...interface{}) {
	if !d.debugMode || d.sendChan == nil {
		return
	}

	d.sendChan <- NewPrintMessage(PM_Print, "", v...)
}

func (d *Debugger) Write(p []byte) (n int, err error) {
	if !d.debugMode || d.sendChan == nil {
		return 0, nil
	}

	d.sendChan <- NewPrintMessage(PM_Print, "", p)

	return len(p), nil
}
