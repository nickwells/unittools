package main

import (
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramset"
	"github.com/nickwells/unittools/internal/utparams"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// makeParamSet generates the param set ready for parsing
func makeParamSet(prog *prog) *param.PSet {
	pName := utparams.ProgNameUnittags

	return paramset.New(
		versionparams.AddParams,

		addParams(prog),
		utparams.AddRefs(pName),

		utparams.SetProgramDescription(pName),
	)
}
