package main

// unittags

import (
	"fmt"
	"os"
	"sort"

	"github.com/nickwells/col.mod/v3/col"
	"github.com/nickwells/col.mod/v3/col/colfmt"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/twrap.mod/twrap"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/unitsetter.mod/v4/unitsetter"
)

// Created: Sat Jul 24 12:59:55 2021

// Prog groups the parameter values for the Prog program
type Prog struct {
	tag         units.Tag
	showDetails bool
}

// NewProg returns a new Prog instance with the default values set
func NewProg() *Prog {
	return &Prog{}
}

func main() {
	prog := NewProg()
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
	max := 0

	for _, tag := range tags {
		if len(tag) > max {
			max = len(tag)
		}
	}
	return max
}

// listTagNames lists the available tag listTagNames
func (prog Prog) listTagNames() {
	tags := units.GetTagNames()

	sort.Strings(tags)

	max := maxTagNameLen(tags)
	extraCols := []*col.Col{}
	if prog.showDetails {
		extraCols = append(extraCols,
			col.New(&colfmt.WrappedString{W: 50}, "Notes"))
	}
	rpt := col.StdRpt(col.New(&colfmt.String{W: max}, "Tag"), extraCols...)

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
func (prog Prog) showTagDetails() {
	twc := twrap.NewTWConfOrPanic()

	twc.WrapPrefixed("  Tag: ", string(prog.tag), 0)
	twc.WrapPrefixed("Notes: ", prog.tag.Notes(), 0)
}

// addParams will add parameters to the passed ParamSet
func addParams(prog *Prog) param.PSetOptFunc {
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
