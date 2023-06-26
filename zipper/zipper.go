// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.
package zipper

import (
	"archive/zip"
	"compress/flate"
	"io"
	"time"
)

// Wrapper for archive/zip.Writer.  Keeps track of its error state.
type Zipper struct {
	ZipWriter *zip.Writer
	Error     error
}

// Wraps `dst` with a `zip.Writer`.  `Deflate` always uses `BestCompression`.
func Make(dst io.Writer) Zipper {
	zw := Zipper{zip.NewWriter(dst), nil}
	zw.ZipWriter.RegisterCompressor(zip.Deflate, makeBestFlateWriter)
	return zw
}

func makeBestFlateWriter(w io.Writer) (io.WriteCloser, error) {
	return flate.NewWriter(w, flate.BestCompression)
}

// Close the underlying `ZipWriter`
func (zw *Zipper) Close() {
	err := zw.ZipWriter.Close()
	if zw.Error == nil {
		zw.Error = err
	}
}

func (zw *Zipper) create(name string, method uint16, mod time.Time) io.Writer {
	if !mod.IsZero() {
		mod = mod.UTC()
	}
	if zw.Error == nil {
		var w io.Writer
		if w, zw.Error = zw.ZipWriter.CreateHeader(&zip.FileHeader{
			Name:     name,
			Modified: mod,
			Method:   method,
		}); zw.Error == nil {
			return w
		}
	}
	return nil
}

// Add a new file to the zip archive (with deflate).  Returns nil on error.
func (zw *Zipper) CreateDeflate(name string, mod time.Time) io.Writer {
	return zw.create(name, zip.Deflate, mod)
}

// Add a new file to the zip archive (with store).  Returns nil on error.
func (zw *Zipper) CreateStore(name string, mod time.Time) io.Writer {
	return zw.create(name, zip.Store, mod)
}
