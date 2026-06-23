package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/mathutil.mod/v2/mathutil"
	"github.com/nickwells/units.mod/v2/units"
	"github.com/nickwells/verbose.mod/verbose"
)

const (
	esBadConversion = 1 + iota
)

type converted struct {
	vu              units.ValUnit
	absWholeNumDiff float64
	absLogVal       float64
}

// prog holds program parameters and status
type prog struct {
	exitStatus int
	stack      *verbose.Stack
	// parameters
	unitFamily *units.Family

	unitFromName string
	unitToNames  []string

	unitFrom units.Unit
	unitTo   []units.Unit

	val float64

	nearestVal        bool
	nearestCount      int
	nearestPrecision  float64
	nearestIgnoreTags []units.Tag

	justVal        bool
	roughly        bool
	roughPrecision float64

	displayWidth int
	displayPrec  int
}

// newProg returns a new Prog instance with the default values set
func newProg() *prog {
	const (
		dfltDisplayPrec = 6

		dfltNearestCount     = 5
		dfltNearestPrecision = 0.01
	)

	return &prog{
		stack: &verbose.Stack{},

		val:          1,
		displayWidth: 0,
		displayPrec:  dfltDisplayPrec,

		nearestCount:     dfltNearestCount,
		nearestPrecision: dfltNearestPrecision,
	}
}

// setExitStatus sets the exit status to the new value. It will not do this
// if the exit status has already been set to a non-zero value.
func (prog *prog) setExitStatus(es int) {
	if prog.exitStatus == 0 {
		prog.exitStatus = es
	}
}

// getUnitFrom populates the unitFrom member
func (prog *prog) getUnitFrom() error {
	var err error

	if prog.unitFamily != nil {
		prog.unitFrom, err = prog.unitFamily.GetUnit(prog.unitFromName)
		return err
	}

	fNames := []string{}

	for _, f := range units.GetFamilies() {
		u, err := f.GetUnit(prog.unitFromName)
		if err == nil {
			prog.unitFamily = f
			prog.unitFrom = u

			fNames = append(fNames, f.Name())
		}
	}

	switch len(fNames) {
	case 0:
		return fmt.Errorf(
			"there are %d unit-families with a unit called %q: %s",
			len(fNames),
			prog.unitFromName,
			english.JoinQuoted(fNames, ", ", " and "))
	case 1:
		return nil
	default:

	}

	return nil
}

// calcAbsWholeNumDiff calculates the absolute difference between the value and
// the nearest whole number.
func calcAbsWholeNumDiff(v float64) float64 {
	multiples := []float64{2, 3, 4, 5, 8, 10}
	wholeNum := math.Round(v)
	awnd := math.Abs(v - wholeNum)

	for _, m := range multiples {
		vm := v * m
		wholeNum := math.Round(vm)
		altAwnd := math.Abs(vm - wholeNum)
		awnd = min(altAwnd, awnd)
	}

	return awnd
}

// calcAbsLog calculates the absolute log value of the supplied value. This
// will give a value that is smallest when the value is equal to 1.
func calcAbsLog(v float64) float64 {
	return math.Abs(math.Log(v))
}

// cmpAbsLog returns -1, 0 or 1 depending on whether the absLogVal of 'a'
// and 'b' are less than, equal to or greater than each other.
func cmpAbsLog(a, b converted) int {
	if a.absLogVal < b.absLogVal {
		return -1
	}

	if a.absLogVal > b.absLogVal {
		return 1
	}

	return 0
}

// makeCmpConvertedFunc returns a function that will compare the two
// converted values firstly by how close each is to a whole number (note that
// the value also includes several fractions - see calcAbsWholeNumDiff). Then
// if they are the same or only differ by a small amount they are compared by
// how close they are to one.
//
// This is a generated function so that the small difference value can use
// the nearestPrecision value from the prog struct.
func (prog *prog) makeCmpConvertedFunc() func(converted, converted) int {
	return func(a, b converted) int {
		if a.absWholeNumDiff < b.absWholeNumDiff {
			if b.absWholeNumDiff-a.absWholeNumDiff < prog.nearestPrecision {
				return cmpAbsLog(a, b)
			}

			return -1
		}

		if a.absWholeNumDiff > b.absWholeNumDiff {
			if a.absWholeNumDiff-b.absWholeNumDiff < prog.nearestPrecision {
				return cmpAbsLog(a, b)
			}

			return 1
		}

		return cmpAbsLog(a, b)
	}
}

// findNearestVals populates the unitToNames slice with names such that the
// converted values are the closest to small, preferably whole-number
// values. The units are ordered by small, whole numbers first and fractional
// values second.
func (prog *prog) findNearestVals() error {
	if err := prog.getUnitFrom(); err != nil {
		return err
	}

	allUnits := prog.unitFamily.GetUnits()
	unitVals := make([]converted, 0, len(allUnits))
	fromVal := units.ValUnit{V: prog.val, U: prog.unitFrom}

	for _, u := range allUnits {
		vu := fromVal.ConvertOrPanic(u)
		c := converted{
			vu:              vu,
			absWholeNumDiff: calcAbsWholeNumDiff(vu.V),
			absLogVal:       calcAbsLog(vu.V),
		}
		unitVals = append(unitVals, c)
	}

	slices.SortFunc(unitVals, prog.makeCmpConvertedFunc())

AvailableUnits:
	for _, c := range unitVals {
		if units.Equals(c.vu.U, prog.unitFrom) {
			continue AvailableUnits
		}

		for _, tag := range prog.nearestIgnoreTags {
			if c.vu.U.HasTag(tag) {
				continue AvailableUnits
			}
		}

		prog.unitTo = append(prog.unitTo, c.vu.U)
		prog.unitToNames = append(prog.unitToNames, c.vu.U.ID())

		if len(prog.unitTo) >= prog.nearestCount {
			break AvailableUnits
		}
	}

	return nil
}

// showNearest shows the alternative units most likely to be the value.
func (prog *prog) showNearest(v units.ValUnit, fmtStr string) {
	for i, unitTo := range prog.unitTo {
		converted, err := v.Convert(unitTo)
		if err != nil {
			fmt.Println(err)
			prog.setExitStatus(esBadConversion)

			return
		}

		if prog.roughly {
			converted.V = mathutil.Roughly(converted.V, prog.roughPrecision)
		}

		fmt.Printf(fmtStr, converted, prog.unitToNames[i])
	}
}

// formatString returns the format string to display a ValUnit
func (prog *prog) formatString() string {
	fmtIntro := "%" +
		strconv.Itoa(prog.displayWidth) +
		"." +
		strconv.Itoa(prog.displayPrec)

	if prog.justVal {
		return fmtIntro + "f"
	}

	return fmtIntro + "u"
}

// run is the starting point for the program, it is called from main()
// after the command-line parameters have been parsed.
func (prog *prog) run() {
	v := units.ValUnit{V: prog.val, U: prog.unitFrom}

	fmtStr := prog.formatString()

	var s string
	if !prog.justVal {
		s = fmt.Sprintf(fmtStr+" = ", v)
		fmt.Println(s)
	}

	indent := strings.Repeat(" ", len(s))

	if prog.nearestVal {
		prog.showNearest(v, indent+fmtStr+"\t%s\n")

		return
	}

	for i, unitTo := range prog.unitTo {
		converted, err := v.Convert(unitTo)
		if err != nil {
			fmt.Println(err)
			prog.setExitStatus(esBadConversion)

			return
		}

		if prog.roughly {
			converted.V = mathutil.Roughly(converted.V, prog.roughPrecision)
		}

		if i != len(prog.unitTo)-1 {
			intPart := math.Floor(converted.V)
			fracPart := converted.V - intPart
			converted.V = intPart
			backVal := units.ValUnit{V: fracPart, U: unitTo}

			convertedBack, err := backVal.Convert(prog.unitFrom)
			if err != nil {
				fmt.Println(err)
				prog.setExitStatus(esBadConversion)

				return
			}

			v.V = convertedBack.V
		}

		fmt.Printf(fmtStr, converted)
		fmt.Println()
	}
}
