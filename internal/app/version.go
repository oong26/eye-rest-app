package app

const (
	Name = "EyeRest"

	// DefaultVersion follows semantic versioning: MAJOR.MINOR.PATCH.
	DefaultVersion = "1.1.0"
)

func DisplayName(version string) string {
	return Name + " v" + version
}
