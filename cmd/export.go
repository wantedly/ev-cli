package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	e "github.com/wantedly/ev-cli/cmd/export"
	"github.com/wantedly/ev-cli/util"
)

var exportOpts = struct {
	namespace string
}{}

func init() {
	exportCmd := &cobra.Command{
		Use:   "export <target>",
		Short: "(Used only for debugging ex-export) Export evaluation result files in a target to bigquery",
		RunE:  export,
	}

	exportCmd.PersistentFlags().StringVarP(&exportOpts.namespace, "namespace", "n", util.GetWorkingDirectoryName(), "target application name")

	RootCmd.AddCommand(exportCmd)
}

func export(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("<target> should be specified!\n")
	}

	t := args[0]
	e.PrintStart(t, exportOpts.namespace)
	if err := e.Export(t, exportOpts.namespace); err != nil {
		return err
	}
	e.PrintSuccess(t, exportOpts.namespace)

	return nil
}
