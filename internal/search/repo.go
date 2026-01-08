package search

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"

	"github.com/go-git/go-git/v5"
)

type Repo struct {
	Repository *git.Repository
	Path       string
	Name       string
	Pinned     bool
	PinIndex   int // -1 if not pinned, otherwise the index in the pins list (lower = higher priority)
}

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

type RepoSearchOptions struct {
	SearchFolders []string
	Query         string
}

func DefaultRepoSearchOptions() RepoSearchOptions {
	return RepoSearchOptions{
		SearchFolders: config.GetSearchFolders(),
		Query:         "",
	}
}

func ForRepo(query string) RepoSearchOptions {
	return RepoSearchOptions{
		SearchFolders: config.GetSearchFolders(),
		Query:         query,
	}
}

func InFolders(searchFolders []string) RepoSearchOptions {
	return RepoSearchOptions{
		SearchFolders: searchFolders,
		Query:         "",
	}
}

func Repositories(opts RepoSearchOptions) []Repo {
	found := []Repo{}
	pins := config.GetPins()

	for _, folder := range opts.SearchFolders {
		found = append(found, findReposRecursively(folder, opts.Query, pins)...)
	}

	// Sort alphabetically by name
	slices.SortFunc(found, func(a, b Repo) int {
		return strings.Compare(a.Name, b.Name)
	})

	return found
}

// SortByPins sorts repos with pinned repos first (by pin order), then alphabetically.
// This is useful for commands like jump that want to prioritize pinned repos.
func SortByPins(repos []Repo) {
	slices.SortFunc(repos, func(a, b Repo) int {
		if a.Pinned && !b.Pinned {
			return -1
		}
		if !a.Pinned && b.Pinned {
			return 1
		}
		// Both pinned: sort by pin index (lower index = higher priority)
		if a.Pinned && b.Pinned {
			return a.PinIndex - b.PinIndex
		}
		// Neither pinned: sort alphabetically
		return strings.Compare(a.Name, b.Name)
	})
}

func findReposRecursively(folder string, query string, pins []string) []Repo {
	results := []Repo{}

	entries, err := os.ReadDir(folder)
	if err != nil {
		log.Error("Error reading directory", "path", folder, "error", err)
		return results
	}

	for _, entry := range entries {
		path := filepath.Join(folder, entry.Name())

		if !entry.IsDir() {
			continue
		}

		repo, err := git.PlainOpen(path)
		if err != nil {
			results = append(results, findReposRecursively(path, query, pins)...)
			continue
		}

		if strings.Contains(entry.Name(), query) {
			pinIndex := slices.Index(pins, path)
			results = append(results, Repo{
				Repository: repo,
				Path:       path,
				Name:       entry.Name(),
				Pinned:     pinIndex >= 0,
				PinIndex:   pinIndex,
			})
		}
	}

	return results
}
