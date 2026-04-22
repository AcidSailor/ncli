package main

import (
	"fmt"
	"os"

	"github.com/acidsailor/ncli/cmd"
)

var (
	version = "none"
	commit  = "none"
	date    = "none"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
