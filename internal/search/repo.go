package search

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/kernelle-soft/gimme/internal/repo"

	"github.com/go-git/go-git/v5"
)

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

func Repositories(opts RepoSearchOptions) []repo.Repo {
	found := []repo.Repo{}
	pins := config.GetPinnedRepos()

	for _, folder := range opts.SearchFolders {
		found = append(found, findReposRecursively(folder, opts.Query, pins)...)
	}

	// Sort alphabetically by name
	slices.SortFunc(found, func(a, b repo.Repo) int {
		return strings.Compare(a.Name, b.Name)
	})

	return found
}

// SortByPins sorts repos with pinned repos first (by pin order), then alphabetically.
// This is useful for commands like jump that want to prioritize pinned repos.
func SortByPins(repos []repo.Repo) {
	slices.SortFunc(repos, func(a, b repo.Repo) int {
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

func findReposRecursively(folder string, query string, pins []string) []repo.Repo {
	results := []repo.Repo{}

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

		gitRepo, err := git.PlainOpen(path)
		if err != nil {
			results = append(results, findReposRecursively(path, query, pins)...)
			continue
		}

		if strings.Contains(entry.Name(), query) {
			pinIndex := slices.Index(pins, path)
			if pinIndex >= 0 {
				results = append(results, repo.NewPinnedRepo(gitRepo, path, entry.Name(), pinIndex))
			} else {
				results = append(results, repo.NewRepo(gitRepo, path, entry.Name()))
			}
		}
	}

	return results
}

// FindRepoForPath finds the repository that contains the given path
func FindRepoForPath(path string) *repo.Repo {
	repos := Repositories(DefaultRepoSearchOptions())
	for _, r := range repos {
		if strings.HasPrefix(path, r.Path) {
			return &r
		}
	}
	return nil
}
