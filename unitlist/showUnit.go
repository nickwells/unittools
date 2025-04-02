package main

import (
	"fmt"
	"os"

	"github.com/nickwells/twrap.mod/twrap"
	"github.com/nickwells/units.mod/v2/units"
)

type prefixedVal struct {
	pfx string
	val string
}

type unitVal struct {
	labels []string
	values []prefixedVal
}

// aliasList returns a list of prefixed values for the aliasList of the unit.
func aliasList(aliases map[string]string) []prefixedVal {
	rval := []prefixedVal{}
	maxAliasNameLen := 0

	for alias := range aliases {
		if len(alias) > maxAliasNameLen {
			maxAliasNameLen = len(alias)
		}
	}

	for alias, notes := range aliases {
		rval = append(rval,
			prefixedVal{
				pfx: fmt.Sprintf("%*s: ", maxAliasNameLen, alias),
				val: notes,
			})
	}

	return rval
}

// maxLabelLen calculates the length of the longest label
func maxLabelLen(uvList []unitVal) int {
	maximum := 0

	for _, uv := range uvList {
		for _, l := range uv.labels {
			if len(l) > maximum {
				maximum = len(l)
			}
		}
	}

	return maximum
}

func unitTags(u units.Unit) string {
	tags := ""
	sep := ""

	for _, t := range u.Tags() {
		tags += sep + string(t)
		sep = ", "
	}

	return tags
}

// showUnit displays full details of the named Unit
func (prog prog) showUnit() {
	u, err := prog.family.GetUnit(prog.uName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q is not a %s\n",
			prog.uName, prog.family.Description())
		os.Exit(1)
	}

	fmt.Printf("%s/%s", prog.family.Name(), prog.uName)

	unitName := u.ID()
	if prog.uName != unitName {
		fmt.Printf(" (= %s)", unitName)
	}

	fmt.Println()

	uvList := []unitVal{
		{
			labels: []string{"Abbreviation"},
			values: []prefixedVal{{val: u.Abbrev()}},
		},
		{
			labels: []string{"Name"},
			values: []prefixedVal{{val: u.Name()}},
		},
		{
			labels: []string{"Plural"},
			values: []prefixedVal{{val: u.NamePlural()}},
		},
		{
			labels: []string{"Aliases"},
			values: aliasList(u.Aliases()),
		},
		{
			labels: []string{"Notes"},
			values: []prefixedVal{{val: u.Notes()}},
		},
		{
			labels: []string{"Base Unit"},
			values: []prefixedVal{{val: prog.family.BaseUnitName()}},
		},
		{
			labels: []string{"To convert", "from base units"},
			values: []prefixedVal{{val: u.ConversionFormula()}},
		},
		{
			labels: []string{"Unit tags"},
			values: []prefixedVal{{val: unitTags(u)}},
		},
	}

	maxLabelLen := maxLabelLen(uvList)
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
			twc.WrapPrefixed(
				fmt.Sprintf("%*s%s",
					maxLabelLen+2, label, //nolint:mnd
					v.pfx),
				v.val,
				0)

			label = ""
		}
	}
}
