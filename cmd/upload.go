package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wantedly/ev-cli/aws/s3"
	e "github.com/wantedly/ev-cli/cmd/export"
	"github.com/wantedly/ev-cli/config"
	"github.com/wantedly/ev-cli/consts"
	"github.com/wantedly/ev-cli/target"
	"github.com/wantedly/ev-cli/util"
	"io"
	"os"
	"time"
)

const (
	currentBranch = "[current branch]"
	currentCommit = "[current commit]"
)

var uploadOpts = struct {
	branch         string
	commit         string
	out            string
	metrics        string
	hyperparameter string
	namespace      string
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
	uploadCmd.PersistentFlags().StringVarP(&uploadOpts.hyperparameter, "hyperparameter", "p", "", "target hyper parameter file")
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

	if err := uploadFile(uploadOpts.namespace, t, uploadOpts.out, "evaluate/out.txt"); err != nil {
		return err
	}
	if err := uploadFile(uploadOpts.namespace, t, uploadOpts.metrics, "metrics.json"); err != nil {
		return err
	}
	if err := uploadContext(uploadOpts.namespace, t, uploadOpts.branch, uploadOpts.commit, datetime); err != nil {
		return err
	}
	if err := uploadHyperParameter(uploadOpts.namespace, t, uploadOpts.hyperparameter); err != nil {
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

func uploadFile(namespace, t, src, dst string) error {
	key := consts.ReportDir + "/" + namespace + "/" + target.ToPath(t) + "/" + dst

	file, err := os.Open(src)
	if err != nil {
		return err
	}

	fmt.Printf("uploading %s...\n", src)
	if err = s3.Upload(config.Bucket, key, file); err != nil {
		return err
	}
	return nil
}

func uploadHyperParameter(namespace, t, f string) error {
	var r io.Reader
	if f == "" {
		r = bytes.NewBufferString("{}")
	} else {
		file, err := os.Open(f)
		r = file
		if err != nil {
			return err
		}
	}

	fmt.Printf("uploading %s...\n", "hyperparameter")
	key := consts.ReportDir + "/" + namespace + "/" + target.ToPath(t) + "/hyperparameter.json"
	if err := s3.Upload(config.Bucket, key, r); err != nil {
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
	if err := s3.Upload(config.Bucket, key, bytes.NewBuffer(j)); err != nil {
		return err
	}

	return nil
}
