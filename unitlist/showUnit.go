package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nickwells/twrap.mod/twrap"
	"github.com/nickwells/units.mod/units"
)

type prefixedVal struct {
	pfx string
	val string
}

type unitVal struct {
	labels []string
	values []prefixedVal
}

// formulaText returns a textual description for converting the unit into its
// base units
func formulaText(u units.Unit) string {
	formula := ""
	if u.ConvPostAdd != 0 {
		if u.ConvPostAdd > 0 {
			formula += fmt.Sprintf(" subtract %g", u.ConvPostAdd)
		} else {
			formula += fmt.Sprintf(" add %g", -u.ConvPostAdd)
		}
	}
	if u.ConvFactor != 1.0 {
		formula += fmt.Sprintf(" multiply by %g", u.ConvFactor)
	}
	if u.ConvPreAdd != 0 {
		if u.ConvPreAdd > 0 {
			formula += fmt.Sprintf(" subtract %g", u.ConvPreAdd)
		} else {
			formula += fmt.Sprintf(" add %g", -u.ConvPreAdd)
		}
	}
	if formula == "" {
		formula = "no conversion required (already in base unit)"
	}
	return strings.TrimSpace(formula)
}

// aliases returns the unitVal for the aliases of the unit.
func aliases(fName, unitName string) unitVal {
	ud := units.GetUnitDetailsOrPanic(fName)
	aliases := unitVal{labels: []string{"Aliases"}}
	maxAliasNameLen := 0
	for _, alias := range units.GetAliases(fName, unitName) {
		if len(alias) > maxAliasNameLen {
			maxAliasNameLen = len(alias)
		}
	}
	for _, alias := range units.GetAliases(fName, unitName) {
		pfx := fmt.Sprintf("%*s: ", maxAliasNameLen, alias)
		val := ud.Aliases[alias].Notes
		aliases.values = append(aliases.values, prefixedVal{pfx, val})
	}
	return aliases
}

// getUnitName returns the underlying unit name which will be the uName
// unless it's an alias
func getUnitName(fName, uName string) string {
	ud := units.GetUnitDetailsOrPanic(fName)
	if alias, ok := ud.Aliases[uName]; ok {
		return alias.UnitName
	}
	return uName
}

// showUnitHeader reports the family and unit (and shows if it is an alias)
func showUnitHeader(fName, uName string) {
	fmt.Printf("%s/%s", fName, uName)
	unitName := getUnitName(fName, uName)
	if uName != unitName {
		fmt.Printf(" (= %s)", unitName)
	}
	fmt.Println()
}

// showUnit displays full details of the named showUnit
func showUnit(fName, uName string) {
	ud := units.GetUnitDetailsOrPanic(fName)
	u, err := units.GetUnit(fName, uName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q is not a unit of %s\n", uName, fName)
		os.Exit(1)
	}
	showUnitHeader(fName, uName)
	unitName := getUnitName(fName, uName)

	uvList := []unitVal{
		{
			labels: []string{"Abbreviation"},
			values: []prefixedVal{{val: u.Abbrev}},
		},
		{
			labels: []string{"Name"},
			values: []prefixedVal{{val: u.Name}},
		},
		{
			labels: []string{"Plural"},
			values: []prefixedVal{{val: u.NamePlural}},
		},
		aliases(fName, unitName),
		{
			labels: []string{"Notes"},
			values: []prefixedVal{{val: u.Notes}},
		},
		{
			labels: []string{"Base Unit"},
			values: []prefixedVal{{val: ud.Fam.BaseUnitName}},
		},
		{
			labels: []string{"To convert", "to base units"},
			values: []prefixedVal{{val: formulaText(u)}},
		},
	}
	maxLabelLen := 0
	for _, uv := range uvList {
		for _, l := range uv.labels {
			if len(l) > maxLabelLen {
				maxLabelLen = len(l)
			}
		}
	}

	twc := twrap.NewTWConfOrPanic()
	for _, uv := range uvList {

		if len(uv.labels) > 1 {
			for _, l := range uv.labels[:len(uv.labels)-1] {
				fmt.Printf("%*s\n", maxLabelLen, l)
			}
		}
		label := ": "
		if len(uv.labels) > 0 {
			label = uv.labels[len(uv.labels)-1] + ": "
		}
		for _, v := range uv.values {
			twc.WrapPrefixed(fmt.Sprintf("%*s%s", maxLabelLen+2, label, v.pfx),
				v.val,
				0)
			label = ""
		}
	}
}
