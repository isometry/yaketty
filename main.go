package main

import (
	"fmt"
	"os"

	"github.com/isometry/yaketty/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	versionString := fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
	rootCmd := cmd.New(versionString)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
