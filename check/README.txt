package check // import "github.com/HalCanary/facility/check"

Copyright 2022 Hal Canary Use of this program is governed by the file LICENSE.

FUNCTIONS

func Assert(condition bool)
    Runtime check that condition is true. If not, log failure and exit.

func Check(err error)
    Runtime check that err is nil. If not, log err and exit.

