package spinner // import "github.com/HalCanary/facility/spinner"


TYPES

type TerminalSpinner struct {
	// Has unexported fields.
}

func NewTerminalSpinner() TerminalSpinner

func (t *TerminalSpinner) Printf(format string, a ...any)

