package main

// unittags

import (
	"fmt"
	"os"
	"sort"

	"github.com/nickwells/col.mod/v3/col"
	"github.com/nickwells/col.mod/v3/col/colfmt"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/twrap.mod/twrap"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
	"github.com/nickwells/unittools/internal/utparams"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// Created: Sat Jul 24 12:59:55 2021

// unittags groups the parameter values for the unittags program
type unittags struct {
	tag         units.Tag
	showDetails bool
}

func main() {
	var ut unittags
	ps := paramset.NewOrDie(addParams(&ut),
		utparams.AddRefUnitlist,
		utparams.AddRefUnitconv,
		versionparams.AddParams,
		param.SetProgramDescription(utparams.ProgDescUnittags),
	)

	ps.Parse()

	if ut.tag == "" {
		listTagNames(ut)
	} else {
		showTagDetails(ut)
	}
}

// maxTagNameLen returns the maximum length of the tag names
func maxTagNameLen(tags []string) int {
	max := 0

	for _, tag := range tags {
		if len(tag) > max {
			max = len(tag)
		}
	}
	return max
}

// listTagNames lists the available tag listTagNames
func listTagNames(ut unittags) {
	tags := units.GetTagNames()

	sort.Strings(tags)

	max := maxTagNameLen(tags)
	extraCols := []*col.Col{}
	if ut.showDetails {
		extraCols = append(extraCols,
			col.New(&colfmt.WrappedString{W: 50}, "Notes"))
	}
	rpt := col.StdRpt(col.New(&colfmt.String{W: max}, "Tag"), extraCols...)

	var err error
	for _, name := range tags {
		vals := []any{name}
		if ut.showDetails {
			vals = append(vals, units.Tag(name).Notes())
		}

		err = rpt.PrintRow(vals...)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Error found while printing the list of tag names: %v\n",
				err)
			os.Exit(1)
		}
	}
}

// showTagDetails displays the details for just the given showTagDetails
func showTagDetails(ut unittags) {
	twc := twrap.NewTWConfOrPanic()

	twc.WrapPrefixed("  Tag: ", string(ut.tag), 0)
	twc.WrapPrefixed("Notes: ", ut.tag.Notes(), 0)
}

// addParams will add parameters to the passed ParamSet
func addParams(ut *unittags) func(ps *param.PSet) error {
	return func(ps *param.PSet) error {
		ps.Add("long", psetter.Bool{Value: &ut.showDetails},
			"show the full details when displaying the tag",
			param.AltNames("l"),
		)

		ps.Add("tag", unitsetter.TagSetter{Value: &ut.tag},
			"show the full details of just this tag",
			param.AltNames("t"),
		)

		return nil
	}
}
