package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/wantedly/ev/aws/session"
	"strings"
)

func Download(bucket string, key string) ([]byte, error) {
	buff := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(session.Session())
	// NOTE: numBytes is not necessary, so _ is used
	_, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return []byte{}, err
	}
	return buff.Bytes(), nil
}

func ListFiles(bucket string, prefix string) ([]string, error) {
	resp, err := listObjects(bucket, prefix, "")
	if err != nil {
		return []string{}, err
	}

	var r []string
	for _, c := range resp.Contents {
		r = append(r, strings.TrimPrefix(*c.Key, prefix))
	}

	return r, nil
}

func ListPaths(bucket string, prefix string) ([]string, error) {
	resp, err := listObjects(bucket, prefix, "/")
	if err != nil {
		return []string{}, err
	}

	var r []string
	for _, p := range resp.CommonPrefixes {
		k := *p.Prefix
		s := strings.TrimPrefix(k, prefix)
		r = append(r, s[:len(s)-1]) // NOTE: Trim `/` suffix
	}

	return r, nil
}

func listObjects(bucket string, prefix string, delimiter string) (*s3.ListObjectsV2Output, error) {
	cli := s3.New(session.Session())

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String(delimiter),
	}

	resp, err := cli.ListObjectsV2(input)

	return resp, err
}
