package humanize

// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.

import "testing"

func TestHumanize(t *testing.T) {
	for v, s := range map[int64]string{
		-9223372036854775808: "-9223372036854775808 B",
		-1:                   "-1 B",
		0:                    "0 B",
		1:                    "1 B",
		1023:                 "1023 B",
		1024:                 "1024 B",
		1025:                 "1025 B",
		2047:                 "2047 B",
		2048:                 "2048 B",
		2049:                 "2049 B",
		9999:                 "9999 B",
		10000:                "9 KB",
		10238976:             "9999 KB",
		10240000:             "9 MB",
		1073741823:           "1023 MB",
		1073741824:           "1024 MB",
		1073741825:           "1024 MB",
		9223372036854775807:  "8191 PB",
	} {
		u := Humanize(v)
		if u != s {
			t.Errorf("(%d) %q != %q\n", v, u, s)
		}
	}
}
