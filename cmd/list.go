package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wantedly/ev/consts"
	"github.com/wantedly/ev/aws/s3"
	"github.com/wantedly/ev/util"
)

var listOpts = struct {
	namespace string
}{}

func init() {
	listCmd := &cobra.Command{
		Use:   "ls",
		Short: "List targets in a namespace",
		RunE:  list,
	}

	listCmd.PersistentFlags().StringVarP(&listOpts.namespace, "namespace", "n", util.GetWorkingDirectoryName(), "target application name")

	RootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	fmt.Printf("showing targets in \"%s\" namespace\n", listOpts.namespace)
	keyPrefix := consts.ReportDir + "/" + listOpts.namespace + "/"

	paths, err := s3.ListPaths(consts.BucketName, keyPrefix)
	if err != nil {
		return err
	}
	for _, p := range paths {
		fmt.Printf("%s\n", p)
	}
	return nil
}
