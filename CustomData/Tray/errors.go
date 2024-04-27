package Tray

type ErrEmptyTray struct{}

func (e *ErrEmptyTray) Error() string {
	return "tray is empty"
}

func NewErrEmptyTray() *ErrEmptyTray {
	return &ErrEmptyTray{}
}
