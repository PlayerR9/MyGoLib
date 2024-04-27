package Document

type Document struct {
}

func Regularize() {
	// 0x00 (null, NUL, \0, ^@), originally intended to be an ignored character, but now used by many
	// programming languages including C to mark the end of a string.

	// 0x07 (bell, BEL, \a, ^G), which may cause the device to emit a warning such as a bell or beep
	// sound or the screen flashing.

	// 0x08 (backspace, BS, \b, ^H), may overprint the previous character.

	// 0x09 (horizontal tab, HT, \t, ^I), moves the printing position right to the next tab stop.

	// 0x0A (line feed, LF, \n, ^J), moves the print head down one line, or to the left edge and down.
	// Used as the end of line marker in most UNIX systems and variants.

	// 0x0B (vertical tab, VT, \v, ^K), vertical tabulation.

	// 0x0C (form feed, FF, \f, ^L), to cause a printer to eject paper to the top of the next page,
	// or a video terminal to clear the screen.

	// 0x0D (carriage return, CR, \r, ^M), moves the printing position to the start of the line,
	// allowing overprinting. Used as the end of line marker in Classic Mac OS, OS-9, FLEX (and variants).
	// A CR+LF pair is used by CP/M-80 and its derivatives including DOS and Windows, and by
	// Application Layer protocols such as FTP, SMTP, and HTTP.

	// 0x1A (Control-Z, SUB, ^Z). Acts as an end-of-file for the Windows text-mode file i/o.

	// 0x1B (escape, ESC, \e (GCC only), ^[). Introduces an escape sequence.
}
