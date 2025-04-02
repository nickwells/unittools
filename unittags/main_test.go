package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/param.mod/v6/paramtest"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

// cmpUnitTagsStruct compares the value with the expected value and returns
// an error if they differ
func cmpUnitTagsStruct(iVal, iExpVal any) error {
	val, ok := iVal.(*prog)
	if !ok {
		return errors.New("Bad value: not a pointer to unittags")
	}

	expVal, ok := iExpVal.(*prog)
	if !ok {
		return errors.New("Bad expected value: not a pointer to unittags")
	}

	if val.tag != expVal.tag {
		return fmt.Errorf("The Tag values differ: %q != %q",
			val.tag, expVal.tag)
	}

	if val.showDetails != expVal.showDetails {
		return fmt.Errorf("The ShowDetails values differ: %t != %t",
			val.showDetails, expVal.showDetails)
	}

	return nil
}

// TestAddParams will use the paramtest.Parser to make sure that the
// behaviour of the parameter setting is as expected.
func TestAddParams(t *testing.T) {
	var ut1, ut2, ut3 prog

	testCases := []paramtest.Parser{
		{
			ID:        testhelper.MkID("set long option"),
			Ps:        paramset.NewNoHelpNoExitNoErrRptOrPanic(addParams(&ut1)),
			Val:       &ut1,
			ExpVal:    &prog{tag: "", showDetails: true},
			CheckFunc: cmpUnitTagsStruct,
			Args:      []string{"-long"},
		},
		{
			ID:        testhelper.MkID("set tag"),
			Ps:        paramset.NewNoHelpNoExitNoErrRptOrPanic(addParams(&ut2)),
			Val:       &ut2,
			ExpVal:    &prog{tag: "historic", showDetails: false},
			CheckFunc: cmpUnitTagsStruct,
			Args:      []string{"-tag", "historic"},
		},
		{
			ID:        testhelper.MkID("tag setting error"),
			Ps:        paramset.NewNoHelpNoExitNoErrRptOrPanic(addParams(&ut3)),
			Val:       &ut3,
			ExpVal:    &prog{tag: "", showDetails: false},
			CheckFunc: cmpUnitTagsStruct,
			Args:      []string{"-tag", "hystoric"},
			ExpParseErrors: errutil.ErrMap{
				"tag": []error{
					errors.New(`there is no unit tag called "hystoric".` +
						` Did you mean: "historic"?` +
						"\n" +
						"At: [command line]:" +
						` Supplied Parameter:2: "-tag" "hystoric"`),
				},
			},
		},
	}

	for _, tc := range testCases {
		_ = tc.Test(t)
	}
}
