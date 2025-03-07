package main

import "github.com/nickwells/param.mod/v6/param"

// addExamples adds some examples of how the program might be used
func addExamples(ps *param.PSet) error {
	ps.AddExample("unitconv -from pint -to litre",
		"This will show how many litres in a pint")
	ps.AddExample("unitconv -from chain -to mile",
		"This will show how many chains in a mile")
	ps.AddExample("unitconv -from chain -to mile -val 80",
		"This will show 80 chains in miles")
	ps.AddExample("unitconv -from chain -to m -val 80",
		"This will show 80 chains in metres")
	ps.AddExample("unitconv -from chain -to m -val 80 -just-val",
		"This will show 80 chains in metres. Only the value is "+
			"shown and no surrounding explanatory text")
	ps.AddExample("unitconv -from chain -to m -val 80 -roughly",
		"This will show 80 chains in metres. The value is "+
			"adjusted to show the nearest multiple of 5 or 10")
	return nil
}
