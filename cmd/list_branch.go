package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wantedly/ev-cli/aws/s3"
	"github.com/wantedly/ev-cli/config"
	"github.com/wantedly/ev-cli/consts"
	"github.com/wantedly/ev-cli/target"
	"github.com/wantedly/ev-cli/util"
)

var listBranchOpts = struct {
	namespace string
}{}

func init() {
	listBranchCmd := &cobra.Command{
		Use:   "ls-branch",
		Short: "List branches in a namespace",
		RunE:  listBranch,
	}

	listBranchCmd.PersistentFlags().StringVarP(&listBranchOpts.namespace, "namespace", "n", util.GetWorkingDirectoryName(), "target application name")

	RootCmd.AddCommand(listBranchCmd)
}

func listBranch(cmd *cobra.Command, args []string) error {
	fmt.Printf("Showing branches in \"%s\" namespace\n", listBranchOpts.namespace)
	keyPrefix := consts.ReportDir + "/" + listBranchOpts.namespace + "/"

	paths, err := s3.ListPaths(config.Bucket, keyPrefix)
	if err != nil {
		return err
	}
	branches, err := target.Branches(paths)
	if err != nil {
		return err
	}
	for _, b := range branches {
		fmt.Printf("%s\n", b)
	}
	return nil
}
