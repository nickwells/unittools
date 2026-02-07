package main

// unittags

import (
	"fmt"
	"os"
	"sort"

	"github.com/nickwells/col.mod/v6/col"
	"github.com/nickwells/col.mod/v6/colfmt"
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/twrap.mod/twrap"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
)

// Created: Sat Jul 24 12:59:55 2021

// prog groups the parameter values for the prog program
type prog struct {
	tag         units.Tag
	showDetails bool
}

// newProg returns a new Prog instance with the default values set
func newProg() *prog {
	return &prog{}
}

func main() {
	prog := newProg()
	ps := makeParamSet(prog)
	ps.Parse()

	if prog.tag == "" {
		prog.listTagNames()
	} else {
		prog.showTagDetails()
	}
}

// maxTagNameLen returns the maximum length of the tag names
func maxTagNameLen(tags []string) int {
	maxLen := 0

	for _, tag := range tags {
		maxLen = max(len(tag), maxLen)
	}

	return maxLen
}

// listTagNames lists the available tag listTagNames
//
//nolint:mnd
func (prog prog) listTagNames() {
	tags := units.GetTagNames()

	sort.Strings(tags)

	maximum := maxTagNameLen(tags)
	extraCols := []*col.Col{}

	if prog.showDetails {
		extraCols = append(extraCols,
			col.New(&colfmt.WrappedString{W: 50}, "Notes"))
	}

	rpt := col.StdRpt(col.New(&colfmt.String{W: maximum}, "Tag"), extraCols...)

	var err error

	for _, name := range tags {
		vals := []any{name}
		if prog.showDetails {
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
func (prog prog) showTagDetails() {
	twc := twrap.NewTWConfOrPanic()

	twc.WrapPrefixed("  Tag: ", string(prog.tag), 0)
	twc.WrapPrefixed("Notes: ", prog.tag.Notes(), 0)
}

// addParams will add parameters to the passed ParamSet
func addParams(prog *prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.Add("long", psetter.Bool{Value: &prog.showDetails},
			"show the full details when displaying the tag",
			param.AltNames("l"),
		)

		ps.Add("tag", unitsetter.TagSetter{Value: &prog.tag},
			"show the full details of just this tag",
			param.AltNames("t"),
		)

		return nil
	}
}
