package utparams

import (
	"fmt"
	"maps"
	"slices"

	"github.com/nickwells/param.mod/v6/param"
)

const (
	ProgNameUnitlist = "unitlist"
	ProgNameUnitconv = "unitconv"
	ProgNameUnittags = "unittags"

	repo = "github.com/nickwells/unittools/"
)

// descriptions holds the recognised programs and their associated
// descriptions. It is used to generate references, program descriptions and
// to validate the program names
var descriptions = map[string]string{
	ProgNameUnitlist: "This can be used to list" +
		" the available families of units and" +
		" the units in those families." +
		" It can also give full details of a particular unit",

	ProgNameUnitconv: "This will convert a value" +
		" between units of the same family.",

	ProgNameUnittags: "This lists the available unit tags" +
		" and gives an explanation of their meaning",
}

// AddRefs returns a PSetOptFunc which adds the references for all the
// unittools programs with a different name to the supplied value. It will
// return an error if the supplied name does not exist in the map of
// descriptions.
func AddRefs(selfName string) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		prognames := slices.Sorted(maps.Keys(descriptions))
		foundSelf := false

		for _, pName := range prognames {
			if pName == selfName {
				foundSelf = true
				continue
			}

			ps.AddReference(pName,
				descriptions[pName]+
					"\n\n"+
					"To get this program:"+
					"\n\n"+
					"go install "+repo+pName+"@latest")
		}

		if !foundSelf {
			return fmt.Errorf("the program name %q is not recognised",
				selfName)
		}

		return nil
	}
}

// SetProgramDescription returns a PSetOptFunc which sets the program
// description for the named program. It will return an error if the program
// name is not found in the map of program descriptions.
func SetProgramDescription(selfName string) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		desc, ok := descriptions[selfName]
		if !ok {
			return fmt.Errorf("the program name %q is not recognised",
				selfName)
		}

		return param.SetProgramDescription(desc)(ps)
	}
}
