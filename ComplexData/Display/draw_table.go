package Display

import (
	"fmt"
	"strings"

	"github.com/PlayerR9/MyGoLib/ListLike/Stack"
	"github.com/gdamore/tcell"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// dtCell represents a cell in a data table. It contains a rune which
// represents the content of the cell, and a tcell.Style which represents
// the style of the cell.
type DtCell struct {
	// Content is the rune that represents the content of the cell.
	Content rune

	// Style is the tcell.Style that represents the style of the cell.
	Style tcell.Style
}

func (cell *DtCell) String() string {
	return string(cell.Content)
}

func NewDtCell(content rune, style tcell.Style) *DtCell {
	return &DtCell{
		Content: content,
		Style:   style,
	}
}

type DrawTable struct {
	table         [][]*DtCell
	width, height int
}

func (dt *DrawTable) GetWidth() int {
	return dt.width
}

func (dt *DrawTable) GetHeight() int {
	return dt.height
}

func (dt *DrawTable) WriteAt(x, y int, cell *DtCell) error {
	if x < 0 || x >= dt.width {
		return ers.NewErrInvalidParameter("x").
			WithReason(ers.NewErrOutOfBound(0, dt.width, x))
	}

	if y < 0 || y >= dt.height {
		return ers.NewErrInvalidParameter("y").
			WithReason(ers.NewErrOutOfBound(0, dt.height, y))
	}

	dt.data.table[y+dt.absoluteY][x+dt.absoluteX] = cell

	return nil
}

func (dt *DrawTable) GetAt(x, y int) (*DtCell, error) {
	if x < 0 || x >= dt.width {
		return nil, ers.NewErrInvalidParameter("x").
			WithReason(ers.NewErrOutOfBound(0, dt.width, x))
	}

	if y < 0 || y >= dt.height {
		return nil, ers.NewErrInvalidParameter("y").
			WithReason(ers.NewErrOutOfBound(0, dt.height, y))
	}

	return dt.data.table[y+dt.absoluteY][x+dt.absoluteX], nil
}

func (dt *DrawTable) ClearTable() {
	startY := dt.absoluteY

	for i := 0; i < dt.height; i++ {
		dt.data.table[i+startY] = make([]*DtCell, dt.width)
	}
}

func (dt *DrawTable) WriteVerticalSequence(x, y int, sequence []*DtCell) error {
	if len(sequence) == 0 {
		return nil
	}

	if x < 0 || x >= dt.width {
		return ers.NewErrInvalidParameter("x").
			WithReason(ers.NewErrOutOfBound(0, dt.width, x))
	}

	if y < 0 || y >= dt.height {
		return ers.NewErrInvalidParameter("y").
			WithReason(ers.NewErrOutOfBound(0, dt.height, y))
	}

	if y+len(sequence) > dt.height {
		return ers.NewErrInvalidParameter("sequence").
			WithReason(fmt.Errorf("sequence length (%d) exceeds the height (%d)", len(sequence), dt.height-y))
	}

	x += dt.absoluteX
	y += dt.absoluteY

	for i, cell := range sequence {
		dt.data.table[y+i][x] = cell
	}

	return nil
}

func (dt *DrawTable) WriteHorizontalSequence(x, y int, sequence []*DtCell) error {
	if len(sequence) == 0 {
		return nil
	}

	if x < 0 || x >= dt.width {
		return ers.NewErrInvalidParameter("x").
			WithReason(ers.NewErrOutOfBound(0, dt.width, x))
	}

	if y < 0 || y >= dt.height {
		return ers.NewErrInvalidParameter("y").
			WithReason(ers.NewErrOutOfBound(0, dt.height, y))
	}

	if x+len(sequence) > dt.width {
		return ers.NewErrInvalidParameter("sequence").
			WithReason(fmt.Errorf("sequence length (%di) exceeds the width (%di)", len(sequence), dt.width-x))
	}

	x += dt.absoluteX
	y += dt.absoluteY

	for j, cell := range sequence {
		dt.data.table[y][x+j] = cell
	}

	return nil
}

// GetLines generates and returns a slice of strings representing
// the content of the drawTable. Each string in the slice corresponds
// to a line in the drawTable, and each character in a string
// corresponds to a cell in the line. If a cell contains a rune,
// the rune is included in the string. If a cell is nil, a space
// character is included in the string.
//
// The function returns a slice of strings. The length of the slice
// is equal to the height of the drawTable, and the length of each
// string is equal to the width of the drawTable.
func (dt *DrawTable) GetLines() []string {
	lines := make([]string, 0, dt.height)
	var builder strings.Builder

	x := dt.absoluteX
	y := dt.absoluteY

	for i := 0; i < dt.height; i++ {
		for j := 0; j < dt.width; j++ {
			if cell := dt.data.table[i+y][j+x]; cell != nil {
				fmt.Fprintf(&builder, "%v", cell)
			} else {
				builder.WriteRune(' ')
			}
		}

		lines = append(lines, builder.String())
		builder.Reset()
	}

	return lines
}

func (dt *DrawTable) ExecuteCells(f func(int, int, *DtCell) error) error {
	x := dt.absoluteX
	y := dt.absoluteY

	for i := 0; i < dt.height; i++ {
		for j := 0; j < dt.width; j++ {
			err := f(j, i, dt.data.table[y+i][x+j])
			if err != nil {
				return fmt.Errorf("failed to execute function on cell (%d, %d): %v", j, i, err)
			}
		}
	}

	return nil
}

func InitializePartition[DtCell PtCell](width, height int) *DrawTable {
	rt := &rootTable[DtCell]{
		width:  width,
		height: height,
		table:  make([][]*DtCell, height),
	}

	for i := 0; i < height; i++ {
		rt.table[i] = make([]*DtCell, width)
	}

	return &DrawTable{
		data:     rt,
		parent:   nil,
		children: make([]*DrawTable, 0),

		width:     width,
		height:    height,
		startX:    0,
		startY:    0,
		absoluteX: 0,
		absoluteY: 0,
	}
}

func Deallocate[DtCell PtCell](target *DrawTable) error {
	// 1. Check parameters
	if target == nil {
		return nil // Nothing to do
	}

	if target.parent == nil {
		return ers.NewErrInvalidParameter("target").
			WithReason(fmt.Errorf("root cannot be deallocated"))
	}

	// 2. Use a stack to deallocate the target
	S := Stack.NewArrayStack(target)

	for !S.IsEmpty() {
		node := S.Peek()

		// 1. Check if the node has children
		if node.children != nil {
			// Evaluate the children first
			// by pushing them into the stack above
			// the current node

			for i := 0; i < len(node.children); i++ {
				S.Push(node.children[i])
				node.children[i] = nil
			}

			node.children = nil

			continue
		}

		// Otherwise, remove the node from the stack
		S.Pop()

		// Clear the node
		node.data = nil
		node.parent = nil

		// Call to cleanup() to the element
		err := node.element.Cleanup()
		if err != nil {
			return fmt.Errorf("could not deallocate element: %v", err)
		}
	}

	// 3. Call to compress() to the parent in order to
	// optimize the memory
	_, err := compress(target.parent)
	if err != nil {
		return fmt.Errorf("could not compress table: %v", err)
	}

	// FIXME: Remember to force draw the highest modified

	if err := target.element.Cleanup(); err != nil {
		return fmt.Errorf("could not cleanup element: %v", err)
	}

	return nil
}

func (dt *DrawTable) ShiftUp() (*DrawTable, error) {
	if len(dt.children) == 0 {
		return nil, fmt.Errorf("cannot shift up: no children")
	}

	// Get the first child
	firstChild := dt.children[0]

	// Shift the child to the up
	var prevChild *DrawTable

	for i, child := range dt.children[1:] {
		prevChild = dt.children[i-1]

		// FIX ME: Remember to shift not by the
		// slice order but by the actual position

		prevChild.data = child.data
		prevChild.width = child.width
		prevChild.height = child.height
		prevChild.startX = child.startX
		prevChild.startY = child.startY
	}

	// FIX ME: Remember to get back the memory
	// of the first child

	return firstChild, nil
}

func (dt *DrawTable) GetPathToRoot() []*DrawTable {
	path := make([]*DrawTable, 0)

	for node := dt; node != nil; node = node.parent {
		path = append([]*DrawTable{node}, path...)
	}

	return path
}

func GetHighestCommonParent[DtCell PtCell](node1, node2 *DrawTable) *DrawTable {
	if node1 == nil || node2 == nil {
		return nil
	}

	path1, path2 := node1.GetPathToRoot(), node2.GetPathToRoot()

	var commonParent *DrawTable

	for i := 0; i < len(path1) && i < len(path2); i++ {
		if path1[i] != path2[i] {
			break
		}

		commonParent = path1[i]
	}

	return commonParent
}

type DrawTable struct {
	tableHeight int // Height of the table
	*p.PtTable[*DtCell]
}

func (dt *DrawTable) Init(dt p.Partitionable[*DtCell]) *DrawTable {
	dt.tableHeight = dt.GetHeight()
	dt.PtTable = new(p.PtTable[*DtCell]).Init(dt)
	dt.untriggeredEvent = nil

	return dt
}

/// FIX THIS

func (dt *DrawTable) HasEvent() bool {
	return dt.untriggeredEvent != nil
}

func (dt *DrawTable) SetEvent(ev *DtEvent) {
	dt.untriggeredEvent = ev
}

func (dt *DrawTable) WriteLineAt(x, y int, contents string, style tcell.Style) error {
	if len(contents) == 0 {
		return nil // Nothing to do
	}

	width, height := dt.GetWidth(), dt.GetHeight()

	if x < 0 || x >= width {
		return ers.NewErrInvalidParameter("x").
			WithReason(ers.NewErrOutOfBound(0, width, x))
	}

	if y < 0 || y >= height {
		return ers.NewErrInvalidParameter("y").
			WithReason(ers.NewErrOutOfBound(0, height, y))
	}

	if x+len(contents) > width {
		return ers.NewErrInvalidParameter("contents").
			WithReason(fmt.Errorf("length (%d) exceeds table's width by %d", len(contents), x+len(contents)-width))
	}

	sequence := make([]*DtCell, 0, len(contents))
	for _, content := range contents {
		sequence = append(sequence, &DtCell{Content: content, Style: style})
	}

	dt.WriteSequence(x, y, sequence, false)

	return nil
}

func (dt *drawTable) TriggerEvent() error {
	if dt.untriggeredEvent == nil {
		return fmt.Errorf("no event to trigger")
	}

	if dt.Parent == nil {
		return fmt.Errorf("cannot trigger event: table doesn't support triggering events")
	}

	dt.FirstParent.enqueueEvent(dt.untriggeredEvent)
	dt.untriggeredEvent = nil

	return nil
}
