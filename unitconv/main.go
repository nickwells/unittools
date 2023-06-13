package main

// units

import (
	"fmt"
	"os"

	"github.com/nickwells/mathutil.mod/v2/mathutil"
	"github.com/nickwells/units.mod/v2/units"
)

// Created: Sat Aug 29 16:52:07 2020

// Prog holds the values describing the conversion to perform
type Prog struct {
	unitFrom units.Unit
	unitTo   units.Unit

	val float64

	justVal        bool
	roughly        bool
	roughPrecision float64
}

// NewProg returns a new Prog instance with the default values set
func NewProg() *Prog {
	return &Prog{
		val: 1,
	}
}

func main() {
	prog := NewProg()
	ps := makeParamSet(prog)

	ps.Parse()

	v := units.ValUnit{V: prog.val, U: prog.unitFrom}
	converted, err := v.Convert(prog.unitTo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if prog.roughly {
		converted.V = mathutil.Roughly(converted.V, prog.roughPrecision)
	}

	if prog.justVal {
		fmt.Println(converted.V)
	} else {
		fmt.Println(v, "=", converted)
	}
}
