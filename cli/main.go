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
		// Check if it's an ExitError with a specific exit code
		if exitErr, ok := err.(*cmd.ExitError); ok {
			// Print error message only if it's not empty
			if exitErr.Message != "" {
				fmt.Fprintf(os.Stderr, "Error: %v\n", exitErr.Message)
			}
			os.Exit(exitErr.Code)
		}
		// For other errors, print them
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
