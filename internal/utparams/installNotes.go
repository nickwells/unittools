package utparams

// installNotes returns text describing how to install the named program
func installNotes(pName string) string {
	return "\n\n" +
		"To get this program:" +
		"\n\n" +
		"go install github.com/nickwells/unittools/" + pName + "@latest"
}
