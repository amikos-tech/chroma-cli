package main

import (
	"chromacli/cmd"
)

var (
	Version   = "0.0.0-development-build" // Replaced at build time
	BuildDate = "9999-12-31"              // Replace with the actual build date
)

func main() {
	cmd.Execute(Version, BuildDate)
}
