package main

import (
	"bytes"
	"os/exec"
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
