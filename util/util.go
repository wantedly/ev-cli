package util

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

var (
	jstTimeZone = time.FixedZone("Asia/Tokyo", 9*60*60)
)

func GetWorkingDirectoryName() string {
	if dir, err := os.Getwd(); err == nil {
		return path.Base(dir)
	}
	return ""
}

func GetCurrentBranch() (string, error) {
	// NOTE: `git symbolic-ref --short HEAD` can be used for getting current branch
	// cf. https://qiita.com/sugyan/items/83e060e895fa8ef2038c
	out, err := exec.Command("git", "symbolic-ref", "--short", "HEAD").Output()
	if err != nil {
		return "", err
	}
	s := strings.TrimRight(string(out), "\r\n")
	return strings.TrimRight(s, "\r\n"), err
}

func GetCurrentCommit() (string, error) {
	// NOTE: `git rev-parse HEAD` can be used for getting current commit
	// cf. https://stackoverflow.com/questions/949314/how-to-retrieve-the-hash-for-the-current-commit-in-git
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}
	s := strings.TrimRight(string(out), "\r\n")
	return s, err
}

func CurrentJSTTime() time.Time {
	return time.Now().In(jstTimeZone)
}
