package utparams

import "github.com/nickwells/param.mod/v5/param"

const (
	ProgDescUnitlist = "This can be used to list" +
		" the available families of units and" +
		" the units in those families." +
		" It can also give full details of a particular unit"

	ProgDescUnitconv = "This will convert a value" +
		" between units of the same family."

	ProgDescUnittags = "This lists the available unit tags" +
		" and gives an explanation of their meaning"
)

// AddRefUnitlist adds a reference to the unitlist program
func AddRefUnitlist(ps *param.PSet) error {
	pName := "unitlist"
	ps.AddReference(pName,
		ProgDescUnitlist+
			installNotes(pName))

	return nil
}

// AddRefUnitconv adds a reference to the unitconv program
func AddRefUnitconv(ps *param.PSet) error {
	pName := "unitconv"
	ps.AddReference(pName,
		ProgDescUnitconv+
			installNotes(pName))

	return nil
}

// AddRefUnittags adds a reference to the unittags program
func AddRefUnittags(ps *param.PSet) error {
	pName := "unittags"
	ps.AddReference(pName,
		ProgDescUnittags+
			installNotes(pName))

	return nil
}
