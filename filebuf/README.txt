package filebuf // import "github.com/HalCanary/facility/filebuf"

Copyright 2022 Hal Canary Use of this program is governed by the file LICENSE.

TYPES

type FileBuf struct {
	Path string

	// Has unexported fields.
}
    Implements `io.Writer`, `io.StringWriter`, `io.WriteCloser`. When `Close()`
    is called, writes buffer to `Path`, if it differs from the contents of the
    file at `Path` or there is no file at `Path`. `Path` must be populated by
    the caller.

func (b FileBuf) Changed() bool
    After `Close` is called, this will return `true` if the file changed.

func (b *FileBuf) Close() error
    Close writes buffer to `Path`, if it differs from the contents of the file
    at `Path` or there is no file at `Path`.

func (b *FileBuf) Len() int
    Len returns the number of bytes in the buffer.

func (b *FileBuf) Reset() error
    Reset returns the internal buffer to empty.

func (b *FileBuf) Write(p []byte) (int, error)
    Write implements the standard Write interface

func (b *FileBuf) WriteString(s string) (int, error)
    WriteString writes the contents of s.

