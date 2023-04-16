package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ravener/chip8/ui"
)

func main() {
	// Parse command line flags.
	flag.Parse()

	// Check if a rom file is passed.
	file := flag.Arg(0)
	if file == "" {
		fmt.Fprintln(os.Stderr, "Usage: chip8 <rom.c8>")
		os.Exit(1)
	}

	// Start the UI.
	ui.Run(file)
}
