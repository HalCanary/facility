package spinner

import (
	"bytes"
	"fmt"
	"os"
)

type TerminalSpinner struct {
	on    bool
	width int
}

func NewTerminalSpinner() TerminalSpinner {
	stderrStat, _ := os.Stderr.Stat()
	return TerminalSpinner{on: stderrStat.Mode()&os.ModeCharDevice != 0}
}

func (t *TerminalSpinner) Printf(format string, a ...any) {
	if t.on {
		os.Stderr.Write(bytes.Repeat([]byte{'\010'}, t.width))
		s := fmt.Sprintf(format, a...)
		os.Stderr.WriteString(s)
		t.width = len(s)
	}
}
