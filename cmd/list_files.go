package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wantedly/ev-cli/aws/s3"
	"github.com/wantedly/ev-cli/config"
	"github.com/wantedly/ev-cli/consts"
	"github.com/wantedly/ev-cli/target"
	"github.com/wantedly/ev-cli/util"
)

var listFilesOpts = struct {
	namespace string
}{}

func init() {
	listFilesCmd := &cobra.Command{
		Use:   "ls-files <target or branch>",
		Short: "List files in a target or branch",
		RunE:  listFiles,
	}

	listFilesCmd.PersistentFlags().StringVarP(&listFilesOpts.namespace, "namespace", "n", util.GetWorkingDirectoryName(), "target namespace name")

	RootCmd.AddCommand(listFilesCmd)
}

func listFiles(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("<target or branch> should be specified!\n")
	}

	t, err := target.Target(args[0], listFilesOpts.namespace)
	if err != nil {
		return err
	}

	fmt.Printf("Showing files in \"%s\" target in \"%s\" namespace\n", t, listFilesOpts.namespace)
	keyPrefix := consts.ReportDir + "/" + listFilesOpts.namespace + "/" + target.ToPath(t) + "/"

	files, err := s3.ListFiles(config.Bucket, keyPrefix)
	if err != nil {
		return err
	}
	for _, f := range files {
		fmt.Printf("%s\n", f)
	}

	return nil
}
