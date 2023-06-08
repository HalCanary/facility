// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.
package zipper

import (
	"archive/zip"
	"compress/flate"
	"io"
	"time"
)

type Zipper struct {
	ZipWriter *zip.Writer
	Error     error
}

func Make(dst io.Writer) Zipper {
	zw := Zipper{zip.NewWriter(dst), nil}
	zw.ZipWriter.RegisterCompressor(zip.Deflate, makeBestFlateWriter)
	return zw
}

func makeBestFlateWriter(w io.Writer) (io.WriteCloser, error) {
	return flate.NewWriter(w, flate.BestCompression)
}

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

func (zw *Zipper) CreateDeflate(name string, mod time.Time) io.Writer {
	return zw.create(name, zip.Deflate, mod)
}

func (zw *Zipper) CreateStore(name string, mod time.Time) io.Writer {
	return zw.create(name, zip.Store, mod)
}
