package FScreen

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

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
)

var (
	screenWidth  int
	screenHeight int
)

var (
	screen        tcell.Screen
	titleScreen   string
	elementScreen []string
)

var (
	receiveChannel chan Messager
	isInitialized  bool = false
	isRunning      bool = false
)

func Initialize(title string) error {
	if isRunning {
		SendMessages(
			MessageError(
				"FScreen is running.",
				"Stop before initializing",
			),
		)

		return nil
	} else if isInitialized {
		return fmt.Errorf("FScreen is already initialized")
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

	receiveChannel = make(chan Messager, BufferSize)

	titleScreen = title

	isInitialized = true

	elementScreen = make([]string, 0) // Set to the capacity of the screen

	return nil
}

func Run() error {
	if !isInitialized {
		return fmt.Errorf("FScreen is not initialized")
	} else if isRunning {
		SendMessages(
			MessageWarning(
				"FScreen is already running",
				"Ignoring call to Run()",
			),
		)

		return nil
	}

	for msg := range receiveChannel {
		msg.Apply()

		currentStyle := msg.GetStyle()
	}

	return nil
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

func printScreen() {
	// Clear the screen
	screen.Clear()

	j := 0

	screenWidth, screenHeight = screen.Size()

	// Print the title

	// 1. Found middle of the line
	startPos := (screenWidth - len(titleScreen)) / 2

	// 2. Print the title
	for i, char := range titleScreen {
		screen.SetContent(i+startPos, j, rune(char), nil, NormalStyle)
	}

	j++

	// Print the current process

	if currentProcess != "" {
		j++

		for i, char := range currentProcess {
			screen.SetContent(i, j, rune(char), nil, NormalStyle)
		}
	}

	// Print the counters

	j++

	for j, counter := range counters {
		for i, char := range counter.String() {
			screen.SetContent(i, j, rune(char), nil, NormalStyle)
		}

		j++
	}

	// Draw the box and the elements

	// 1. Corners
	screen[0][0] = '+'
	screen[0][screenWidth-1] = '+'
	screen[screenHeight-1][0] = '+'
	screen[screenHeight-1][screenWidth-1] = '+'

	// 2. Top and bottom borders
	for j := 1; j < screenWidth-1; j++ {
		screen[0][j] = '-'
		screen[screenHeight-1][j] = '-'
	}

	// 3. Left and right borders
	for i := 1; i < screenHeight-1; i++ {
		screen[i][0] = '|'
		screen[i][screenWidth-1] = '|'
	}

	// Clear the screen
	for i := 1; i < screenHeight-1; i++ {
		for j := 1; j < screenWidth-1; j++ {
			screen[i][j] = ' '
		}
	}
}

//// OLD CODE ////

const (
	Padding      = 2
	PaddingWidth = 4 // 2 * Padding

	Hellip    = "..."
	HellipLen = 3

	Space     = " "
	Separator = "="

	WaitTime = 100 * time.Millisecond
)

func SetShiftUp(shift int) error {
	if isRunning {
		SendMessages(
			MessageError(
				"FScreen is running.",
				"Stop before changing the shift up",
			),
		)

		return nil
	}

	if shift < 1 {
		return fmt.Errorf("shift up must be at least 1. Got %d instead", shift)
	} else if shift > screenHeight-Padding {
		return fmt.Errorf("shift up must be at most %d. Got %d instead", screenHeight-Padding, shift)
	}

	shiftUp = shift

	isInitialized = false

	return nil
}

var (
	screenHeight     int = 10
	screenWidth      int = 80
	innerScreenWidth int = screenWidth - PaddingWidth
)

var (
	title          string = ""
	currentProcess string = ""
	counters       []ScreenCounter
	separator      string = ""
	lineBreak      string = ""
)

var (
	emptyLine int = 1
	shiftUp   int = 5
)

var (
	messageChannel chan Messager  = make(chan Messager, BufferSize)
	shouldPrint    PrintCondition = NotPrinted
	wg             sync.WaitGroup
)

type PrintCondition int

const (
	JustPrinted PrintCondition = iota
	YetToPrint
	NotPrinted
	MustPrint
)

func Defaults() {
	if isRunning {
		SendMessages(
			MessageFatal(
				"FScreen is running.",
				"Stop before initializing",
			),
		)

		return
	}

	screenHeight = 10
	screenWidth = 80
	innerScreenWidth = screenWidth - PaddingWidth

	shiftUp = 5

	isInitialized = false
}

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

func clearScreen() {

}

func printScreen() {
	screen.Clear()

	screenWidth, _ := screen.Size()

	startPos := (screenWidth - len(title)) / 2

	// Print the title
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	for i, char := range title {
		screen.SetContent(i+startPos, 0, rune(char), nil, style)
	}

	// Print the current process
	if currentProcess != "" {
		for i, char := range currentProcess {
			screen.SetContent(i, 1, rune(char), nil, style)
		}
	}

	// Print the counters
	for j, counter := range counters {
		for i, char := range counter.String() {
			screen.SetContent(i, 2+j, rune(char), nil, style)
		}
	}

	// Draw the box

	// Print the screen
	for i := 0; i < screenHeight; i++ {
		for j := 0; j < screenWidth; j++ {
			fmt.Print(string(screen[i][j]))
		}

		fmt.Println()
	}

	screen.Show()
}

func shiftScreenUp() {
	if emptyLine <= screenHeight-Padding {
		return
	}

	for i := 1; i < screenHeight-shiftUp-1; i++ {
		screen[i] = screen[i+shiftUp-1]
	}

	for i := screenHeight - shiftUp - 1; i < screenHeight-1; i++ {
		for j := 1; j < screenWidth-1; j++ {
			screen[i][j] = ' '
		}

		screen[i][0] = '|'
		screen[i][screenWidth-1] = '|'
	}

	emptyLine -= shiftUp
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

func writeLines(contents []string) {
	trimIndex := make([]int, len(contents))
	for i, content := range contents {
		if len(content) > innerScreenWidth {
			trimIndex[i] = strings.LastIndex(content[:innerScreenWidth-HellipLen], " ")
		} else {
			trimIndex[i] = -1
		}
	}

	index := lineInsert(contents[0], trimIndex[0], 2)

	for i, content := range contents[1:] {
		if len(content)+index+1 <= innerScreenWidth {
			// Add a space
			screen[emptyLine][index] = ' '
			index++

			// Copy the content
			copy(screen[emptyLine][index:], []rune(content))

			index += len(content)
		} else {
			emptyLine++

			index = lineInsert(content, trimIndex[i], 4)
		}
	}

	emptyLine++

	shouldPrint = MustPrint
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
