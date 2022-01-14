package main

import (
	"os"

	"github.com/hapefarts/hapefarts/cmd/v2/internal/cli"
)

var version string

func main() {
	os.Exit((&cli.CLI{
		Version:  version,
		Thinking: false,
	}).Run(os.Args[1:]))
}
