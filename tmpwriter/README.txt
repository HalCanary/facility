package tmpwriter // import "github.com/HalCanary/facility/tmpwriter"


FUNCTIONS

func Touch(path string) error
    Touch a file, like the Unix `touch` utility.


TYPES

type TmpWriter struct {
	// Has unexported fields.
}
    Implements `io.Writer`, `io.StringWriter`, `io.WriteCloser`. Writes to a
    temp file so that calling `Reset` instead of `Close` leaves the original
    file in place.

func Make(path string) (TmpWriter, error)
    Create a new TmpWriter.

func (f *TmpWriter) Close() error
    Close the file. The file change is atomic, via `os.Rename`.

func (f *TmpWriter) Len() int
    Len returns the current file length.

func (f *TmpWriter) Reset() error
    Discard any changes to the file.

func (f *TmpWriter) Write(p []byte) (int, error)
    Write implements the standard Write interface.

func (f *TmpWriter) WriteString(s string) (int, error)
    WriteString implements the standard WriteString interface.

