package main

import "eye-rest-app/internal/app"

var version = app.DefaultVersion

func main() {
	app.Run(version)
}
