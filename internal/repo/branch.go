package repo

import (
	"os/exec"
	"strings"
)

// IsMerged checks if a branch is merged into any of the target branches.
// Returns true if the branch's commit is an ancestor of any target.
func (r *Repo) IsMerged(branch string, targets []string) bool {
	for _, target := range targets {
		if r.isMergedInto(branch, target) {
			return true
		}
	}
	return false
}

// isMergedInto checks if branch is merged into target.
// Uses git merge-base --is-ancestor which returns exit code 0 if true.
func (r *Repo) isMergedInto(branch, target string) bool {
	cmd := exec.Command("git", "merge-base", "--is-ancestor", branch, target)
	cmd.Dir = r.Path
	err := cmd.Run()
	return err == nil // exit code 0 means branch is ancestor of target
}

// IsStale checks if a branch tracks a remote that no longer exists.
// A branch is stale if it has an upstream that's marked as "gone".
func (r *Repo) IsStale(branch string) bool {
	cmd := exec.Command("git", "for-each-ref",
		"--format=%(upstream:track)",
		"refs/heads/"+branch)
	cmd.Dir = r.Path

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// If the upstream is gone, output contains "[gone]"
	return strings.Contains(string(output), "[gone]")
}

// HasWorktree checks if a branch has an active worktree.
func (r *Repo) HasWorktree(branch string) bool {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = r.Path

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Parse porcelain output - each worktree has:
	// worktree /path/to/worktree
	// HEAD <commit>
	// branch refs/heads/<branch>
	// (blank line)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "branch refs/heads/") {
			wtBranch := strings.TrimPrefix(line, "branch refs/heads/")
			if wtBranch == branch {
				return true
			}
		}
	}

	return false
}

// ListBranches returns all local branch names in the repository.
func (r *Repo) ListBranches() []string {
	cmd := exec.Command("git", "for-each-ref",
		"--format=%(refname:short)",
		"refs/heads/")
	cmd.Dir = r.Path

	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}
	}
	return lines
}
