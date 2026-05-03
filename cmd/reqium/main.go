package main

import (
	"fmt"
	"os"

	"reqium/internal/cli"
)

func main() {
	root := cli.NewRootCommand()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
