package main

// units

import (
	"fmt"
	"math"
	"os"

	"github.com/nickwells/mathutil.mod/v2/mathutil"
	"github.com/nickwells/units.mod/v2/units"
)

// Created: Sat Aug 29 16:52:07 2020

// prog holds the values describing the conversion to perform
type prog struct {
	unitFamily *units.Family

	unitFromName string
	unitToNames  []string

	unitFrom units.Unit
	unitTo   []units.Unit

	val float64

	justVal        bool
	roughly        bool
	roughPrecision float64
}

// newProg returns a new Prog instance with the default values set
func newProg() *prog {
	return &prog{
		val: 1,
	}
}

func main() {
	prog := newProg()
	ps := makeParamSet(prog)

	ps.Parse()

	v := units.ValUnit{V: prog.val, U: prog.unitFrom}

	if !prog.justVal {
		fmt.Print(v, "=")
	}

	sep := ""

	for i, unitTo := range prog.unitTo {
		converted, err := v.Convert(unitTo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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
				os.Exit(1)
			}

			v.V = convertedBack.V
		}

		fmt.Print(sep)

		if prog.justVal {
			sep = " "

			fmt.Print(converted.V)
		} else {
			sep = ", "

			fmt.Print(converted)
		}
	}

	fmt.Println()
}
