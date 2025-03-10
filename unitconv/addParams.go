package main

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/param.mod/v6/paction"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
)

const (
	paramNameFrom        = "from"
	paramNameTo          = "to"
	paramNameRoughly     = "roughly"
	paramNameVeryRoughly = "very-roughly"
)

const (
	roughPrecisionValue     = 1
	veryRoughPrecisionValue = 10
)

// addParams will add parameters to the passed ParamSet
func addParams(prog *Prog) func(ps *param.PSet) error {
	return func(ps *param.PSet) error {
		ps.Add(paramNameFrom, psetter.String[string]{Value: &prog.unitFromName},
			"The units the value is in."+
				" It must be in the same family of units"+
				" as the '"+paramNameTo+"' units.",
			param.Attrs(param.MustBeSet),
		)
		ps.Add(paramNameTo, psetter.StrList[string]{
			Value: &prog.unitToNames,
			Checks: []check.ValCk[[]string]{
				check.SliceLength[[]string, string](check.ValGT[int](0)),
				check.SliceAll[[]string, string](
					check.StringLength[string](check.ValGT[int](0))),
			},
		},
			"The units to convert the value into."+
				" They must all be in the same family of units as"+
				" the '"+paramNameFrom+"' unit.",
			param.Attrs(param.MustBeSet),
		)
		ps.Add("family",
			unitsetter.FamilySetter{
				Value: &prog.unitFamily,
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
			fmt.Sprintf("just show the result rounded to the nearest"+
				" multiple of 10 or 5 within %d%% of the original value.",
				roughPrecisionValue),
			param.PostAction(paction.SetVal(&prog.roughly, true)),
			param.PostAction(paction.SetVal(&prog.roughPrecision,
				roughPrecisionValue)),
			param.SeeAlso(paramNameVeryRoughly),
		)

		ps.Add(paramNameVeryRoughly, psetter.Nil{},
			fmt.Sprintf("just show the result rounded to the nearest"+
				" multiple of 10 or 5 within %d%% of the original value.",
				veryRoughPrecisionValue),
			param.PostAction(paction.SetVal(&prog.roughly, true)),
			param.PostAction(paction.SetVal(&prog.roughPrecision,
				veryRoughPrecisionValue)),
			param.SeeAlso(paramNameRoughly),
		)

		ps.AddFinalCheck(func() error {
			if prog.unitFamily != nil {
				return populateTargetUnitsFromFamily(prog)
			}

			for _, fName := range units.GetFamilyNames() {
				prog.unitFamily = units.GetFamilyOrPanic(fName)
				if err := populateTargetUnitsFromFamily(prog); err == nil {
					return nil
				}
			}

			return fmt.Errorf("There is no unit-family having both %q and %s",
				prog.unitFromName,
				english.JoinQuoted(prog.unitToNames, ", ", " and ", `"`, `"`))
		})

		return nil
	}
}

// populateTargetUnitsFromFamily finds the units in the supplied family
func populateTargetUnitsFromFamily(prog *Prog) error {
	var err error

	prog.unitFrom, err = prog.unitFamily.GetUnit(prog.unitFromName)
	if err != nil {
		return err
	}

	prog.unitTo = []units.Unit{}
	for _, unitName := range prog.unitToNames {
		u, err := prog.unitFamily.GetUnit(unitName)
		if err != nil {
			return err
		}

		prog.unitTo = append(prog.unitTo, u)
	}

	return nil
}
