package spinner

import (
	"bytes"
	"fmt"
	"os"
)

// If os.Stderr is attached to a terminal, calling `Printf` will write a short
// message to stderr, overwriting the last such message.
type TerminalSpinner struct {
	on    bool
	init  bool
	width int
}

// Return a new Termial Spinner.
func NewTerminalSpinner() TerminalSpinner {
	return TerminalSpinner{}
}

// Printf will write a short message to stderr (if stderr is iattached to a
// terminal), overwriting the last such message, if that message was a
// contained no newline bytes.
func (t *TerminalSpinner) Printf(format string, a ...any) {
	if !t.init {
		stderrStat, _ := os.Stderr.Stat()
		t.init, t.on = true, stderrStat.Mode()&os.ModeCharDevice != 0
	}
	if t.on {
		os.Stderr.Write(bytes.Repeat([]byte{'\010'}, t.width))
		s := fmt.Sprintf(format, a...)
		os.Stderr.WriteString(s)
		t.width = len(s)
	}
}
