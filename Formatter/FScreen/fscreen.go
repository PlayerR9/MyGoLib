package FScreen

import (
	"fmt"
	"strings"
	"sync"
	"time"

	rws "github.com/PlayerR9/MyGoLib/CustomData/RWSafe"
	"github.com/gdamore/tcell"
)

type wScreen struct {
	title, currentProcess rws.RWSafe[string]
	counters              []*ScreenCounter

	mb *MessageBox

	screenWidth  rws.RWSafe[int]
	screenHeight rws.RWSafe[int]

	closeSignal rws.RWSafe[bool]
	hasPrinted  sync.Cond
	isAdding    rws.RWSafe[bool]
	canAdd      sync.Cond
	isShifting  rws.RWSafe[bool]
	canShiftUp  sync.Cond
	canPrint    sync.Cond
}

func (ws *wScreen) AddElement(contents []string, style tcell.Style) {
	ws.canAdd.L.Lock()
	for ws.isAdding.Get() && !ws.closeSignal.Get() {
		ws.canAdd.Wait()
	}
	defer ws.canAdd.L.Unlock()

	if ws.closeSignal.Get() {
		return
	}

	ws.isAdding.Set(true)

	remaining := ws.mb.EnqueueContents(contents, style, false)

	for len(remaining) > 0 {
		for ws.mb.firstEmptyLine >= ws.mb.height && !ws.closeSignal.Get() {
			ws.canShiftUp.Wait() // FIXME: This is not working
		}

		remaining = ws.mb.EnqueueContents(remaining, style, true)
	}

	ws.isAdding.Set(false)
}

func DrawScreen(title, currentProcess string, counters []*ScreenCounter, mb *MessageBox) (int, int, int, int) {
	DrawHorizontalBorder := func(y, width int) {
		screen.SetContent(0, y, '+', nil, NormalStyle) // Left corner

		for x := 1; x < width-1; x++ {
			screen.SetContent(x, y, '-', nil, NormalStyle)
		}

		screen.SetContent(width-1, y, '+', nil, NormalStyle) // Right corner
	}

	// Initialize variables
	width, height := screen.Size()
	innerWidth := width - 4   // Padding
	innerHeight := height - 2 // Padding

	y := 0 // y coordinate

	// Print the title
	startPos := (innerWidth - len(title)) / 2
	for x, char := range title {
		screen.SetContent(startPos+x, y, char, nil, NormalStyle)
	}

	y += 2

	// Print the current process (if any)
	if currentProcess != "" {
		for x, char := range currentProcess {
			screen.SetContent(x, y, char, nil, NormalStyle)
		}

		y += 2
	}

	// Print all the counters (if any)
	if len(counters) != 0 {
		for _, counter := range counters {
			for x, char := range counter.String() {
				screen.SetContent(x, y, char, nil, NormalStyle)
			}

			y++
		}

		y++
	}

	// Draw the message box
	DrawHorizontalBorder(y, width) // Top border
	y++

	for ; y < innerHeight; y++ {
		screen.SetContent(0, y, '|', nil, NormalStyle) // Left border
		mb.DrawLine(1, y)
		screen.SetContent(width-1, y, '|', nil, NormalStyle) // Right border
	}

	DrawHorizontalBorder(y, width) // Bottom border

	return width, height, innerWidth, innerHeight
}

func (ws *wScreen) PrintScreen() {
	screen.Clear()

	DrawScreen(
		ws.title.Get(),
		ws.currentProcess.Get(),
		ws.counters,
		ws.mb,
	)

	// Print the screen
	screen.Show()

	if ws.mb.CanShiftUp() {
		ws.canShiftUp.Signal()
	}

	ws.hasPrinted.Signal()

	if ws.isAdding.Get() {

	} else {
		// Wait a bit before printing the next screen
		time.Sleep(100 * time.Millisecond)
	}
}

const (
	BufferSize int = 100
)

// Styles
var (
	// Normal
	NormalStyle tcell.Style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGhostWhite)

	// Debug
	DebugStyle tcell.Style = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorSlateGray).Background(tcell.ColorGhostWhite)

	// Warning
	WarningStyle tcell.Style = tcell.StyleDefault.Foreground(tcell.ColorDarkOrange).Background(tcell.ColorGhostWhite)

	// Error
	ErrorStyle tcell.Style = tcell.StyleDefault.Foreground(tcell.ColorFireBrick).Background(tcell.ColorGhostWhite)

	// Fatal
	FatalStyle tcell.Style = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDarkRed).Background(tcell.ColorGhostWhite)

	// Success
	SuccessStyle tcell.Style = tcell.StyleDefault.Foreground(tcell.ColorDarkGreen).Background(tcell.ColorGhostWhite)
)

type PrintCondition int

const (
	JustPrinted PrintCondition = iota
	YetToPrint
	NotPrinted
	MustPrint
)

var (
	fscreen        *wScreen
	screen         tcell.Screen
	receiveChannel chan Messager
	isInitialized  bool           = false
	isRunning      bool           = false
	shouldPrint    PrintCondition = NotPrinted
	wgFScreen      sync.WaitGroup
)

func Initialize(titleScreen string) error {
	if isRunning {
		SendMessages(
			MessageWarning(
				"FScreen is running; ignoring call to Initialize()",
			),
		)

		return nil
	} else if isInitialized {
		return fmt.Errorf("FScreen is already initialized")
	} else if strings.TrimSpace(titleScreen) == "" {
		return fmt.Errorf("titleScreen cannot be empty")
	}

	var err error

	screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}

	err = screen.Init()
	if err != nil {
		return err
	}
	screen.Clear()

	receiveChannel = make(chan Messager, BufferSize)

	width, height := screen.Size()

	fscreen = new(wScreen)
	fscreen.title = rws.NewRWSafe(titleScreen)
	fscreen.currentProcess = rws.NewRWSafe("")
	fscreen.counters = make([]*ScreenCounter, 0)
	fscreen.table = make([][]*rws.RWSafe[rune], 0)
	fscreen.styles = make([]*rws.RWSafe[tcell.Style], 0)
	fscreen.maxLines = height - 2
	fscreen.firstEmptyLine = rws.NewRWSafe(0)
	fscreen.screenWidth = rws.NewRWSafe(width)
	fscreen.screenHeight = rws.NewRWSafe(height)
	fscreen.innerScreenWidth = rws.NewRWSafe(width - 4)
	fscreen.closeSignal = rws.NewRWSafe(false)
	fscreen.hasPrinted = *sync.NewCond(new(sync.Mutex))
	fscreen.isAdding = rws.NewRWSafe(false)
	fscreen.canAdd = *sync.NewCond(new(sync.Mutex))
	fscreen.isShifting = rws.NewRWSafe(false)
	fscreen.canShiftUp = *sync.NewCond(new(sync.Mutex))
	fscreen.canPrint = *sync.NewCond(new(sync.Mutex))

	isInitialized = true

	return nil
}

func Run() {
	if isRunning {
		SendMessages(
			MessageWarning(
				"FScreen is already running; ignoring call to Run()",
			),
		)

		return
	} else if !isInitialized {
		panic("FScreen is not initialized")
	}

	var wg sync.WaitGroup

	for msg := range receiveChannel {
		msg.Apply()
		currentStyle := msg.GetStyle()

		if shouldPrint == MustPrint {
			wg.Wait() // Wait for the previous screen to be printed

			wg.Add(1)

			go func() {
				fscreen.PrintScreen()
				wg.Done()
			}()
		}
	}

	wg.Wait() // Wait for the last screen to be printed
}

func SendMessages(msg Messager, optionals ...Messager) {
	if !isRunning {
		panic("FScreen is not running")
	}

	receiveChannel <- msg

	for _, m := range optionals {
		receiveChannel <- m
	}
}

//// OLD CODE ////

const (
	Padding      = 2
	PaddingWidth = 4 // 2 * Padding

	Hellip    = "..."
	HellipLen = 3

	WaitTime = 100 * time.Millisecond
)

var (
	innerScreenWidth int = screenWidth - PaddingWidth
)

var (
	separator string = ""
	lineBreak string = ""
)

var (
	emptyLine int = 1
)

func Initialize(name string) error {
	if isRunning {
		SendMessages(
			MessageFatal(
				"FScreen is running.",
				"Stop before initializing",
			),
		)

		return nil
	}

	var err error

	screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}
	if err = screen.Init(); err != nil {
		return err
	}
	screen.Clear()

	title = name

	isInitialized = true

	title = name

	isInitialized = true

	separator = strings.Repeat(Separator, screenWidth-Padding)
	lineBreak = strings.Repeat(Space, screenWidth-Padding)

	return nil
}

func lineInsert(line string, trimIndex int, baseIndex int) int {
	shiftScreenUp()

	index := baseIndex

	if trimIndex == -1 {
		copy(screen[emptyLine][index:], []rune(line))

		index += len(line)
	} else {
		copy(screen[emptyLine][index:], []rune(line[:trimIndex]))
		copy(screen[emptyLine][screenWidth-5:], []rune(Hellip))
	}

	return index
}

func Start() {
	if !isInitialized {
		panic("FScreen is not initialized")
	} else if isRunning {
		SendMessages(
			MessageWarning(
				"FScreen is already running",
				"Ignoring call to Run()",
			),
		)

		return
	}

	isRunning = true

	wg.Add(1)

	go func() {
		run()
		wg.Done()
	}()
}

func Wait() {
	if !isRunning {
		return
	}

	wg.Wait()
}

func run() {
	printScreen()

	shouldPrint = JustPrinted

	for message := range messageChannel {
		if _, ok := message.(EmptyMSG); ok {
			continue
		}

		message.Apply()

		if shouldPrint == MustPrint {
			printScreen()

			shouldPrint = JustPrinted
		}

		switch message.(type) {
		case ImportantMSG:
			// Wait a bit before displaying the next message
			time.Sleep(WaitTime)
		case CloseMSG, FatalMSG:
			isRunning = false
			return
		}
	}

	if shouldPrint == YetToPrint {
		printScreen()
	}
}
