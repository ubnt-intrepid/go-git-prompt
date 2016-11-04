package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func parseBranch(line string) (string, int, int, error) {
	splitted := strings.Split(line, " ")
	branch := strings.Split(splitted[1], "...")[0]

	var ahead, behind int
	if len(splitted) >= 3 {
		joined := strings.Join(splitted[2:len(splitted)], " ")
		ahead = parsePattern(`ahead (\d+)`, joined)
		behind = parsePattern(`behind (\d+)`, joined)
	}

	return branch, ahead, behind, nil
}

func parsePattern(pattern string, s string) int {
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

func countChanges(lines []string) (int, int, int, int) {
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

// Status ...
type Status struct {
	branch    string
	ahead     int
	behind    int
	staged    int
	conflicts int
	changed   int
	untracked int
	stashs    int
}

func newStatus() Status {
	return Status{"", 0, 0, 0, 0, 0, 0, 0}
}

func (s Status) String() string {
	return fmt.Sprintf("%s %d %d %d %d %d %d %d", s.branch, s.ahead, s.behind, s.staged, s.conflicts, s.changed, s.untracked, s.stashs)
}

// GetCurrentStatus ...
func GetCurrentStatus() (Status, error) {
	var lines []string
	{
		stdout, stderr, err := Communicate("git", "status", "--porcelain", "--branch")
		if err != nil {
			return newStatus(), err
		} else if strings.Contains(stderr, "fatal") {
			return newStatus(), errors.New("")
		}
		if len(stdout) > 0 {
			lines = strings.Split(stdout[0:len(stdout)-1], "\n")
		}
	}

	var stashs []string
	{
		stdout, stderr, err := Communicate("git", "stash", "list")
		if err != nil {
			return newStatus(), err
		} else if strings.Contains(stderr, "fatal") {
			return newStatus(), errors.New("")
		}
		if len(stdout) > 0 {
			stashs = strings.Split(stdout[0:len(stdout)-1], "\n")
		}
	}

	branch, ahead, behind, err := parseBranch(lines[0])
	if err != nil {
		return newStatus(), err
	}
	staged, conflicts, changed, untracked := countChanges(lines[1:len(lines)])

	return Status{branch, ahead, behind, staged, conflicts, changed, untracked, len(stashs)}, nil
}
