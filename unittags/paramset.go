package main

import (
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/unittools/internal/utparams"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// makeParamSet generates the param set ready for parsing
func makeParamSet(prog *Prog) *param.PSet {
	pName := utparams.ProgNameUnittags

	return paramset.NewOrPanic(
		versionparams.AddParams,

		addParams(prog),
		utparams.AddRefs(pName),

		utparams.SetProgramDescription(pName),
	)
}
