package target

import (
	"errors"
	"fmt"
	"github.com/wantedly/ev/consts"
	"github.com/wantedly/ev/aws/s3"
	"regexp"
)

const (
	validTimestamp = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+09:00`
	validCommit    = `[1-9a-z]{7}`
)

var (
	// Ex. 2018-07-22T02:06:04+09:00-master-9e75514/
	validTarget = regexp.MustCompile(`^` + validTimestamp + `-(?P<branch>.+)-` + validCommit + `$`)
)

func Target(t string, namespace string) (string, error) {
	if validTarget.MatchString(t) {
		return t, nil
	}

	// NOTE: t is branch
	branch := t
	keyPrefix := consts.ReportDir + "/" + namespace + "/"
	targets, err := s3.ListPaths(consts.BucketName, keyPrefix)
	if err != nil {
		return "", err
	}

	tt, err := findLatestMatch(targets, branch)
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s in \"%s\" namespace", err.Error(), namespace))
	}
	return tt, nil
}

func Branch(t string) (string, error) {
	m := validTarget.FindStringSubmatch(t)
	if len(m) > 0 {
		return m[1], nil
	}
	return "", errors.New(fmt.Sprintf("Invalid target: %s", t))
}

func Branches(targets []string) ([]string, error) {
	visit := map[string]bool{}
	r := []string{}
	for _, t := range targets {
		b, err := Branch(t)
		if err != nil {
			return []string{}, err
		}
		if !visit[b] {
			visit[b] = true
			r = append(r, b)
		}
	}
	return r, nil
}

func findLatestMatch(targets []string, branch string) (string, error) {
	for i := len(targets) - 1; i >= 0; i-- {
		t := targets[i]
		m := validTarget.FindStringSubmatch(t)
		if len(m) > 0 && (m[1] == branch) {
			return t, nil
		}
	}
	return "", errors.New(fmt.Sprintf("There is no target in \"%s\" branch", branch))
}
