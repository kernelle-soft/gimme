package cmd

import (
	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/kernelle-soft/gimme/internal/repo"
	"github.com/kernelle-soft/gimme/internal/search"
	"github.com/spf13/cobra"
)

var listDescription = Description{
	Short: "List this workstation's visible repositories",
	Long:  "List this workstations's visible repositories. Lists repositories under each search group recursively.",
}

var listCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   listDescription.Short,
	Long:    listDescription.Long,
	Run:     listReposRun,
}

var listReposRun = func(cmd *cobra.Command, args []string) {

	var query string
	if len(args) > 0 {
		query = args[0]
	}

	// Get all repos to find pinned ones
	allRepos := search.Repositories(search.ForRepo(query))

	// Show pinned repos first as their own group
	pinnedRepos := []repo.Repo{}
	for _, r := range allRepos {
		if r.Pinned {
			pinnedRepos = append(pinnedRepos, r)
		}
	}

	if len(pinnedRepos) > 0 {
		// Sort pinned repos by pin index
		search.SortByPins(pinnedRepos)
		log.Print("pinned repositories:")
		for _, r := range pinnedRepos {
			log.Print("- {} ({})", r.Name, r.CurrentBranch())
		}
		log.Print("") // Empty line for spacing
	}

	// Show repos by search folder
	for _, folder := range config.GetSearchFolders() {
		log.Print("{}/", folder)
		repos := search.Repositories(search.RepoSearchOptions{
			Query:         query,
			SearchFolders: []string{folder},
		})

		for _, r := range repos {
			if r.Pinned {
				log.Print("📌 {} ({})", r.Name, r.CurrentBranch())
			} else {
				log.Print("-  {} ({})", r.Name, r.CurrentBranch())
			}
		}
	}
}
