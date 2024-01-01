package FScreen

import (
	"fmt"
	"strings"
	"sync"
	"time"

	rws "github.com/PlayerR9/MyGoLib/CustomData/RWSafe"
	"github.com/gdamore/tcell"
)

const (
	IndentLevel int = 2
)

type wScreen struct {
	title, currentProcess *rws.RWSafe[string]
	counters              []*ScreenCounter

	table         [][]*rws.RWSafe[rune]
	styles        []*rws.RWSafe[tcell.Style]
	maxLines      int
	lastEmptyLine rws.RWSafe[int]

	// Screen Information
	screenWidth      rws.RWSafe[int]
	screenHeight     rws.RWSafe[int]
	innerScreenWidth rws.RWSafe[int]

	closeSignal rws.RWSafe[bool]
	hasPrinted  sync.Cond
	isAdding    rws.RWSafe[bool]
	canAdd      sync.Cond
	isShifting  rws.RWSafe[bool]
	canShiftUp  sync.Cond
	canPrint    sync.Cond
}

// Returns the first empty index after the last character
func (ws *wScreen) writeWord(x, y int, word string) int {
	row := ws.table[y]

	for _, char := range word {
		row[x].Set(char)
		x++
	}

	return x
}

// Returns the first empty index after the last character
func (ws *wScreen) canWriteWord(x int, word string) bool {
	return x+len(word) <= ws.innerScreenWidth.Get()
}

func (ws *wScreen) WriteHorizontalBorder(y int) {
	screen.SetContent(0, y, '+', nil, NormalStyle) // Left corner

	for x := 1; x < ws.screenWidth.Get()-1; x++ {
		screen.SetContent(x, y, '-', nil, NormalStyle)
	}

	screen.SetContent(ws.screenWidth.Get()-1, y, '+', nil, NormalStyle) // Right corner
}

func (ws *wScreen) writeFields(x, y int, fields []string) (int, int) {
	if len(fields) == 1 {
		if ws.canWriteWord(x, fields[0]) {
			return ws.writeWord(x, y, fields[0]), y
		}

		x = ws.writeWord(x, y, fields[0][:ws.innerScreenWidth.Get()-x-HellipLen])
		return ws.writeWord(x, y, Hellip), y
	} else if !ws.canWriteWord(x, fields[0]) {
		x = ws.writeWord(x, y, fields[0][:ws.innerScreenWidth.Get()-x-HellipLen])
		x = ws.writeWord(x, y, Hellip)

		return IndentLevel, y + 1
	}

	x = ws.writeWord(x, y, fields[0])

	index := -1
	for i, field := range fields[1:] {
		if !ws.canWriteWord(x+2, field) {
			index = i
			break
		}

		x = ws.writeWord(x+2, y, field)
	}

	if index != -1 {
		ws.writeWord(x, y, Hellip)

		return IndentLevel, y + 1
	}

	return x, y
}

// -1 = no fields, >= 0 index of the last field that fits
func (ws *wScreen) canWriteFields(x int, fields []string) int {
	if !ws.canWriteWord(x, fields[0]) {
		return -1
	}

	x += len(fields[0])

	for i, field := range fields[1:] {
		if !ws.canWriteWord(x+2, field) {
			return i
		}

		x += len(field) + 2
	}

	return len(fields) - 1
}

func (ws *wScreen) WriteElement(yCoord int, contents []string, style tcell.Style, isIndented bool) (int, []string) {
	if yCoord < 1 || yCoord >= ws.screenHeight.Get()-1 {
		panic(fmt.Sprintf("y must be between 1 and %d", ws.screenHeight.Get()-1))
	}

	if yCoord >= len(ws.table) {
		for i := len(ws.table); i < yCoord; i++ {
			row := make([]*rws.RWSafe[rune], ws.innerScreenWidth.Get())

			for j := 0; j < ws.innerScreenWidth.Get(); j++ {
				row[j] = new(rws.RWSafe[rune])
				row[j].Set(' ')
			}

			ws.table = append(ws.table, row)

			style := new(rws.RWSafe[tcell.Style])
			style.Set(NormalStyle)

			ws.styles = append(ws.styles, style)
		}
	}

	fields := strings.Fields(contents[0])
	if len(fields) == 0 {
		return yCoord, nil
	}

	xCoord := 0
	if isIndented {
		xCoord = IndentLevel
	}

	var newYCoord int
	xCoord, newYCoord = ws.writeFields(xCoord, yCoord, fields)
	if newYCoord != yCoord {
		return newYCoord, contents[1:]
	}

	// Try to see if the next content fits on the same line
	for index, content := range contents[1:] {
		fields = strings.Fields(content)
		if len(fields) == 0 {
			continue // Skip empty lines
		}

		lastValidField := ws.canWriteFields(xCoord, fields)
		if lastValidField == -1 {
			return yCoord + 1, contents[index:]
		}

		for _, field := range fields[:lastValidField+1] {
			xCoord = ws.writeWord(xCoord+2, yCoord, field)
		}

		if lastValidField != len(fields)-1 {
			return yCoord + 1, append([]string{strings.Join(fields[lastValidField+1:], " ")}, contents[index+1:]...)
		}
	}

	return yCoord, nil
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

	yCoord, remaining := ws.WriteElement(ws.lastEmptyLine.Get(), contents, style, false)
	ws.lastEmptyLine.Set(yCoord)

	for len(remaining) > 0 {
		for ws.lastEmptyLine.Get() >= ws.maxLines && !ws.closeSignal.Get() {
			ws.canShiftUp.Wait() // FIXME: This is not working
		}

		yCoord, remaining = ws.WriteElement(ws.lastEmptyLine.Get(), remaining, style, true)
		ws.lastEmptyLine.Set(yCoord)
	}

	ws.isAdding.Set(false)
}

func (ws *wScreen) PrintScreen() {
	screen.Clear()

	// Initialize variables
	width, height := screen.Size()
	ws.screenWidth.Set(width)
	ws.screenHeight.Set(height)
	ws.innerScreenWidth.Set(width - 4) // Padding

	y := 0 // y coordinate

	// Print the title
	title := ws.title.Get()
	startPos := (ws.screenWidth.Get() - len(title)) / 2

	for x, char := range title {
		screen.SetContent(startPos+x, y, char, nil, NormalStyle)
	}

	y += 2

	// Print the current process (if any)
	currentProcess := ws.currentProcess.Get()

	if currentProcess != "" {
		for x, char := range currentProcess {
			screen.SetContent(x, y, char, nil, NormalStyle)
		}

		y += 2
	}

	// Print all the counters (if any)
	if len(ws.counters) != 0 {
		for _, counter := range ws.counters {
			for x, char := range counter.String() {
				screen.SetContent(x, y, char, nil, NormalStyle)
			}

			y++
		}

		y++
	}

	// Draw the box and the elements
	ws.WriteHorizontalBorder(y) // Top border
	y++

	// 2. Element
	for i, row := range ws.table {
		screen.SetContent(0, y, '|', nil, NormalStyle) // Left border

		styleOfRow := ws.styles[i].Get()
		for x, char := range row {
			screen.SetContent(x+2, y, char.Get(), nil, styleOfRow)
		}

		screen.SetContent(ws.screenWidth.Get()-1, y, '|', nil, NormalStyle) // Right border

		y++

		if i >= ws.lastEmptyLine.Get() {
			break
		}
	}

	// Fill the rest of the screen with empty lines
	for y < ws.maxLines {
		screen.SetContent(0, y, '|', nil, NormalStyle)                      // Left border
		screen.SetContent(ws.screenWidth.Get()-1, y, '|', nil, NormalStyle) // Right border

		y++
	}

	ws.WriteHorizontalBorder(y) // Bottom border

	// Print the screen
	screen.Show()

	if ws.CanShiftUp() {
		ws.canShiftUp.Signal()
	}

	ws.hasPrinted.Signal()

	if ws.isAdding.Get() {

	} else {
		// Wait a bit before printing the next screen
		time.Sleep(100 * time.Millisecond)
	}
}

func (ws *wScreen) CanShiftUp() bool {
	if ws.isShifting.Get() {
		return ws.lastEmptyLine.Get() != 0
	}

	return ws.lastEmptyLine.Get() > ws.maxLines/2
}

func (ws *wScreen) ShiftUp() {
	for {
		ws.canShiftUp.L.Lock()
		for !ws.CanShiftUp() && !ws.closeSignal.Get() {
			ws.canShiftUp.Wait()
		}

		if ws.closeSignal.Get() {
			return
		}

		ws.isShifting.Set(true)

		// Shift up
		limit := ws.lastEmptyLine.Get()

		for i := 0; i < ws.maxLines-limit; i++ {
			ws.table[i] = ws.table[i+limit]
			ws.styles[i] = ws.styles[i+limit]
		}

		ws.lastEmptyLine.Set(ws.maxLines - limit)

		ws.isShifting.Set(false)
		ws.canShiftUp.L.Unlock()

		// Wait some time before shifting up again
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

	fscreen = &wScreen{
		title:          &wElement{content: titleScreen, style: NormalStyle},
		currentProcess: &wElement{content: "", style: NormalStyle},
		counters:       make([]*ScreenCounter, 0),
		elements:       make([]*wElement, 0),
	}

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

func writeLines(contents []string) {
	innerScreenWidth, _ := screen.Size()
	innerScreenWidth -= 4 // Padding

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
