package main

// units

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/units.mod/units"
)

// Created: Sat Aug 29 16:52:07 2020

// conv holds the values describing the conversion to perform
type conv struct {
	unitFrom units.Unit
	unitTo   units.Unit

	val float64

	justVal bool
}

func main() {
	convVals := conv{val: 1.0}
	ps := paramset.NewOrDie(
		addParams(&convVals),
		addExamples,
		addReferences,
		param.SetProgramDescription(
			"This will convert a value between units of the same family"),
	)

	ps.Parse()

	v := units.ValWithUnit{Val: convVals.val, U: convVals.unitFrom}
	converted, err := v.Convert(convVals.unitTo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if convVals.justVal {
		fmt.Println(converted.Val)
	} else {
		fmt.Println(v, "=", converted)
	}
}

// addParams will add parameters to the passed ParamSet
func addParams(convVals *conv) func(ps *param.PSet) error {
	return func(ps *param.PSet) error {
		var unitFromName string
		var unitToName string

		var unitFamily string = "-"
		validFamilies := units.GetFamilyNames()
		avals := make(map[string]string)
		for _, f := range validFamilies {
			avals[f] = units.GetUnitDetailsOrPanic(f).Fam.Description
		}
		avals["-"] = "find the unit family from the unit names"

		ps.Add("from", psetter.String{Value: &unitFromName},
			"The units the value is in."+
				" It must be in the same family of units as the 'to' units.",
			param.Attrs(param.MustBeSet),
		)
		ps.Add("to", psetter.String{Value: &unitToName},
			"The units to convert the value into."+
				" It must be in the same family of units as the 'from' units.",
			param.Attrs(param.MustBeSet),
		)
		ps.Add("family",
			psetter.Enum{
				Value:       &unitFamily,
				AllowedVals: psetter.AllowedVals(avals),
			},
			"the family of units to use",
		)

		ps.Add("val", psetter.Float64{Value: &convVals.val},
			"the value to be converted.",
			param.AltName("v"),
		)

		ps.Add("just-val", psetter.Bool{Value: &convVals.justVal},
			"just show the result of the conversion and not"+
				" the from and to units as well."+
				" This flag will make the result easier to use in"+
				" scripts as only the result is shown.",
			param.AltName("short"),
			param.AltName("s"),
		)

		ps.AddFinalCheck(func() error {
			var err error

			if unitFamily != "-" {
				convVals.unitFrom, err = units.GetUnit(
					unitFamily, unitFromName)
				if err != nil {
					return err
				}

				convVals.unitTo, err = units.GetUnit(
					unitFamily, unitToName)
				if err != nil {
					return err
				}

				return nil
			}

			for _, family := range validFamilies {
				convVals.unitFrom, err = units.GetUnit(family, unitFromName)
				if err == nil {
					convVals.unitTo, err = units.GetUnit(family, unitToName)
					if err == nil {
						return nil
					}
				}
			}
			return fmt.Errorf(
				"There is no family of units having both %q and %q",
				unitFromName, unitToName)
		})

		return nil
	}
}

// addExamples adds some examples of how the program might be used
func addExamples(ps *param.PSet) error {
	ps.AddExample("unitconv -from pint -to litre",
		"This will show how many litres in a pint")
	ps.AddExample("unitconv -from chain -to mile",
		"This will show how many chains in a mile")
	ps.AddExample("unitconv -from chain -to mile -val 80",
		"This will show 80 chains in miles")
	ps.AddExample("unitconv -from chain -to m -val 80",
		"This will show 80 chains in metres")
	ps.AddExample("unitconv -from chain -to m -val 80 -just-val",
		"This will show 80 chains in metres. Only the value is "+
			"shown and no surrounding explanatory text")
	return nil
}

func addReferences(ps *param.PSet) error {
	ps.AddReference("unitlist",
		"This program can be used to list available units and"+
			" families of units. It can also give full details of"+
			" a particular unit"+
			"\n\n"+
			"To get this program:"+
			"\n\n"+
			"go install github.com/nickwells/unittools/unitlist@latest")

	return nil
}
