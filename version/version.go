package version

import "fmt"

// Build metadata
var (
	VERSION   = "UNKNOWN"
	COMMIT    = "UNKNOWN"
	BRANCH    = "UNKNOWN"
	BUILDDATE = "UNKNOWN"
)

// Version returns version, branch, commit, build date info
func Version() string {
	return fmt.Sprintf("current version: %s, branch: %s, commit: %s, build date: %s",
		VERSION, BRANCH, COMMIT, BUILDDATE)
}
