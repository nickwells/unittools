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

// showUnit displays full details of the named Unit
func showUnit(f *units.Family, uName string) {
	u, err := f.GetUnit(uName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q is not a %s\n", uName, f.Description())
		os.Exit(1)
	}

	fmt.Printf("%s/%s", family.Name(), uName)
	unitName := u.ID()
	if uName != unitName {
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
			values: []prefixedVal{{val: f.BaseUnitName()}},
		},
		{
			labels: []string{"To convert", "to base units"},
			values: []prefixedVal{{val: u.ConversionFormula()}},
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
