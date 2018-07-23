package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wantedly/ev/consts"
	"github.com/wantedly/ev/aws/s3"
	"github.com/wantedly/ev/target"
	"github.com/wantedly/ev/util"
)

var listFilesOpts = struct {
	namespace string
}{}

func init() {
	listFilesCmd := &cobra.Command{
		Use:   "ls-files <target>",
		Short: "List files in a target",
		RunE:  listFiles,
	}

	listFilesCmd.PersistentFlags().StringVarP(&listFilesOpts.namespace, "namespace", "n", util.GetWorkingDirectoryName(), "target namespace name")

	RootCmd.AddCommand(listFilesCmd)
}

func listFiles(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("<target> should be specified!\n")
	}

	t, err := target.Target(args[0], listFilesOpts.namespace)
	if err != nil {
		return err
	}

	fmt.Printf("showing files in \"%s\" target\n", listFilesOpts.namespace+"/"+t)
	keyPrefix := consts.ReportDir + "/" + listFilesOpts.namespace + "/" + t + "/"

	files, err := s3.ListFiles(consts.BucketName, keyPrefix)
	if err != nil {
		return err
	}
	for _, f := range files {
		fmt.Printf("%s\n", f)
	}

	return nil
}
