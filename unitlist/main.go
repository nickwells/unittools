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
	"github.com/nickwells/unittools/internal/utparams"
)

// Created: Fri Dec 25 18:42:35 2020

type unitlist struct {
	family *units.Family
	uName  string

	mustHaveTags    []units.Tag
	mustNotHaveTags []units.Tag

	orderByName bool
	showDetail  bool
	noHeader    bool
}

func main() {
	ul := unitlist{}
	ps := paramset.NewOrDie(addParams(&ul),
		addExamples,
		utparams.AddRefUnitconv,
		utparams.AddRefUnittags,
		param.SetProgramDescription(utparams.ProgDescUnitlist),
	)

	ps.Parse()

	if ul.family == nil {
		listFamilies(ul)
		return
	}
	if ul.uName == "" {
		listUnits(ul)
		return
	}
	showUnit(ul)
}

// getUnitIDs gets a sorted list of unit getUnitIDs
func getUnitIDs(ul unitlist) []string {
	unitIDs := ul.family.GetUnitNames()

	if ul.orderByName {
		sort.Strings(unitIDs)
	} else {
		sort.Slice(unitIDs, func(i, j int) bool {
			iu, err := ul.family.GetUnit(unitIDs[i])
			if err != nil {
				return false
			}
			ju, err := ul.family.GetUnit(unitIDs[j])
			if err != nil {
				return false
			}
			if iu.ConvFactor() == ju.ConvFactor() {
				return iu.Name() < ju.Name()
			}
			return iu.ConvFactor() < ju.ConvFactor()
		})
	}
	return unitIDs
}

// getUnitTags returns the tags for this getUnitTags
func getUnitTags(u units.Unit) string {
	tags := u.Tags()
	rval := ""
	sep := ""
	for _, tag := range tags {
		rval += sep + string(tag)
		sep = "\n"
	}
	return rval
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

// makeUnitListRpt generates the appropriate report taking into account the
// noHeader and showDetail flags
func makeUnitListRpt(ul unitlist) *col.Report {
	hdr := col.NewHeaderOrPanic()
	if ul.noHeader {
		hdr = col.NewHeaderOrPanic(col.HdrOptDontPrint)
	}

	if ul.showDetail {
		return col.NewReport(hdr, os.Stdout,
			col.New(&colfmt.String{}, "Base", "Unit"),
			col.New(&colfmt.WrappedString{W: 20}, "Unit Name"),
			col.New(&colfmt.WrappedString{W: 20}, "Tags"),
			col.New(&colfmt.Float{
				W:                        20,
				Prec:                     9,
				TrimTrailingZeroes:       true,
				ReformatOutOfBoundValues: true,
			}, "Conversion", "Factor"),
			col.New(&colfmt.WrappedString{W: 40}, "Notes"),
		)
	}

	return col.NewReport(hdr, os.Stdout, col.New(&colfmt.String{}, "Unit Name"))
}

// printUnitRow prints the row in the unit list report. It returns false if
// the unit cannot be found, true otherwise.
func printUnitRow(ul unitlist, rpt *col.Report, uName string) bool {
	u, err := ul.family.GetUnit(uName)
	if err != nil {
		return false
	}

	for _, tag := range ul.mustHaveTags {
		if !u.HasTag(tag) {
			return true
		}
	}
	for _, tag := range ul.mustNotHaveTags {
		if u.HasTag(tag) {
			return true
		}
	}

	if !ul.showDetail {
		err = rpt.PrintRow(uName)
	} else {
		intro := ""
		if uName == ul.family.BaseUnitName() {
			intro = ">>>"
		}

		err = rpt.PrintRow(
			intro, uName, getUnitTags(u), u.ConvFactor(), getUnitNotes(u))
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error found while printing the %q units: %v\n",
			ul.family.Name(), err)
		os.Exit(1)
	}

	return true
}

// listUnits reports on the available units in the given family
func listUnits(ul unitlist) {
	unitIDs := getUnitIDs(ul)

	rpt := makeUnitListRpt(ul)

	badUnits := []string{}
	for _, uName := range unitIDs {
		if !printUnitRow(ul, rpt, uName) {
			badUnits = append(badUnits, uName)
		}
	}
	if len(badUnits) != 0 {
		fmt.Println("These units could not be found in the unit family:")
		fmt.Println(strings.Join(badUnits, "\n"))
	}
}

// makeFamilyListRpt generates the appropriate report taking into account the
// noHeader and showDetail flags
func makeFamilyListRpt(ul unitlist) *col.Report {
	hdr := col.NewHeaderOrPanic()
	if ul.noHeader {
		hdr = col.NewHeaderOrPanic(col.HdrOptDontPrint)
	}

	if ul.showDetail {
		maxW := 0
		validFamilies := units.GetFamilyNames()
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
		return col.NewReport(hdr, os.Stdout,
			col.New(&colfmt.String{W: maxW}, "Unit", "Family"),
			col.New(&colfmt.WrappedString{W: maxAliasW}, "Aliases"),
			col.New(&colfmt.String{}, "Description"),
		)
	}

	return col.NewReport(hdr, os.Stdout,
		col.New(&colfmt.String{}, "Unit Family"))
}

// printFamilyRow prints the row in the family list report.
func printFamilyRow(ul unitlist, rpt *col.Report, fName string) {
	var err error

	if ul.showDetail {
		f := units.GetFamilyOrPanic(fName)
		err = rpt.PrintRow(
			fName,
			strings.Join(f.FamilyAliases(), "\n"),
			f.Description())
	} else {
		err = rpt.PrintRow(fName)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr,
			"Error found while printing the list of unit families:", err)
		os.Exit(1)
	}
}

// listFamilies reports on the available families of units
func listFamilies(ul unitlist) {
	validFamilies := units.GetFamilyNames()
	sort.Strings(validFamilies)

	rpt := makeFamilyListRpt(ul)

	for _, fName := range validFamilies {
		printFamilyRow(ul, rpt, fName)
	}
}

// addParams will add parameters to the passed ParamSet
func addParams(ul *unitlist) func(ps *param.PSet) error {
	return func(ps *param.PSet) error {
		familyParam := ps.Add("family",
			unitsetter.FamilySetter{Value: &ul.family},
			"the family of units to use."+
				" If this is given without a unit then all the units for"+
				" the family will be listed."+
				"\n\n"+
				"If this is not given then a list of available families"+
				" will be shown.",
			param.AltNames("f"),
		)

		unitParam := ps.Add("unit", psetter.String{Value: &ul.uName},
			"the name of the unit to show. If this is given then"+
				" a family name must also be given."+
				" Full details of the unit will be displayed.",
			param.AltNames("u"),
		)

		orderParam := ps.Add("by-name", psetter.Bool{Value: &ul.orderByName},
			"sort the units in alpabetical order not in size order."+
				"\n\n"+
				"This should only be given when listing all the units"+
				" for a single family.",
		)

		ps.Add("tagged",
			unitsetter.TagListAppender{Value: &ul.mustHaveTags},
			"only show units which have the given tag."+
				" This should only be given when listing all the units"+
				" Repetitions of this parameter"+
				" will add to the list of tags that must be present."+
				"\n\n"+
				"This should only be given when listing all the units"+
				" for a single family.",
		)

		ps.Add("not-tagged",
			unitsetter.TagListAppender{Value: &ul.mustNotHaveTags},
			"only show units which do not have the given tag."+
				" for a single family. Repetitions of this parameter"+
				" will add to the list of tags that must be missing."+
				"\n\n"+
				"This should only be given when listing all the units"+
				" for a single family.",
		)

		detailsParam := ps.Add("show-details",
			psetter.Bool{Value: &ul.showDetail},
			"show details when listing."+
				"\n\n"+
				"This should not be given when"+
				" showing details for a single unit.",
			param.AltNames("show-detail", "l"),
		)

		noHdrParam := ps.Add("no-header",
			psetter.Bool{Value: &ul.noHeader},
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

				if ul.hasTagConstraints() {
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

				if ul.hasTagConstraints() {
					return errors.New("constraining the units to list" +
						" only has an effect when listing units in a family")
				}
			}

			return ul.checkTagLists()
		})

		return nil
	}
}

// hasTagConstraints returns true if there are any entries in either of the
// lists of tags to check when constraining the units to show.
func (ul unitlist) hasTagConstraints() bool {
	return len(ul.mustHaveTags) > 0 || len(ul.mustNotHaveTags) > 0
}

// checkTagLists returns an error if the same tag appears in both the list of
// mandatory and forbidden tags
func (ul unitlist) checkTagLists() error {
	for _, mht := range ul.mustHaveTags {
		for _, mnht := range ul.mustNotHaveTags {
			if mht == mnht {
				return fmt.Errorf(
					"Tag %q is in both the mandatory and forbidden tag lists",
					mht)
			}
		}
	}
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
