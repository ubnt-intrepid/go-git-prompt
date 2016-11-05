package prompt

import "testing"

func TestParseBranchLine(t *testing.T) {
	line := "## master...origin/master"
	branch, detached, hasremote, ahead, behind, err := ParseBranchLine(line)
	if err != nil {
		t.Fatal("error during parsing:", err)
	}

	if branch != "master" {
		t.Fatalf("'branch' should be equal to  %s, actual: %s", "master", branch)
	}
	if detached != false {
		t.Fatalf("'detach' should be false")
	}
	if hasremote != true {
		t.Fatalf("'hasremote' should be true")
	}
	if ahead != 0 {
		t.Fatalf("'ahead' should be 0, actual: %d", ahead)
	}
	if behind != 0 {
		t.Fatalf("'behind' should be 0, actual: %d", behind)
	}
}
