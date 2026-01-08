package search

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"

	"github.com/go-git/go-git/v5"
)

type Repo struct {
	Repository *git.Repository
	Path       string
	Name       string
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

	for _, folder := range opts.SearchFolders {
		found = append(found, findReposRecursively(folder, opts.Query)...)
	}

	return found
}

func findReposRecursively(folder string, query string) []Repo {
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
			results = append(results, findReposRecursively(path, query)...)
			continue
		}

		if strings.Contains(entry.Name(), query) {
			results = append(results, Repo{
				Repository: repo,
				Path:       path,
				Name:       entry.Name(),
			})
		}
	}

	return results
}
