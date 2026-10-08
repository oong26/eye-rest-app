package main

const appName = "EyeRest"

// appVersion follows semantic versioning: MAJOR.MINOR.PATCH.
// Release builds can override it with:
// go build -ldflags "-X main.appVersion=1.2.3"
var appVersion = "1.0.0"

func appDisplayName() string {
	return appName + " v" + appVersion
}
