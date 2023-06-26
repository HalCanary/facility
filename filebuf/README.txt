package filebuf // import "github.com/HalCanary/facility/filebuf"

Copyright 2022 Hal Canary Use of this program is governed by the file LICENSE.

TYPES

type FileBuf struct {
	Path string

	// Has unexported fields.
}

func (b FileBuf) Changed() bool

func (b *FileBuf) Close() error

func (b *FileBuf) Len() int

func (b *FileBuf) Reset() error

func (b *FileBuf) Write(p []byte) (int, error)

func (b *FileBuf) WriteString(s string) (int, error)

