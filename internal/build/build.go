package build

import "fmt"

var (
	Version   = "development" //nolint
	CommitSHA = "unknown"     //nolint
	Date      = "unknown"     //nolint
)

func VersionMessage() string {
	return fmt.Sprintf("version %s (commit %s) released %s", Version, CommitSHA, Date)
}
