package main

import "github.com/nickwells/param.mod/v7/param"

// addExamples adds some examples of how the program might be used
func addExamples(ps *param.PSet) error {
	ps.AddExample("unitlist",
		"This will show the available families of units")
	ps.AddExample("unitlist -f temperature",
		"This will show the available units of temperature")
	ps.AddExample("unitlist -f temperature -u K",
		"This will show details of the 'K' unit of temperature")

	return nil
}
