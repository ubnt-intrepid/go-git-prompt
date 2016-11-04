package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func getTagOrHash() string {
	tag, _, _ := Communicate("git", "describe", "--exact-match")
	if tag != "" {
		return tag[0 : len(tag)-1]
	}

	hash, _, _ := Communicate("git", "rev-parse", "--short", "HEAD")
	return ":" + strings.TrimSpace(hash[0:len(hash)-1])
}

func parseBranch(line string) (string, bool, bool, int, int, error) {
	var branch string
	var detached, hasremote bool
	var ahead, behind int

	if strings.Contains(line, "no branch") {
		detached = true
		hasremote = false
		branch = getTagOrHash()
	} else if strings.Contains(line, "...") {
		detached = false
		hasremote = true

		splitted := strings.Split(line, " ")
		branch = strings.Split(splitted[1], "...")[0]

		if len(splitted) >= 3 {
			joined := strings.Join(splitted[2:len(splitted)], " ")

			ahead = parsePattern(`ahead (\d+)`, joined)
			behind = parsePattern(`behind (\d+)`, joined)
		}
	} else {
		detached = false
		hasremote = false
		branch = strings.TrimSpace(strings.Split(line, " ")[1])
	}

	return branch, detached, hasremote, ahead, behind, nil
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
	detached  bool
	hasremote bool
	ahead     int
	behind    int
	staged    int
	conflicts int
	changed   int
	untracked int
	stashs    int
}

func newStatus() Status {
	return Status{"", false, false, 0, 0, 0, 0, 0, 0, 0}
}

func (s Status) String() string {
	return fmt.Sprintf("%s %v %v %d %d %d %d %d %d %d", s.branch, s.detached, s.hasremote, s.ahead, s.behind, s.staged, s.conflicts, s.changed, s.untracked, s.stashs)
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

	branch, detached, hasremote, ahead, behind, err := parseBranch(lines[0])
	if err != nil {
		return newStatus(), err
	}
	staged, conflicts, changed, untracked := countChanges(lines[1:len(lines)])

	return Status{branch, detached, hasremote, ahead, behind, staged, conflicts, changed, untracked, len(stashs)}, nil
}
