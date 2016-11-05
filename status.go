package main

import (
	"fmt"
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
	ret := ansi.Color("[", "yellow")

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
		ret += ansi.Color(fmt.Sprintf(" +%d", s.staged), "red")
		ret += ansi.Color(fmt.Sprintf(" ~%d", s.changed), "red")
		ret += ansi.Color(fmt.Sprintf(" ?%d", s.untracked), "red")
		ret += ansi.Color(fmt.Sprintf(" !%d", s.conflicts), "red")
	}

	if s.stashs > 0 {
		ret += ansi.Color(fmt.Sprintf(" | s%d", s.stashs), "green")
	}

	ret += ansi.Color("]", "yellow")

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
	lines, err := GetLines("git", "status", "--porcelain", "--branch")
	if err != nil {
		return newStatus(), err
	}

	numStashes, err := getStashCount()
	if err != nil {
		return newStatus(), err
	}

	branch, detached, hasremote, ahead, behind, err := ParseBranchLine(lines[0])
	if err != nil {
		return newStatus(), err
	}
	staged, conflicts, changed, untracked := CollectChanges(lines[1:len(lines)])

	return Status{branch, detached, hasremote, ahead, behind, staged, conflicts, changed, untracked, numStashes}, nil
}
