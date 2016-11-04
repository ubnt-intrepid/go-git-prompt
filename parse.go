package main

import (
	"regexp"
	"strconv"
	"strings"
)

func getTagOrHash() string {
	tag, _, _ := Communicate("git", "describe", "--exact-match")
	if tag != "" {
		return strings.TrimSpace(tag[0 : len(tag)-1])
	}

	hash, _, _ := Communicate("git", "rev-parse", "--short", "HEAD")
	return strings.TrimSpace(hash[0 : len(hash)-1])
}

// ParseBranchLine ...
func ParseBranchLine(line string) (string, bool, bool, int, int, error) {
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

// CollectChanges ...
func CollectChanges(lines []string) (int, int, int, int) {
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
