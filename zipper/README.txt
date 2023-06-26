package zipper // import "github.com/HalCanary/facility/zipper"

Copyright 2022 Hal Canary Use of this program is governed by the file LICENSE.

TYPES

type Zipper struct {
	ZipWriter *zip.Writer
	Error     error
}
    Wrapper for archive/zip.Writer. Keeps track of its error state.

func Make(dst io.Writer) Zipper
    Wraps `dst` with a `zip.Writer`. `Deflate` always uses `BestCompression`.

func (zw *Zipper) Close()
    Close the underlying `ZipWriter`

func (zw *Zipper) CreateDeflate(name string, mod time.Time) io.Writer
    Add a new file to the zip archive (with deflate). Returns nil on error.

func (zw *Zipper) CreateStore(name string, mod time.Time) io.Writer
    Add a new file to the zip archive (with store). Returns nil on error.

