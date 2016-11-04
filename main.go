package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func communicate(name string, arg ...string) (string, string, error) {
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

func _parseBranch(line string) (string, int, int, error) {
	splitted := strings.Split(line, " ")
	branch := strings.Split(splitted[1], "...")[0]

	var ahead, behind int
	if len(splitted) >= 3 {
		joined := strings.Join(splitted[2:len(splitted)], " ")
		ahead = _parsePattern(`ahead (\d+)`, joined)
		behind = _parsePattern(`behind (\d+)`, joined)
	}

	return branch, ahead, behind, nil
}

func _parsePattern(pattern string, s string) int {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		ahead, err := strconv.Atoi(matches[1])
		if err == nil {
			return ahead
		}
	}
	return 0
}

func _countChanges(lines []string) (int, int, int, int) {
	staged, conflicts, changed, untracked := 0, 0, 0, 0
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		idxStatus := line[0]
		wtStatus := line[1]

		// " M  hoge.txt" , "AM ahoo.png" , ...
		if wtStatus != ' ' && wtStatus != '?' {
			changed++
		}

		// "MT hoge.cpp" , "A  fuga.txt" , ...
		if idxStatus != ' ' && idxStatus != '?' {
			staged++
		}

		// "?? hoge.txt", ...
		if idxStatus == '?' && wtStatus == '?' {
			untracked++
		}

		// "UU hogehoge.txt" ...
		if idxStatus == 'U' && wtStatus == 'U' {
			conflicts++
		}
	}

	return staged, conflicts, changed, untracked
}

type status struct {
	branch    string
	ahead     int
	behind    int
	staged    int
	conflicts int
	changed   int
	untracked int
}

func newStatus() status {
	return status{"", 0, 0, 0, 0, 0, 0}
}

func (s status) String() string {
	return fmt.Sprintf("%s %d %d %d %d %d %d", s.branch, s.ahead, s.behind, s.staged, s.conflicts, s.changed, s.untracked)
}

func getCurrentStatus() (status, error) {
	stdout, stderr, err := communicate("git", "status", "--porcelain", "--branch")
	if err != nil {
		return newStatus(), err
	} else if strings.Contains(stderr, "fatal") {
		return newStatus(), errors.New("")
	}
	lines := strings.Split(stdout, "\n")

	branch, ahead, behind, err := _parseBranch(lines[0])
	if err != nil {
		return newStatus(), err
	}
	staged, conflicts, changed, untracked := _countChanges(lines[1:len(lines)])

	return status{branch, ahead, behind, staged, conflicts, changed, untracked}, nil
}

func main() {
	status, err := getCurrentStatus()
	if err != nil {
		return
	}
	fmt.Println(status)
}
