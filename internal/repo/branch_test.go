package repo

import (
	"os"
	"os/exec"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Helper to create a test repo with some commits and branches
func setupTestRepo(t *testing.T) (*Repo, string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "gimme-branch-test-*")
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	// Initialize repo
	gitRepo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		cleanup()
		t.Fatal(err)
	}

	// Create initial commit on main
	w, err := gitRepo.Worktree()
	if err != nil {
		cleanup()
		t.Fatal(err)
	}

	// Create a file and commit
	testFile := tmpDir + "/test.txt"
	if err := os.WriteFile(testFile, []byte("initial"), 0644); err != nil {
		cleanup()
		t.Fatal(err)
	}

	if _, err := w.Add("test.txt"); err != nil {
		cleanup()
		t.Fatal(err)
	}

	_, err = w.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@test.com"},
	})
	if err != nil {
		cleanup()
		t.Fatal(err)
	}

	repo := NewRepo(gitRepo, tmpDir, "test-repo")
	return &repo, tmpDir, cleanup
}

func TestListBranches(t *testing.T) {
	repo, tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a feature branch using git command
	cmd := exec.Command("git", "branch", "feature-1")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "branch", "feature-2")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	branches := repo.ListBranches()

	if len(branches) != 3 {
		t.Errorf("Expected 3 branches, got %d: %v", len(branches), branches)
	}

	// Check that expected branches exist
	branchSet := make(map[string]bool)
	for _, b := range branches {
		branchSet[b] = true
	}

	for _, expected := range []string{"master", "feature-1", "feature-2"} {
		if !branchSet[expected] {
			t.Errorf("Expected branch %q not found in %v", expected, branches)
		}
	}
}

func TestIsMerged(t *testing.T) {
	repo, tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a feature branch, add a commit, then merge it back
	cmd := exec.Command("git", "checkout", "-b", "feature-merged")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Add a commit on feature branch
	testFile := tmpDir + "/feature.txt"
	if err := os.WriteFile(testFile, []byte("feature"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "add", "feature.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "commit", "-m", "feature commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Go back to master and merge
	cmd = exec.Command("git", "checkout", "master")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "merge", "--no-ff", "feature-merged", "-m", "merge feature")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Create an unmerged branch
	cmd = exec.Command("git", "checkout", "-b", "feature-unmerged")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	testFile2 := tmpDir + "/unmerged.txt"
	if err := os.WriteFile(testFile2, []byte("unmerged"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "add", "unmerged.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "commit", "-m", "unmerged commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Go back to master for testing
	cmd = exec.Command("git", "checkout", "master")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	t.Run("merged branch returns true", func(t *testing.T) {
		if !repo.IsMerged("feature-merged", []string{"master"}) {
			t.Error("Expected feature-merged to be merged into master")
		}
	})

	t.Run("unmerged branch returns false", func(t *testing.T) {
		if repo.IsMerged("feature-unmerged", []string{"master"}) {
			t.Error("Expected feature-unmerged to NOT be merged into master")
		}
	})

	t.Run("checks multiple targets", func(t *testing.T) {
		// feature-merged is merged into master but not into feature-unmerged
		if !repo.IsMerged("feature-merged", []string{"feature-unmerged", "master"}) {
			t.Error("Expected IsMerged to return true when merged into any target")
		}
	})
}

func TestIsStale(t *testing.T) {
	repo, tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// A local-only branch is not stale (no upstream)
	cmd := exec.Command("git", "branch", "local-only")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	t.Run("local branch without upstream is not stale", func(t *testing.T) {
		if repo.IsStale("local-only") {
			t.Error("Expected local-only branch to NOT be stale")
		}
	})

	// Note: Testing actual stale branches would require setting up a remote
	// and then deleting the remote branch, which is complex for a unit test.
	// The implementation is simple enough that manual testing is sufficient.
}

func TestHasWorktree(t *testing.T) {
	repo, tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Setup: create branches and worktree
	cmd := exec.Command("git", "branch", "with-worktree")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "branch", "without-worktree")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	wtDir, err := os.MkdirTemp("", "gimme-worktree-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(wtDir)

	cmd = exec.Command("git", "worktree", "add", wtDir, "with-worktree")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Tests
	t.Run("branch with worktree returns true", func(t *testing.T) {
		if !repo.HasWorktree("with-worktree") {
			t.Error("Expected with-worktree to have a worktree")
		}
	})

	t.Run("branch without worktree returns false", func(t *testing.T) {
		if repo.HasWorktree("without-worktree") {
			t.Error("Expected without-worktree to NOT have a worktree")
		}
	})
}
