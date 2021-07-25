package main

// unitlist

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/nickwells/col.mod/v3/col"
	"github.com/nickwells/col.mod/v3/col/colfmt"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
)

// Created: Fri Dec 25 18:42:35 2020

var (
	family *units.Family
	uName  string

	orderBySize bool
)

func main() {
	ps := paramset.NewOrDie(addParams,
		addExamples,
		addReferences,
		param.SetProgramDescription(
			"This describes available units and families of units."),
	)

	ps.Parse()

	if family == nil {
		listFamilies()
		return
	}
	if uName == "" {
		listUnits(family)
		return
	}
	showUnit(family, uName)
}

// getUnitIDs gets a sorted list of unit getUnitIDs
func getUnitIDs(f *units.Family) []string {
	unitIDs := f.GetUnitNames()

	if orderBySize {
		sort.Slice(unitIDs, func(i, j int) bool {
			iu, err := f.GetUnit(unitIDs[i])
			if err != nil {
				return false
			}
			ju, err := f.GetUnit(unitIDs[j])
			if err != nil {
				return false
			}
			if iu.ConvFactor() == ju.ConvFactor() {
				return iu.Name() < ju.Name()
			}
			return iu.ConvFactor() < ju.ConvFactor()
		})
	} else {
		sort.Strings(unitIDs)
	}
	return unitIDs
}

// getUnitNotes builds the notes column value for the given unit
func getUnitNotes(u units.Unit) string {
	notes := u.Notes()

	aliases := u.Aliases()
	if len(aliases) > 0 {
		aliasNames := []string{}
		for k := range aliases {
			aliasNames = append(aliasNames, k)
		}
		sort.Strings(aliasNames)
		notes += "\n\nAliases:\n"
		sep := ""
		for _, aName := range aliasNames {
			notes += sep + "    " + aName
			sep = "\n"
		}
	}

	if u.ConvPreAdd() != 0 || u.ConvPostAdd() != 0 {
		notes += "\n\n" +
			"The conversion is not a simple multiplication," +
			" show the full unit details for a full explanation."
	}

	return notes
}

// listUnits reports on the available units in the given family
func listUnits(f *units.Family) {
	unitIDs := getUnitIDs(f)

	rpt := col.StdRpt(
		col.New(&colfmt.String{}, "Base", "Unit"),
		col.New(&colfmt.WrappedString{W: 20}, "Unit Name"),
		col.New(&colfmt.Float{
			W:                        20,
			Prec:                     9,
			TrimTrailingZeroes:       true,
			ReformatOutOfBoundValues: true,
		}, "Conversion", "Factor"),
		col.New(&colfmt.WrappedString{W: 40}, "Notes"),
	)

	badUnits := []string{}
	for _, name := range unitIDs {
		intro := ""
		if name == f.BaseUnitName() {
			intro = ">>>"
		}
		u, err := f.GetUnit(name)
		if err != nil {
			badUnits = append(badUnits, name)
			continue
		}
		notes := getUnitNotes(u)
		err = rpt.PrintRow(intro, name, u.ConvFactor(), notes)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Error found while printing the %q units: %v\n",
				f.Name(), err)
			os.Exit(1)
		}
	}
	if len(badUnits) != 0 {
		fmt.Println("These units could not be found in the unit family:")
		fmt.Println(strings.Join(badUnits, "\n"))
	}
}

// listFamilies reports on the available families of units
func listFamilies() {
	validFamilies := units.GetFamilyNames()
	maxW := 0
	for _, f := range validFamilies {
		if len(f) > maxW {
			maxW = len(f)
		}
	}

	maxAliasW := 0
	aliases := units.GetFamilyAliases()
	for _, a := range aliases {
		if len(a) > maxAliasW {
			maxAliasW = len(a)
		}
	}
	rpt := col.StdRpt(
		col.New(&colfmt.String{W: maxW}, "Unit", "Family"),
		col.New(&colfmt.WrappedString{W: maxAliasW}, "Aliases"),
		col.New(&colfmt.String{}, "Description"),
	)
	for _, fName := range validFamilies {
		f := units.GetFamilyOrPanic(fName)
		err := rpt.PrintRow(
			fName,
			strings.Join(f.FamilyAliases(), "\n"),
			f.Description())
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Error found while printing the list of unit families:", err)
			os.Exit(1)
		}
	}
}

// addParams will add parameters to the passed ParamSet
func addParams(ps *param.PSet) error {
	familyParam := ps.Add("family",
		unitsetter.FamilySetter{Value: &family},
		"the family of units to use."+
			" If this is given without a unit then all the units for"+
			" the family will be listed."+
			"\n\n"+
			"If this is not given then a list of available families"+
			" will be shown.",
		param.AltName("f"),
	)

	unitParam := ps.Add("unit", psetter.String{Value: &uName},
		"the name of the unit to show. If this is given then"+
			" a family name must also be given."+
			" Full details of the unit will be displayed.",
		param.AltName("u"),
	)

	orderParam := ps.Add("by-size", psetter.Bool{Value: &orderBySize},
		"sort the units in size order not in alpabetical order."+
			" This should only be given when listing all the units"+
			" for a single family.",
	)

	ps.AddFinalCheck(func() error {
		if unitParam.HasBeenSet() {
			if !familyParam.HasBeenSet() {
				return errors.New("if a unit name is given" +
					" a family name must also be given")
			}
			if orderParam.HasBeenSet() {
				return errors.New("specifying the order of units" +
					" has no effect when showing a single unit")
			}
			return nil
		}

		if orderParam.HasBeenSet() && !familyParam.HasBeenSet() {
			return errors.New("specifying the order of units" +
				" has no effect when listing families of units")
		}
		return nil
	})

	return nil
}

// addExamples adds some examples of how the program might be used
func addExamples(ps *param.PSet) error {
	ps.AddExample("unitlist",
		"This will show the available families of units")
	ps.AddExample("unitlist -f temperature",
		"This will show the available units of temperature")
	ps.AddExample("unitlist -f temperature -u K",
		"This will show details of the 'K' unit of temperature")
	return nil
}

func addReferences(ps *param.PSet) error {
	ps.AddReference("unitconv",
		"This program can be used to convert between units."+
			"\n\n"+
			"To get this program:"+
			"\n\n"+
			"go install github.com/nickwells/unittools/unitconv@latest")

	return nil
}
