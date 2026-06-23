package main

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/param.mod/v7/paction"
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
)

const (
	paramNameFrom = "from"
	paramNameTo   = "to"

	paramNameNearest          = "nearest"
	paramNameNearestCount     = "nearest-count"
	paramNameNearestPrecision = "nearest-precision"
	paramNameNearestIgnoreTag = "nearest-ignore-tag"

	paramNameRoughly     = "roughly"
	paramNameVeryRoughly = "very-roughly"

	paramNameFamily    = "family"
	paramNameValue     = "value"
	paramNameJustValue = "just-value"
	paramNameWidth     = "width"
	paramNamePrecision = "precision"
)

const (
	roughPrecisionValue     = 1
	veryRoughPrecisionValue = 10
)

// addParams will add parameters to the passed ParamSet
func addParams(prog *prog) func(ps *param.PSet) error {
	return func(ps *param.PSet) error {
		const familyChoice = "The first family that has both" +
			" the 'from' and 'to'" +
			" units will be used." +
			" You can force the family to use" +
			" with the '" + paramNameFamily + "' parameter."

		var toOrBestCounter paction.Counter

		tOBCAF := toOrBestCounter.MakeActionFunc()

		ps.Add(paramNameFrom,
			psetter.String[string]{Value: &prog.unitFromName},
			"The units the value is in."+
				" It must be in the same family of units"+
				" as the '"+paramNameTo+"' units."+
				"\n\n"+
				familyChoice,
			param.Attrs(param.MustBeSet),
			param.SeeAlso(paramNameFamily, paramNameTo, paramNameNearest),
		)

		ps.Add(paramNameNearest,
			psetter.Bool{Value: &prog.nearestVal},
			"Convert the value into some unit in the same family of"+
				" units such that the quantity in that unit is some small,"+
				" preferably whole number value. A range of alternatives will"+
				" be shown.",
			param.PostAction(tOBCAF),
			param.SeeAlso(
				paramNameTo,
				paramNameNearestCount,
				paramNameNearestPrecision,
				paramNameNearestIgnoreTag,
			),
		)

		nearestCountParam := ps.Add(paramNameNearestCount,
			psetter.Int[int]{
				Value: &prog.nearestCount,
				Checks: []check.ValCk[int]{
					check.ValGE(1),
				},
			},
			"how many 'nearest' values should be shown.",
			param.SeeAlso(
				paramNameNearest,
				paramNameNearestPrecision,
				paramNameNearestIgnoreTag,
			),
		)

		nearestPrecisionParam := ps.Add(paramNameNearestPrecision,
			psetter.Float[float64]{
				Value: &prog.nearestPrecision,
				Checks: []check.ValCk[float64]{
					check.ValGT(0.0),
				},
			},
			"when generating the 'nearest' value,"+
				" how close to a whole number value to we allow"+
				" when comparing converted values.",
			param.SeeAlso(
				paramNameNearest,
				paramNameNearestCount,
				paramNameNearestIgnoreTag,
			),
		)

		nearestIgnoreTagsParam := ps.Add(paramNameNearestIgnoreTag,
			unitsetter.TagListAppender{
				Value: &prog.nearestIgnoreTags,
			},
			"when generating the 'nearest' value,"+
				" ignore any units with these tags.",
			param.SeeAlso(
				paramNameNearest,
				paramNameNearestCount,
				paramNameNearestPrecision,
			),
		)

		ps.Add(paramNameTo,
			psetter.StrList[string]{
				Value: &prog.unitToNames,
				Checks: []check.ValCk[[]string]{
					check.SliceLength[[]string](check.ValGT(0)),
					check.SliceAll[[]string](
						check.StringLength[string](check.ValGT(0))),
				},
			},
			"The units to convert the value into."+
				" They must all be in the same family of units as"+
				" the '"+paramNameFrom+"' unit."+
				"\n\n"+
				familyChoice,
			param.PostAction(tOBCAF),
			param.SeeAlso(paramNameFamily, paramNameFrom, paramNameNearest),
		)

		ps.Add(paramNameFamily,
			unitsetter.FamilySetter{
				Value: &prog.unitFamily,
			},
			"the family of units to use."+
				" The 'to' and 'from' units will be selected from this family.",
			param.AltNames("f", "fam"),
			param.SeeAlso(paramNameTo, paramNameFrom),
		)

		ps.Add(paramNameValue, psetter.Float[float64]{Value: &prog.val},
			"the value to be converted.",
			param.AltNames("v", "val"),
		)

		ps.Add(paramNameWidth, psetter.Int[int]{Value: &prog.displayWidth},
			"the space to allow for the display of the"+
				" converted value (the number part).",
			param.SeeAlso(paramNamePrecision),
		)

		ps.Add(paramNamePrecision, psetter.Int[int]{Value: &prog.displayPrec},
			"the number of digits of precision to allow"+
				" when displaying the"+
				" converted value (the number part).",
			param.AltNames("prec"),
			param.SeeAlso(paramNameWidth),
		)

		ps.Add(paramNameJustValue, psetter.Bool{Value: &prog.justVal},
			"just show the result of the conversion and not"+
				" the from and to units as well."+
				" This flag will make the result easier to use in"+
				" scripts as only the result is shown.",
			param.AltNames("just-val", "value-only", "val-only", "short", "s"),
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
			if toOrBestCounter.Count() != 1 {
				return fmt.Errorf(
					"you must give either %q or %q (not both)",
					paramNameTo, paramNameNearest)
			}

			if prog.nearestVal {
				return prog.findNearestVals()
			}

			if !prog.nearestVal &&
				(nearestCountParam.HasBeenSet() ||
					nearestPrecisionParam.HasBeenSet() ||
					nearestIgnoreTagsParam.HasBeenSet()) {
				return fmt.Errorf(
					"unless the %q parameter is given"+
						" the %q, %q or %q parameters have no effect",
					paramNameNearest,
					paramNameNearestCount,
					paramNameNearestPrecision,
					paramNameNearestIgnoreTag)
			}

			if prog.unitFamily != nil {
				return populateTargetUnitsFromFamily(prog)
			}

			for _, f := range units.GetFamilies() {
				prog.unitFamily = f
				if err := populateTargetUnitsFromFamily(prog); err == nil {
					return nil
				}
			}

			return fmt.Errorf("there is no unit-family having both %q and %s",
				prog.unitFromName,
				english.JoinQuoted(prog.unitToNames, ", ", " and "))
		})

		return nil
	}
}

// populateTargetUnitsFromFamily finds the units in the supplied family
func populateTargetUnitsFromFamily(prog *prog) error {
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
