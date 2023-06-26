// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.
package humanize

import "fmt"

// Humanize converts a byte size to a human readable number, for example: 10240
// becomes "10 KB.  Everything is rounded up to nearest integer."
func Humanize(v int64) string {
	prfx := []string{"", "K", "M", "G", "T", "P", "E"}
	for i, s := range prfx {
		if v <= 9999 || i == len(prfx)-1 {
			return fmt.Sprintf("%d %sB", v, s)
		}
		v = ((v - 1) >> 10) + 1
	}
	return ""
}
