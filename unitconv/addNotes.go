package main

import (
	"github.com/nickwells/param.mod/v7/param"
)

const (
	noteBaseName = "unitconv - "

	noteNameNearest = noteBaseName + "nearest conversion"
)

// addNotes adds the notes for this program.
func addNotes(_ *prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.AddNote(noteNameNearest,
			"if you pass the program the '"+paramNameNearest+"'"+
				" parameter then it will try to find units that,"+
				" when the value is converted, will yield"+
				" small, whole numbers or simple fractions."+
				"\n\n"+
				"It can be useful if you want to understand"+
				" the derivation of a value. For instance,"+
				" given a distance of 10.05534 metres,"+
				" using this parameter (especially with the"+
				" '"+paramNameRoughly+"' parameter) can help"+
				" you identify this as two rods or half a chain.")
		return nil
	}
}
