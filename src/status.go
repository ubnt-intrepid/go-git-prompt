package prompt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mgutz/ansi"
)

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

// Format ...
func (s Status) Format() string {
	ret := ansi.Color("(", "yellow")

	// branch
	if s.detached {
		ret += ansi.Color("("+s.branch+")", "cyan")
	} else {
		ret += ansi.Color(s.branch, "cyan")
	}
	if s.hasremote {
		if s.ahead > 0 && s.behind > 0 {
			ret += ansi.Color(fmt.Sprintf(" A%d B%d", s.ahead, s.behind), "yellow")
		} else if s.ahead > 0 {
			ret += ansi.Color(fmt.Sprintf(" A%d", s.ahead), "green")
		} else if s.behind > 0 {
			ret += ansi.Color(fmt.Sprintf(" B%d", s.behind), "red")
		} else {
			//ret += ansi.Color(" ≡", "cyan")
		}
	}

	if s.staged > 0 || s.changed > 0 || s.conflicts > 0 || s.untracked > 0 {
		ret += fmt.Sprintf(" +%d", s.staged)
		ret += fmt.Sprintf(" ~%d", s.changed)
		ret += fmt.Sprintf(" ?%d", s.untracked)
		ret += fmt.Sprintf(" !%d", s.conflicts)
	}

	if s.stashs > 0 {
		ret += ansi.Color(fmt.Sprintf(" | s%d", s.stashs), "green")
	}

	ret += ansi.Color(")", "yellow")

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
func GetCurrentStatus() (Status, error) {
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
		return newStatus(), err
	}
	if err := <-c2; err != nil {
		return newStatus(), err
	}
	return Status{branch, detached, hasremote, ahead, behind, staged, conflicts, changed, untracked, numStashes}, nil
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
