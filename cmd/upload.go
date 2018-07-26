package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wantedly/ev-cli/aws/s3"
	e "github.com/wantedly/ev-cli/cmd/export"
	"github.com/wantedly/ev-cli/consts"
	"github.com/wantedly/ev-cli/target"
	"github.com/wantedly/ev-cli/util"
	"os"
	"path/filepath"
	"time"
)

const (
	currentBranch = "[current branch]"
	currentCommit = "[current commit]"
)

var uploadOpts = struct {
	branch    string
	commit    string
	out       string
	metrics   string
	namespace string
}{}

func init() {
	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload evaluation result files as a target and export it to bigquery",
		Long: `Upload evaluation result files as a target and export it to bigquery

Example.
$ ev upload --out ./out.txt \
            --metrics ./metrics.json \
            --namespace yashima-category-python`,
		RunE: upload,
	}

	uploadCmd.PersistentFlags().StringVarP(&uploadOpts.branch, "branch", "b", currentBranch, "target branch name")
	uploadCmd.PersistentFlags().StringVarP(&uploadOpts.commit, "commit", "c", currentCommit, "target commithash")
	uploadCmd.PersistentFlags().StringVarP(&uploadOpts.out, "out", "o", "./out.txt", "target out file")
	uploadCmd.PersistentFlags().StringVarP(&uploadOpts.metrics, "metrics", "m", "./metrics.json", "target metrics file")
	uploadCmd.PersistentFlags().StringVarP(&uploadOpts.namespace, "namespace", "n", util.GetWorkingDirectoryName(), "target namespace name")

	RootCmd.AddCommand(uploadCmd)
}

func upload(cmd *cobra.Command, args []string) error {
	if uploadOpts.branch == currentBranch {
		out, err := util.GetCurrentBranch()
		if err != nil {
			return err
		}
		uploadOpts.branch = out
	}
	if uploadOpts.commit == currentCommit {
		out, err := util.GetCurrentCommit()
		if err != nil {
			return err
		}
		uploadOpts.commit = out
	}

	datetime := util.CurrentJSTTime()
	t, err := target.TargetFrom(datetime, uploadOpts.branch, uploadOpts.commit)
	if err != nil {
		return err
	}
	fmt.Printf("Upload files as \"%s\" target to \"%s\" namespace\n", t, uploadOpts.namespace)

	if err := uploadFile(uploadOpts.namespace, t, "evaluate", uploadOpts.out); err != nil {
		return err
	}
	if err := uploadFile(uploadOpts.namespace, t, "", uploadOpts.metrics); err != nil {
		return err
	}
	if err := uploadContext(uploadOpts.namespace, t, uploadOpts.branch, uploadOpts.commit, datetime); err != nil {
		return err
	}
	fmt.Printf("Success! Files in \"%s\" target are uploaded!\n", t)

	// NOTE: Export after upload
	e.PrintStart(t, exportOpts.namespace)
	if err := e.Export(t, uploadOpts.namespace); err != nil {
		return err
	}
	e.PrintSuccess(t, exportOpts.namespace)

	return nil
}

func uploadFile(namespace, t, dir, f string) error {
	keyPrefix := consts.ReportDir + "/" + namespace + "/" + target.ToPath(t)
	_, filename := filepath.Split(f)
	var key string
	if dir == "" {
		key = keyPrefix + "/" + filename
	} else {
		key = keyPrefix + "/" + dir + "/" + filename
	}

	file, err := os.Open(f)
	if err != nil {
		return err
	}

	fmt.Printf("uploading %s...\n", f)
	if err = s3.Upload(consts.BucketName, key, file); err != nil {
		return err
	}
	return nil
}

type Context struct {
	Datetime   string `json:"datetime"`
	Branch     string `json:"branch"`
	Commithash string `json:"commithash"`
}

func uploadContext(namespace, t, branch, commithash string, datetime time.Time) error {
	c := Context{
		Datetime:   target.FormatTime(datetime),
		Branch:     branch,
		Commithash: commithash,
	}
	j, err := json.Marshal(c)
	if err != nil {
		return err
	}

	filename := "context.json"
	key := consts.ReportDir + "/" + namespace + "/" + target.ToPath(t) + "/" + filename
	fmt.Printf("uploading %s...\n", filename)
	if err := s3.Upload(consts.BucketName, key, bytes.NewBuffer(j)); err != nil {
		return err
	}

	return nil
}
