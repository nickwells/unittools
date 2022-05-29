package main

// units

import (
	"fmt"
	"os"

	"github.com/nickwells/mathutil.mod/v2/mathutil"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
	"github.com/nickwells/unittools/internal/utparams"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// Created: Sat Aug 29 16:52:07 2020

// conv holds the values describing the conversion to perform
type conv struct {
	unitFrom units.Unit
	unitTo   units.Unit

	val float64

	justVal bool
	roughly bool
}

func main() {
	convVals := conv{val: 1.0}
	ps := paramset.NewOrDie(
		addParams(&convVals),
		versionparams.AddParams,
		addExamples,
		utparams.AddRefUnitlist,
		utparams.AddRefUnittags,
		param.SetProgramDescription(utparams.ProgDescUnitconv),
	)

	ps.Parse()

	v := units.ValUnit{V: convVals.val, U: convVals.unitFrom}
	converted, err := v.Convert(convVals.unitTo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if convVals.roughly {
		converted.V = mathutil.Roughly(converted.V, 1.0)
	}

	if convVals.justVal {
		fmt.Println(converted.V)
	} else {
		fmt.Println(v, "=", converted)
	}
}

// addParams will add parameters to the passed ParamSet
func addParams(convVals *conv) func(ps *param.PSet) error {
	return func(ps *param.PSet) error {
		var unitFromName string
		var unitToName string

		var unitFamily *units.Family

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
			unitsetter.FamilySetter{
				Value: &unitFamily,
			},
			"the family of units to use",
		)

		ps.Add("val", psetter.Float64{Value: &convVals.val},
			"the value to be converted.",
			param.AltNames("v"),
		)

		ps.Add("just-val", psetter.Bool{Value: &convVals.justVal},
			"just show the result of the conversion and not"+
				" the from and to units as well."+
				" This flag will make the result easier to use in"+
				" scripts as only the result is shown.",
			param.AltNames("short", "s"),
		)

		ps.Add("roughly", psetter.Bool{Value: &convVals.roughly},
			"just show the result rounded to the nearest"+
				" multiple of 10 or 5 within 1% of the original value.",
		)

		ps.AddFinalCheck(func() error {
			var err error

			if unitFamily != nil {
				convVals.unitFrom, err = unitFamily.GetUnit(unitFromName)
				if err != nil {
					return err
				}

				convVals.unitTo, err = unitFamily.GetUnit(unitToName)
				if err != nil {
					return err
				}

				return nil
			}

			familiesHavingFromUnit := 0
			for _, fName := range units.GetFamilyNames() {
				f := units.GetFamilyOrPanic(fName)
				convVals.unitFrom, err = f.GetUnit(unitFromName)
				if err == nil {
					familiesHavingFromUnit++
					convVals.unitTo, err = f.GetUnit(unitToName)
					if err == nil {
						return nil
					}
				}
			}
			if familiesHavingFromUnit == 0 {
				return fmt.Errorf("%q is not a valid unit in any unit-family",
					unitFromName)
			}
			return fmt.Errorf("There is no unit-family having both %q and %q",
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
