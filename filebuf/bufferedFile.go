// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.
package filebuf

import (
	"bytes"
	"os"
	"path/filepath"
)

// Implements `io.Writer`, `io.StringWriter`, `io.WriteCloser`.
// When `Close()` is called, writes buffer to `Path`, if it differs from the
// contents of the file at `Path` or there is no file at `Path`.
// `Path` must be populated by the caller.
type FileBuf struct {
	Path    string
	buf     bytes.Buffer
	changed bool
}

// Write implements the standard Write interface
func (b *FileBuf) Write(p []byte) (int, error) { return b.buf.Write(p) }

// WriteString writes the contents of s.
func (b *FileBuf) WriteString(s string) (int, error) { return b.buf.WriteString(s) }

// Reset returns the internal buffer to empty.
func (b *FileBuf) Reset() error {
	b.buf = bytes.Buffer{}
	return nil
}

// Len returns the number of bytes in the buffer.
func (b *FileBuf) Len() int { return b.buf.Len() }

// Close writes buffer to `Path`, if it differs from the
// contents of the file at `Path` or there is no file at `Path`.
func (b *FileBuf) Close() error {
	d := b.buf.Bytes()
	b.buf = bytes.Buffer{}
	c, e := os.ReadFile(b.Path)
	if e == nil && bytes.Equal(c, d) {
		return nil
	}
	if e = os.MkdirAll(filepath.Dir(b.Path), 0o777); e != nil {
		return e
	}
	f, e := os.Create(b.Path)
	if e != nil {
		return e
	}
	b.changed = true
	_, e = f.Write(d)
	_ = f.Close()
	return e
}

// After `Close` is called, this will return `true` if the file changed.
func (b FileBuf) Changed() bool { return b.changed }
