package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SetVersion(version, revision string) {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of ev",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ev %s (rev=%s)\n", version, revision)
		},
	})
}
