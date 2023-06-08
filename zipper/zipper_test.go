package zipper

// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.

import "testing"
import "bytes"

func TestZipper(t *testing.T) {
	var buffer bytes.Buffer
	z := Make(&buffer)
	z.Close()
	var expected = []byte{'P', 'K', 5, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !bytes.Equal(buffer.Bytes(), expected) {
		t.Errorf("%v!=%v", buffer.Bytes(), expected)
	}
}
