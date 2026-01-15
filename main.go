package main

import (
	"fmt"
	"os"

	"github.com/GoExec/internal/core"
)

func main() {
	// Validate modules at startup
	if err := core.InitializeModules(); err != nil {
		fmt.Fprintf(os.Stderr, "Module initialization error: %v\n", err)
		os.Exit(1)
	}

	core.ExecuteCLI()
}
