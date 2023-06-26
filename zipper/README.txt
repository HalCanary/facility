package zipper // import "github.com/HalCanary/facility/zipper"

Copyright 2022 Hal Canary Use of this program is governed by the file LICENSE.

TYPES

type Zipper struct {
	ZipWriter *zip.Writer
	Error     error
}

func Make(dst io.Writer) Zipper

func (zw *Zipper) Close()

func (zw *Zipper) CreateDeflate(name string, mod time.Time) io.Writer

func (zw *Zipper) CreateStore(name string, mod time.Time) io.Writer

