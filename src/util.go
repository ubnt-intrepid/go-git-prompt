package prompt

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// Communicate ...
func Communicate(name string, arg ...string) (string, string, error) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		return "", "", err
	}
	stdout := outbuf.String()
	stderr := errbuf.String()

	return stdout, stderr, nil
}

// GetLines ...
func GetLines(name string, arg ...string) ([]string, error) {
	var lines []string
	stdout, stderr, err := Communicate(name, arg...)
	if err != nil {
		return []string{}, err
	} else if strings.Contains(stderr, "fatal") {
		return []string{}, errors.New(stderr)
	}
	if len(stdout) > 0 {
		lines = strings.Split(stdout[0:len(stdout)-1], "\n")
	}

	return lines, nil
}
