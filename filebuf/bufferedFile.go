// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.
package filebuf

import (
	"bytes"
	"os"
	"path/filepath"
)

type FileBuf struct {
	Path    string
	buf     bytes.Buffer
	changed bool
}

func (b *FileBuf) Write(p []byte) (int, error) { return b.buf.Write(p) }

func (b *FileBuf) WriteString(s string) (int, error) { return b.buf.WriteString(s) }

func (b *FileBuf) Reset() error {
	b.buf = bytes.Buffer{}
	return nil
}

func (b *FileBuf) Len() int { return b.buf.Len() }

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

func (b FileBuf) Changed() bool { return b.changed }
