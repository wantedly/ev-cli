package main

import (
	"github.com/wantedly/ev-cli/cmd"
)

var (
	Version  string
	Revision string
)

func main() {
	cmd.SetVersion(Version, Revision)
	cmd.Execute()
}
