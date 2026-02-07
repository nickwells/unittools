package main

import (
	"errors"

	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
)

const (
	paramNameTagged    = "tagged"
	paramNameNotTagged = "not-tagged"
)

// addParams will add parameters to the passed ParamSet
func addParams(prog *prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		familyParam := ps.Add("family",
			unitsetter.FamilySetter{Value: &prog.family},
			"the family of units to use."+
				" If this is given without a unit then all the units for"+
				" the family will be listed."+
				"\n\n"+
				"If this is not given then a list of available families"+
				" will be shown.",
			param.AltNames("f"),
		)

		unitParam := ps.Add("unit", psetter.String[string]{Value: &prog.uName},
			"the name of the unit to show. If this is given then"+
				" a family name must also be given."+
				" Full details of the unit will be displayed.",
			param.AltNames("u"),
		)

		orderParam := ps.Add("by-name", psetter.Bool{Value: &prog.orderByName},
			"sort the units in alpabetical order not in size order."+
				"\n\n"+
				"This should only be given when listing all the units"+
				" for a single family.",
		)

		ps.Add("tagged",
			unitsetter.TagListAppender{Value: &prog.mustHaveTags},
			"only show units which have the given tag."+
				" This should only be given when listing all the units"+
				" Repetitions of this parameter"+
				" will add to the list of tags that must be present."+
				"\n\n"+
				"This should only be given when listing all the units"+
				" for a single family.",
			param.AltNames("tag"),
			param.SeeAlso(paramNameNotTagged),
		)

		ps.Add(paramNameNotTagged,
			unitsetter.TagListAppender{Value: &prog.mustNotHaveTags},
			"only show units which do not have the given tag."+
				" for a single family. Repetitions of this parameter"+
				" will add to the list of tags that must be missing."+
				"\n\n"+
				"This should only be given when listing all the units"+
				" for a single family.",
			param.SeeAlso(paramNameTagged),
		)

		detailsParam := ps.Add("show-details",
			psetter.Bool{Value: &prog.showDetail},
			"show details when listing."+
				"\n\n"+
				"This should not be given when"+
				" showing details for a single unit.",
			param.AltNames("show-detail", "l"),
		)

		noHdrParam := ps.Add("no-header",
			psetter.Bool{Value: &prog.noHeader},
			"don't show the column headings when listing."+
				"\n\n"+
				"This should not be given when"+
				" showing details for a single unit.",
			param.AltNames("no-hdr"),
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

				if prog.hasTagConstraints() {
					return errors.New(
						"constraining the units to list by tag name" +
							" has no effect when showing a single unit")
				}

				if detailsParam.HasBeenSet() {
					return errors.New("asking to see more detail" +
						" has no effect when showing a single unit")
				}

				if noHdrParam.HasBeenSet() {
					return errors.New("asking to not show headers" +
						" has no effect when showing a single unit")
				}

				return nil
			}

			if !familyParam.HasBeenSet() {
				if orderParam.HasBeenSet() {
					return errors.New("specifying the order of units" +
						" only has an effect when listing units in a family")
				}

				if prog.hasTagConstraints() {
					return errors.New("constraining the units to list" +
						" only has an effect when listing units in a family")
				}
			}

			return prog.checkTagLists()
		})

		return nil
	}
}
