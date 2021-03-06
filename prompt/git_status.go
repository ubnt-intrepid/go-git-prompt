package prompt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ubnt-intrepid/go-git-prompt/color"
)

// Status ...
type GitStatus struct {
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

func newGitStatus() GitStatus {
	return GitStatus{"", false, false, 0, 0, 0, 0, 0, 0, 0}
}

// Format ...
func (s *GitStatus) Prompt(color color.Colored) string {
	ret := color.Yellow("(")

	// branch
	if s.detached {
		ret += color.Cyan("(%s)", s.branch)
	} else {
		ret += color.Cyan("%s", s.branch)
	}
	if s.hasremote {
		if s.ahead > 0 && s.behind > 0 {
			ret += color.Yellow(" A%d B%d", s.ahead, s.behind)
		} else if s.ahead > 0 {
			ret += color.Green(" A%d", s.ahead)
		} else if s.behind > 0 {
			ret += color.Red(" B%d", s.behind)
		} else {
			ret += color.Cyan(" =")
		}
	}

	if s.staged > 0 || s.changed > 0 || s.conflicts > 0 || s.untracked > 0 {
		ret += color.Blue(" +%d", s.staged)
		ret += color.Blue(" ~%d", s.changed)
		ret += color.Blue(" ?%d", s.untracked)
		ret += color.Blue(" !%d", s.conflicts)
	}

	if s.stashs > 0 {
		ret += color.Green(" | s%d", s.stashs)
	}

	ret += color.Yellow(")")

	return ret
}

func getStashCount() (int, error) {
	// check if the repository has stash(es).
	_, _, err := Communicate("git", "rev-parse", "--vefify", "--quiet", "refs/stash")
	if err != nil {
		return 0, nil
	}

	stdout, stderr, err := Communicate("git", "log", "--format=\"%%gd: %%gs\"", "-g", "--first-parent", "-m", "refs/stash", "--")
	if err != nil {
		return 0, err
	} else if strings.Contains(stderr, "fatal") {
		return 0, fmt.Errorf("failed to get the list of stashes: " + stderr)
	}

	return len(strings.Split(stdout, "\n")) - 1, nil
}

// GetCurrentStatus ...
func GetCurrentStatus() (GitStatus, error) {
	var branch string
	var detached, hasremote bool
	var ahead, behind, staged, conflicts, changed, untracked int
	var numStashes int

	c1 := make(chan error, 1)
	go func() {
		lines, err := GetLines("git", "status", "--porcelain", "--branch")
		if err != nil {
			c1 <- err
			return
		}
		branch, detached, hasremote, ahead, behind, err = ParseBranchLine(lines[0])
		if err != nil {
			c1 <- err
			return
		}
		staged, conflicts, changed, untracked = CollectChanges(lines[1:len(lines)])
		c1 <- nil
	}()

	c2 := make(chan error, 1)
	go func() {
		var err error
		numStashes, err = getStashCount()
		c2 <- err
	}()

	if err := <-c1; err != nil {
		return newGitStatus(), err
	}
	if err := <-c2; err != nil {
		return newGitStatus(), err
	}
	return GitStatus{branch, detached, hasremote, ahead, behind, staged, conflicts, changed, untracked, numStashes}, nil
}

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
