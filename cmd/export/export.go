package export

import (
	"bytes"
	"fmt"
	"github.com/wantedly/ev-cli/aws/s3"
	"github.com/wantedly/ev-cli/consts"
	"github.com/wantedly/ev-cli/target"
)

func Export(t, namespace string) error {
	if err := target.ValidateTarget(t); err != nil {
		return err
	}

	key := consts.TriggerDir + "/" + namespace + "/" + target.ToPath(t)

	if err := s3.Upload(consts.BucketName, key, new(bytes.Buffer)); err != nil {
		return err
	}

	return nil
}

func PrintStart(t, namespace string) {
	fmt.Printf("Export \"%s\" target in \"%s\" namespace\n", t, namespace)
	fmt.Printf("exporting...\n")
}

func PrintSuccess(t, namespace string) {
	fmt.Printf("Success! Export of \"%s\" target in \"%s\" namespace is triggered!\n", t, namespace)
}
