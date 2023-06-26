package tmpwriter

import (
	"bytes"
	"os"
	"path/filepath"
	"time"
)

// Touch a file, like the Unix `touch` utility.
func Touch(path string) error {
	if _, e := os.Stat(path); os.IsNotExist(e) {
		f, e := os.Create(path)
		if e != nil {
			return e
		}
		return f.Close()
	}
	t := time.Now()
	return os.Chtimes(path, t, t)
}

////////////////////////////////////////////////////////////////////////////////

// Implements `io.Writer`, `io.StringWriter`, `io.WriteCloser`.
// Writes to a temp file so that calling `Reset` instead of `Close` leaves the
// original file in place.
type TmpWriter struct {
	path string
	file *os.File
}

// Create a new TmpWriter.
func Make(path string) (TmpWriter, error) {
	f, e := os.CreateTemp(filepath.Dir(path), "")
	return TmpWriter{path: path, file: f}, e
}

// Write implements the standard Write interface.
func (f *TmpWriter) Write(p []byte) (int, error) { return f.file.Write(p) }

// WriteString implements the standard WriteString interface.
func (f *TmpWriter) WriteString(s string) (int, error) { return f.file.WriteString(s) }

// Len returns the current file length.
func (f *TmpWriter) Len() int {
	n, _ := f.file.Seek(0, 1)
	return int(n)
}

// Discard any changes to the file.
func (f *TmpWriter) Reset() error {
	t := f.file.Name()
	_ = f.file.Close()
	return os.Remove(t)
}

// Close the file.  The file change is atomic, via `os.Rename`.
func (f *TmpWriter) Close() error {
	t := f.file.Name()
	_ = f.file.Close()
	c, e := os.ReadFile(f.path)
	if e != nil {
		if ct, et := os.ReadFile(t); et == nil && bytes.Equal(ct, c) {
			_ = os.Remove(t)
			return nil
		}
	}
	e = os.Rename(t, f.path)
	os.Chmod(f.path, 0o644)
	return e
}
