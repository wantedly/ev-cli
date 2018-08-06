package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wantedly/ev-cli/aws/s3"
	"github.com/wantedly/ev-cli/consts"
	"github.com/wantedly/ev-cli/target"
	"github.com/wantedly/ev-cli/util"
)

var downloadOpts = struct {
	namespace string
}{}

func init() {
	downloadCmd := &cobra.Command{
		Use:   "download <target or branch> <file>",
		Short: "Download a file in a target or branch",
		RunE:  download,
	}

	downloadCmd.PersistentFlags().StringVarP(&downloadOpts.namespace, "namespace", "n", util.GetWorkingDirectoryName(), "target namespace name")

	RootCmd.AddCommand(downloadCmd)
}

func download(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("<target or branch> and <file> should be specified!\n")
	}
	if len(args) <= 1 {
		return errors.New("<file> should be specified!\n")
	}

	t, err := target.Target(args[0], downloadOpts.namespace)
	if err != nil {
		return err
	}
	file := args[1]

	key := consts.ReportDir + "/" + downloadOpts.namespace + "/" + target.ToPath(t) + "/" + file

	bytes, err := s3.Download(consts.BucketName, key)
	if err != nil {
		return err
	}
	fmt.Printf(string(bytes))

	return nil
}
