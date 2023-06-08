// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.
package check

import (
	"log"
	"os"
)

// Runtime check that err is nil.  If not, log err and exit.
func Check(err error) {
	if err != nil {
		log.Output(2, err.Error())
		os.Exit(1)
	}
}

// Runtime check that condition is true.  If not, log failure and exit.
func Assert(condition bool) {
	if !condition {
		log.Output(2, "Assertial Failure")
		os.Exit(1)
	}
}
