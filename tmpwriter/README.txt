package tmpwriter // import "github.com/HalCanary/facility/tmpwriter"


FUNCTIONS

func Touch(path string) error

TYPES

type TmpWriter struct {
	// Has unexported fields.
}

func Make(path string) (TmpWriter, error)

func (f *TmpWriter) Close() error

func (f *TmpWriter) Len() int

func (f *TmpWriter) Reset() error

func (f *TmpWriter) Write(p []byte) (int, error)

func (f *TmpWriter) WriteString(s string) (int, error)

