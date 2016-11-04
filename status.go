package main

import "fmt"

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
	lines, err := GetLines("git", "status", "--porcelain", "--branch")
	if err != nil {
		return newStatus(), err
	}

	stashs, err := GetLines("git", "stash", "list")
	if err != nil {
		return newStatus(), err
	}

	branch, detached, hasremote, ahead, behind, err := ParseBranchLine(lines[0])
	if err != nil {
		return newStatus(), err
	}
	staged, conflicts, changed, untracked := CollectChanges(lines[1:len(lines)])

	return Status{branch, detached, hasremote, ahead, behind, staged, conflicts, changed, untracked, len(stashs)}, nil
}
