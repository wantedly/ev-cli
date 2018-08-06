package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/wantedly/ev-cli/aws/session"
	"io"
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
		return []byte{}, wrap(err, bucket, key)
	}
	return buff.Bytes(), nil
}

func Upload(bucket string, key string, r io.Reader) error {
	uploader := s3manager.NewUploader(session.Session())
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   r,
	}
	// NOTE: UploadOutput is not necessary, so _ is used
	_, err := uploader.Upload(upParams)
	return wrap(err, bucket, key)
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

// listObjects returns *s3.ListObjectsV2Output object, which contains all
// Contents and all CommonPrefixes returned by all paginated requests.
// Other properties are same with the first response.
func listObjects(bucket string, prefix string, delimiter string) (*s3.ListObjectsV2Output, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String(delimiter),
	}

	o, err := listObjectsImpl(input)
	if err != nil {
		return nil, wrap(err, bucket, prefix)
	}
	return o, nil
}

func listObjectsImpl(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	cli := s3.New(session.Session())

	// NOTE: Try once
	resp, err := cli.ListObjectsV2(input)
	if err != nil {
		return nil, err
	}

	if !(*resp.IsTruncated) {
		return resp, nil
	}

	// NOTE: response is truncated, so loop until response reach last
	input.SetContinuationToken(*resp.NextContinuationToken)
	for {
		r, err := cli.ListObjectsV2(input)
		if err != nil {
			return nil, err
		}

		resp.Contents = append(resp.Contents, r.Contents...)
		resp.CommonPrefixes = append(resp.CommonPrefixes, r.CommonPrefixes...)

		if !(*r.IsTruncated) {
			// NOTE: Reached to last object
			break
		}
		input.SetContinuationToken(*r.NextContinuationToken)
	}

	return resp, nil
}

func wrap(err error, bucket, key string) error {
	if aerr, ok := err.(awserr.Error); ok {
		return errors.Wrapf(aerr, "Error in \"s3://%s/%s\"", bucket, key)
	}
	return err
}
