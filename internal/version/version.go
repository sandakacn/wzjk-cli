package version

// Version is the current version of wzjkctl
// This is set at build time using -ldflags
var Version = "dev"

// GetVersion returns the current version
func GetVersion() string {
	return Version
}
