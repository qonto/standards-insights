package main

import (
	"fmt"
	"os"

	"github.com/qonto/standards-insights/cmd"
)

func main() {
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
