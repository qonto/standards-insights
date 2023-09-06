package main

import (
	"os"

	"github.com/qonto/standards-insights/cmd"
)

func main() {
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}
