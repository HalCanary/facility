package spinner // import "github.com/HalCanary/facility/spinner"


TYPES

type TerminalSpinner struct {
	// Has unexported fields.
}
    If os.Stderr is attached to a terminal, calling `Printf` will write a short
    message to stderr, overwriting the last such message.

func NewTerminalSpinner() TerminalSpinner
    Return a new Termial Spinner.

func (t *TerminalSpinner) Printf(format string, a ...any)
    Printf will write a short message to stderr (if stderr is iattached to
    a terminal), overwriting the last such message, if that message was a
    contained no newline bytes.

