package repo

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/kernelle-soft/gimme/internal/log"
)

// Repo represents a git repository with metadata.
type Repo struct {
	Repository *git.Repository
	Path       string
	Name       string
	Identifier string // Normalized remote URL (e.g. "github.com/user/repo") or path fallback
	Pinned     bool
	PinIndex   int // -1 if not pinned, otherwise the index in the pins list (lower = higher priority)
}

// NewRepo creates a new Repo with the Identifier automatically populated.
// The identifier is derived from the origin remote URL if available,
// otherwise falls back to the path.
func NewRepo(gitRepo *git.Repository, path, name string) Repo {
	return Repo{
		Repository: gitRepo,
		Path:       path,
		Name:       name,
		Identifier: identifierFromRepo(gitRepo, path),
		Pinned:     false,
		PinIndex:   -1,
	}
}

// NewPinnedRepo creates a new Repo that is pinned.
func NewPinnedRepo(gitRepo *git.Repository, path, name string, pinIndex int) Repo {
	return Repo{
		Repository: gitRepo,
		Path:       path,
		Name:       name,
		Identifier: identifierFromRepo(gitRepo, path),
		Pinned:     true,
		PinIndex:   pinIndex,
	}
}

// CurrentBranch returns the name of the current branch, or a short commit hash
// if HEAD is detached.
func (r *Repo) CurrentBranch() string {
	head, err := r.Repository.Head()
	if err != nil {
		log.Error("Error getting current branch", "path", r.Path, "err", err)
		return ""
	}

	// Head could be detached (pointing to a commit, not a branch)
	if !head.Name().IsBranch() {
		return head.Hash().String()[:7]
	}

	return head.Name().Short()
}

// identifierFromRepo returns a stable identifier for a repository.
// It tries to use the normalized origin remote URL, falling back to the path
// if no remote is configured.
func identifierFromRepo(gitRepo *git.Repository, fallbackPath string) string {
	remote, err := gitRepo.Remote("origin")
	if err != nil {
		return fallbackPath
	}

	urls := remote.Config().URLs
	if len(urls) == 0 {
		return fallbackPath
	}

	normalized := NormalizeRemoteURL(urls[0])
	if normalized == "" {
		return fallbackPath
	}

	return normalized
}

// NormalizeRemoteURL converts a git remote URL to a normalized identifier.
// Examples:
//   - git@github.com:user/repo.git     → github.com/user/repo
//   - git@github.com:user/repo         → github.com/user/repo
//   - https://github.com/user/repo.git → github.com/user/repo
//   - https://github.com/user/repo     → github.com/user/repo
//   - ssh://git@github.com/user/repo   → github.com/user/repo
//
// Returns empty string if the URL format is not recognized.
func NormalizeRemoteURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return ""
	}

	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	// Handle SSH format: git@github.com:user/repo
	if strings.HasPrefix(url, "git@") {
		// git@github.com:user/repo → github.com/user/repo
		url = strings.TrimPrefix(url, "git@")
		url = strings.Replace(url, ":", "/", 1)
		return url
	}

	// Handle ssh:// format: ssh://git@github.com/user/repo
	if strings.HasPrefix(url, "ssh://") {
		url = strings.TrimPrefix(url, "ssh://")
		url = strings.TrimPrefix(url, "git@")
		return url
	}

	// Handle https:// format: https://github.com/user/repo
	if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
		return url
	}

	// Handle http:// format: http://github.com/user/repo
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
		return url
	}

	// Unrecognized format
	return ""
}

// IdentifierFromPath returns a stable identifier for a repository given its path.
// Opens the repo to check for an origin remote.
func IdentifierFromPath(repoPath string) string {
	gitRepo, err := git.PlainOpen(repoPath)
	if err != nil {
		return repoPath
	}
	return identifierFromRepo(gitRepo, repoPath)
}
