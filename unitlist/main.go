package main

// unitlist

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/nickwells/col.mod/v2/col"
	"github.com/nickwells/col.mod/v2/col/colfmt"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/units.mod/units"
)

// Created: Fri Dec 25 18:42:35 2020

var fName string
var uName string
var orderBySize bool

func main() {
	ps := paramset.NewOrDie(addParams,
		addExamples,
		addReferences,
		param.SetProgramDescription(
			"This describes available units and families of units."),
	)

	ps.Parse()

	if fName == "" {
		listFamilies()
		return
	}
	if uName == "" {
		listUnits(fName)
		return
	}
	showUnit(fName, uName)
}

// listUnits reports on the available units in the given family
func listUnits(fName string) {
	ud := units.GetUnitDetailsOrPanic(fName)
	names := make([]string, 0, len(ud.AltU))
	maxNameLen := 0
	for name := range ud.AltU {
		names = append(names, name)
		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}
	}
	if orderBySize {
		sort.Slice(names, func(i, j int) bool {
			iu, err := units.GetUnit(fName, names[i])
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"%q is not a unit of %s", uName, fName)
				os.Exit(1)
			}
			ju, err := units.GetUnit(fName, names[j])
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"%q is not a unit of %s", uName, fName)
				os.Exit(1)
			}
			return iu.ConvFactor < ju.ConvFactor
		})
	} else {
		sort.Strings(names)
	}
	h := col.NewHeaderOrPanic()
	rpt := col.NewReportOrPanic(h, os.Stdout,
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
	for _, name := range names {
		intro := ""
		if name == ud.Fam.BaseUnitName {
			intro = ">>>"
		}
		err := rpt.PrintRow(intro,
			name,
			ud.AltU[name].ConvFactor,
			ud.AltU[name].Notes)
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Error found while printing the report:", err)
			os.Exit(1)
		}
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
	h := col.NewHeaderOrPanic()
	rpt := col.NewReportOrPanic(h, os.Stdout,
		col.New(&colfmt.String{W: maxW}, "Unit", "Family"),
		col.New(&colfmt.String{}, "Description"),
	)
	for _, f := range validFamilies {
		err := rpt.PrintRow(f, units.GetUnitDetailsOrPanic(f).Fam.Description)
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Error found while printing the report:", err)
			os.Exit(1)
		}
	}
}

// addParams will add parameters to the passed ParamSet
func addParams(ps *param.PSet) error {
	validFamilies := units.GetFamilyNames()
	avals := make(map[string]string)
	for _, f := range validFamilies {
		avals[f] = units.GetUnitDetailsOrPanic(f).Fam.Description
	}

	familyParam := ps.Add("family",
		psetter.Enum{
			Value:                    &fName,
			AllowedVals:              psetter.AllowedVals(avals),
			AllowInvalidInitialValue: true,
		},
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
			"go get -u github.com/nickwells/unittools/unitconv")

	return nil
}
