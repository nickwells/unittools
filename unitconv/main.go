package main

// units

import (
	"os"
)

// Created: Sat Aug 29 16:52:07 2020

func main() {
	prog := newProg()
	ps := makeParamSet(prog)
	ps.Parse()

	prog.run()
	os.Exit(prog.exitStatus)
}
