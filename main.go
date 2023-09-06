package main

import (
	"os"

	"standards/cmd"
)

func main() {
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}
