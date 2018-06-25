package main

import (
	"fmt"
	"os"
)

func main() {
	var cwd, err = os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get current working directory: %s", err)
		os.Exit(1)
	}

	var s = &server{
		Approot: cwd,
		Bind:    ":8080",
	}
	s.Listen()
}
