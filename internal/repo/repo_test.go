package repo

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func TestNormalizeRemoteURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// GitHub SSH format
		{
			name:     "github ssh with .git",
			input:    "git@github.com:user/repo.git",
			expected: "github.com/user/repo",
		},
		{
			name:     "github ssh without .git",
			input:    "git@github.com:user/repo",
			expected: "github.com/user/repo",
		},

		// GitHub HTTPS format
		{
			name:     "github https with .git",
			input:    "https://github.com/user/repo.git",
			expected: "github.com/user/repo",
		},
		{
			name:     "github https without .git",
			input:    "https://github.com/user/repo",
			expected: "github.com/user/repo",
		},

		// GitHub ssh:// format
		{
			name:     "github ssh:// with .git",
			input:    "ssh://git@github.com/user/repo.git",
			expected: "github.com/user/repo",
		},
		{
			name:     "github ssh:// without .git",
			input:    "ssh://git@github.com/user/repo",
			expected: "github.com/user/repo",
		},

		// Edge cases
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "",
		},
		{
			name:     "input with whitespace",
			input:    "  git@github.com:user/repo.git  ",
			expected: "github.com/user/repo",
		},
		{
			name:     "unrecognized format",
			input:    "not-a-url",
			expected: "",
		},

		// HTTP (less common but valid)
		{
			name:     "http format",
			input:    "http://github.com/user/repo.git",
			expected: "github.com/user/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeRemoteURL(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeRemoteURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewRepo(t *testing.T) {
	t.Run("with origin remote", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "gimme-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		gitRepo, err := git.PlainInit(tmpDir, false)
		if err != nil {
			t.Fatal(err)
		}

		_, err = gitRepo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{"git@github.com:testuser/testrepo.git"},
		})
		if err != nil {
			t.Fatal(err)
		}

		repo := NewRepo(gitRepo, tmpDir, "testrepo")

		if repo.Identifier != "github.com/testuser/testrepo" {
			t.Errorf("Identifier = %q, want %q", repo.Identifier, "github.com/testuser/testrepo")
		}
		if repo.Path != tmpDir {
			t.Errorf("Path = %q, want %q", repo.Path, tmpDir)
		}
		if repo.Name != "testrepo" {
			t.Errorf("Name = %q, want %q", repo.Name, "testrepo")
		}
		if repo.Pinned {
			t.Error("Pinned should be false")
		}
		if repo.PinIndex != -1 {
			t.Errorf("PinIndex = %d, want -1", repo.PinIndex)
		}
	})

	t.Run("without origin remote falls back to path", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "gimme-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		gitRepo, err := git.PlainInit(tmpDir, false)
		if err != nil {
			t.Fatal(err)
		}

		repo := NewRepo(gitRepo, tmpDir, "localrepo")

		if repo.Identifier != tmpDir {
			t.Errorf("Identifier = %q, want %q (path fallback)", repo.Identifier, tmpDir)
		}
	})
}

func TestNewPinnedRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gimme-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	gitRepo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = gitRepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/myorg/myrepo.git"},
	})
	if err != nil {
		t.Fatal(err)
	}

	repo := NewPinnedRepo(gitRepo, tmpDir, "myrepo", 2)

	if repo.Identifier != "github.com/myorg/myrepo" {
		t.Errorf("Identifier = %q, want %q", repo.Identifier, "github.com/myorg/myrepo")
	}
	if !repo.Pinned {
		t.Error("Pinned should be true")
	}
	if repo.PinIndex != 2 {
		t.Errorf("PinIndex = %d, want 2", repo.PinIndex)
	}
}

func TestIdentifierFromPath(t *testing.T) {
	t.Run("repo with origin", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "gimme-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		gitRepo, err := git.PlainInit(tmpDir, false)
		if err != nil {
			t.Fatal(err)
		}

		_, err = gitRepo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{"git@github.com:org/project.git"},
		})
		if err != nil {
			t.Fatal(err)
		}

		result := IdentifierFromPath(tmpDir)
		if result != "github.com/org/project" {
			t.Errorf("IdentifierFromPath() = %q, want %q", result, "github.com/org/project")
		}
	})

	t.Run("non-git directory falls back to path", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "gimme-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		result := IdentifierFromPath(tmpDir)
		if result != tmpDir {
			t.Errorf("IdentifierFromPath() = %q, want %q", result, tmpDir)
		}
	})
}
