package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wantedly/ev-cli/aws/s3"
	"github.com/wantedly/ev-cli/config"
	"github.com/wantedly/ev-cli/consts"
)

var namespacesOpts = struct {
	namespace string
}{}

func init() {
	namespacesCmd := &cobra.Command{
		Use:   "namespaces",
		Short: "List namespaces",
		RunE:  namespaces,
	}

	RootCmd.AddCommand(namespacesCmd)
}

func namespaces(cmd *cobra.Command, args []string) error {
	keyPrefix := consts.ReportDir + "/"

	paths, err := s3.ListPaths(config.Bucket, keyPrefix)
	if err != nil {
		return err
	}
	for _, p := range paths {
		fmt.Printf("%s\n", p)
	}
	return nil
}
