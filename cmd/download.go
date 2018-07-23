package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wantedly/ev/consts"
	"github.com/wantedly/ev/aws/s3"
	"github.com/wantedly/ev/target"
	"github.com/wantedly/ev/util"
	"strings"
)

var downloadOpts = struct {
	namespace string
}{}

func init() {
	downloadCmd := &cobra.Command{
		Use:   "download <target> <file>",
		Short: "Download a file in a target",
		RunE:  download,
	}

	downloadCmd.PersistentFlags().StringVarP(&downloadOpts.namespace, "namespace", "n", util.GetWorkingDirectoryName(), "target namespace name")

	RootCmd.AddCommand(downloadCmd)
}

func download(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("<target> and <file> should be specified!\n")
	}
	if len(args) <= 1 {
		return errors.New("<file> should be specified!\n")
	}

	t, err := target.Target(args[0], downloadOpts.namespace)
	if err != nil {
		return err
	}
	file := args[1]

	key := consts.ReportDir + "/" + downloadOpts.namespace + "/" + t + "/" + file

	bytes, err := s3.Download(consts.BucketName, key)
	if err != nil {
		if err.Error()[:10] == "NoSuchKey:" {
			// NOTE: Handle NoSuchKey error
			s := strings.TrimPrefix(key, consts.ReportDir+"/"+downloadOpts.namespace+"/")
			return errors.New(fmt.Sprintf("NoSuchKey: \"%s\" does not exist", s))
		}
		return err
	}
	fmt.Printf(string(bytes))

	return nil
}
