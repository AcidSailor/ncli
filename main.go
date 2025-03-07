package main

import (
	"github.com/acidsailor/ncli/cmd"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
