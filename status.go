package main

import "fmt"
import "github.com/fatih/color"

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
	ret := "["

	// branch
	if s.detached {
		ret += color.CyanString("(" + s.branch + ")")
	} else {
		ret += color.CyanString(s.branch)
	}
	if s.hasremote {
		if s.ahead > 0 && s.behind > 0 {
			ret += color.YellowString(" ↑%d ↓%d", s.ahead, s.behind)
		} else if s.ahead > 0 {
			ret += color.GreenString(" ↑%d", s.ahead)
		} else if s.behind > 0 {
			ret += color.RedString(" ↓%d", s.behind)
		} else {
			ret += color.CyanString(" ≡")
		}
	}

	if s.staged > 0 || s.changed > 0 || s.conflicts > 0 || s.untracked > 0 {
		ret += color.RedString(" +%d", s.staged)
		ret += color.RedString(" ~%d", s.changed)
		ret += color.RedString(" ?%d", s.untracked)
		ret += color.RedString(" !%d", s.conflicts)
	}

	if s.stashs > 0 {
		ret += color.GreenString(" | s%d", s.stashs)
	}

	ret += "]"

	return ret
}

// GetCurrentStatus ...
func GetCurrentStatus() (Status, error) {
	lines, err := GetLines("git", "status", "--porcelain", "--branch")
	if err != nil {
		return newStatus(), err
	}

	stashes, err := GetLines("git", "stash", "list")
	if err != nil {
		return newStatus(), err
	}

	branch, detached, hasremote, ahead, behind, err := ParseBranchLine(lines[0])
	if err != nil {
		return newStatus(), err
	}
	staged, conflicts, changed, untracked := CollectChanges(lines[1:len(lines)])

	return Status{branch, detached, hasremote, ahead, behind, staged, conflicts, changed, untracked, len(stashes)}, nil
}
