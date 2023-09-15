package main

import (
	"fmt"

	"github.com/nickwells/param.mod/v6/paction"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
)

const (
	paramNameRoughly     = "roughly"
	paramNameVeryRoughly = "very-roughly"
)

// addParams will add parameters to the passed ParamSet
func addParams(prog *Prog) func(ps *param.PSet) error {
	return func(ps *param.PSet) error {
		var unitFromName string
		var unitToName string

		var unitFamily *units.Family

		ps.Add("from", psetter.String[string]{Value: &unitFromName},
			"The units the value is in."+
				" It must be in the same family of units as the 'to' units.",
			param.Attrs(param.MustBeSet),
		)
		ps.Add("to", psetter.String[string]{Value: &unitToName},
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

		ps.Add("val", psetter.Float[float64]{Value: &prog.val},
			"the value to be converted.",
			param.AltNames("v"),
		)

		ps.Add("just-val", psetter.Bool{Value: &prog.justVal},
			"just show the result of the conversion and not"+
				" the from and to units as well."+
				" This flag will make the result easier to use in"+
				" scripts as only the result is shown.",
			param.AltNames("short", "s"),
		)

		ps.Add(paramNameRoughly, psetter.Nil{},
			"just show the result rounded to the nearest"+
				" multiple of 10 or 5 within 1% of the original value.",
			param.PostAction(paction.SetVal(&prog.roughly, true)),
			param.PostAction(paction.SetVal(&prog.roughPrecision, 1.0)),
			param.SeeAlso(paramNameVeryRoughly),
		)

		ps.Add(paramNameVeryRoughly, psetter.Nil{},
			"just show the result rounded to the nearest"+
				" multiple of 10 or 5 within 10% of the original value.",
			param.PostAction(paction.SetVal(&prog.roughly, true)),
			param.PostAction(paction.SetVal(&prog.roughPrecision, 10.0)),
			param.SeeAlso(paramNameRoughly),
		)

		ps.AddFinalCheck(func() error {
			var err error

			if unitFamily != nil {
				prog.unitFrom, err = unitFamily.GetUnit(unitFromName)
				if err != nil {
					return err
				}

				prog.unitTo, err = unitFamily.GetUnit(unitToName)
				if err != nil {
					return err
				}

				return nil
			}

			familiesHavingFromUnit := 0
			for _, fName := range units.GetFamilyNames() {
				f := units.GetFamilyOrPanic(fName)
				prog.unitFrom, err = f.GetUnit(unitFromName)
				if err == nil {
					familiesHavingFromUnit++
					prog.unitTo, err = f.GetUnit(unitToName)
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
