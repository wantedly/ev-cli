package target

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/wantedly/ev-cli/aws/s3"
	"github.com/wantedly/ev-cli/config"
	"github.com/wantedly/ev-cli/consts"
	"regexp"
	"strings"
	"time"
)

const (
	validTimestamp     = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+09:00`
	validCommit        = `[0-9a-z]{7}`
	timestampFormat    = "2006-01-02T15:04:05-07:00"
	validMinCommitSize = 7
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
	paths, err := s3.ListPaths(config.Bucket, keyPrefix)
	if err != nil {
		return "", err
	}

	p, err := findLatestMatch(paths, ToPath(branch))
	if err != nil {
		return "", errors.Wrapf(err, "Error in \"%s\" namespace", namespace)
	}
	return PathTo(p), nil
}

func findLatestMatch(paths []string, branch string) (string, error) {
	for i := len(paths) - 1; i >= 0; i-- {
		p := paths[i]
		m := validTarget.FindStringSubmatch(p)
		if len(m) > 0 && (m[1] == branch) {
			return p, nil
		}
	}
	return "", errors.New(fmt.Sprintf("There is no target in \"%s\" branch", branch))
}

func TargetFrom(t time.Time, branch string, commit string) (string, error) {
	c, err := normalizeCommit(commit)
	if err != nil {
		return "", err
	}
	return FormatTime(t) + "-" + branch + "-" + c, nil
}

func normalizeCommit(commit string) (string, error) {
	if len(commit) < validMinCommitSize {
		return "", errors.New(fmt.Sprintf("commithash size is too short: %s", commit))
	}
	return commit[:validMinCommitSize], nil
}

func FormatTime(t time.Time) string {
	return t.Format(timestampFormat)
}

func ValidateTarget(t string) error {
	if validTarget.MatchString(t) {
		return nil
	} else {
		return errors.New(fmt.Sprintf("The format of target is invalid: %s", t))
	}
}

func Branches(paths []string) ([]string, error) {
	visit := map[string]bool{}
	r := []string{}
	for _, p := range paths {
		b, err := branch(p)
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

func branch(path string) (string, error) {
	m := validTarget.FindStringSubmatch(path)
	if len(m) > 0 {
		return PathTo(m[1]), nil
	}
	return "", errors.New(fmt.Sprintf("Invalid target: %s", path))
}

// `/` should not be used in path. So we replace `/` with `\`
func ToPath(s string) string {
	return strings.Replace(s, "/", "\\", -1)
}

func PathTo(s string) string {
	return strings.Replace(s, "\\", "/", -1)
}
