package main

import (
	"fmt"
	"os"

	"github.com/easel/ddx/cmd"
)

var version = "dev"
var commit = "unknown"
var date = "unknown"

func main() {
	// Set version info for cobra
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
} // Test hooks
