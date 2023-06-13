package main

// unitlist

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/nickwells/col.mod/v3/col"
	"github.com/nickwells/col.mod/v3/col/colfmt"
	"github.com/nickwells/units.mod/v2/units"
)

// Created: Fri Dec 25 18:42:35 2020

type Prog struct {
	family *units.Family
	uName  string

	mustHaveTags    []units.Tag
	mustNotHaveTags []units.Tag

	orderByName bool
	showDetail  bool
	noHeader    bool
}

// NewProg returns a new Prog instance with the default values set
func NewProg() *Prog {
	return &Prog{}
}

func main() {
	prog := NewProg()
	ps := makeParamSet(prog)

	ps.Parse()

	if prog.family == nil {
		prog.listFamilies()
		return
	}
	if prog.uName == "" {
		prog.listUnits()
		return
	}
	prog.showUnit()
}

// getUnitIDs gets a sorted list of unit getUnitIDs
func (prog Prog) getUnitIDs() []string {
	unitIDs := prog.family.GetUnitNames()

	if prog.orderByName {
		sort.Strings(unitIDs)
	} else {
		sort.Slice(unitIDs, func(i, j int) bool {
			iu, err := prog.family.GetUnit(unitIDs[i])
			if err != nil {
				return false
			}
			ju, err := prog.family.GetUnit(unitIDs[j])
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
func (prog Prog) makeUnitListRpt() *col.Report {
	hdr := col.NewHeaderOrPanic()
	if prog.noHeader {
		hdr = col.NewHeaderOrPanic(col.HdrOptDontPrint)
	}

	if prog.showDetail {
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
func (prog Prog) printUnitRow(rpt *col.Report, uName string) bool {
	u, err := prog.family.GetUnit(uName)
	if err != nil {
		return false
	}

	for _, tag := range prog.mustHaveTags {
		if !u.HasTag(tag) {
			return true
		}
	}
	for _, tag := range prog.mustNotHaveTags {
		if u.HasTag(tag) {
			return true
		}
	}

	if !prog.showDetail {
		err = rpt.PrintRow(uName)
	} else {
		intro := ""
		if uName == prog.family.BaseUnitName() {
			intro = ">>>"
		}

		err = rpt.PrintRow(
			intro, uName, getUnitTags(u), u.ConvFactor(), getUnitNotes(u))
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error found while printing the %q units: %v\n",
			prog.family.Name(), err)
		os.Exit(1)
	}

	return true
}

// listUnits reports on the available units in the given family
func (prog *Prog) listUnits() {
	unitIDs := prog.getUnitIDs()

	rpt := prog.makeUnitListRpt()

	badUnits := []string{}
	for _, uName := range unitIDs {
		if !prog.printUnitRow(rpt, uName) {
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
func (prog Prog) makeFamilyListRpt() *col.Report {
	hdr := col.NewHeaderOrPanic()
	if prog.noHeader {
		hdr = col.NewHeaderOrPanic(col.HdrOptDontPrint)
	}

	if prog.showDetail {
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
func (prog *Prog) printFamilyRow(rpt *col.Report, fName string) {
	var err error

	if prog.showDetail {
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
func (prog *Prog) listFamilies() {
	validFamilies := units.GetFamilyNames()
	sort.Strings(validFamilies)

	rpt := prog.makeFamilyListRpt()

	for _, fName := range validFamilies {
		prog.printFamilyRow(rpt, fName)
	}
}

// hasTagConstraints returns true if there are any entries in either of the
// lists of tags to check when constraining the units to show.
func (prog Prog) hasTagConstraints() bool {
	return len(prog.mustHaveTags) > 0 || len(prog.mustNotHaveTags) > 0
}

// checkTagLists returns an error if the same tag appears in both the list of
// mandatory and forbidden tags
func (prog Prog) checkTagLists() error {
	for _, mht := range prog.mustHaveTags {
		for _, mnht := range prog.mustNotHaveTags {
			if mht == mnht {
				return fmt.Errorf(
					"Tag %q is in both the mandatory and forbidden tag lists",
					mht)
			}
		}
	}
	return nil
}
