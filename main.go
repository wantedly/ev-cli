package main

import (
	"github.com/wantedly/ev/cmd"
)

var (
	Version  string
	Revision string
)

func main() {
	cmd.SetVersion(Version, Revision)
	cmd.Execute()
}
