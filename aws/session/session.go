package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	sess *session.Session
)

const (
	DefaultRegion = "ap-northeast-1"
)

func Session() *session.Session {
	if sess == nil {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{Region: aws.String(DefaultRegion)},
		}))
	}
	return sess
}
