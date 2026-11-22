package main

import (
	"testing"

	"github.com/nickwells/testhelper.mod/v3/testhelper"
)

func TestMakeParamSet(t *testing.T) {
	prog := newProg()
	pi := testhelper.PanicSafe(func() {
		_ = makeParamSet(prog)
	})
	testhelper.CheckNoPanic(t, "makeParamSet", pi)
}
