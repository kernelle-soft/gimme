package cmd

import (
	"fmt"

	"github.com/kernelle-soft/gimmetool/internal/config"
	"github.com/kernelle-soft/gimmetool/internal/search"
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

	for _, folder := range config.GetSearchFolders() {
		fmt.Printf("%s\n", folder)
		repos := search.Repositories(search.RepoSearchOptions{
			Query:         query,
			SearchFolders: []string{folder},
		})

		for _, repo := range repos {
			fmt.Printf("- %s (%s)\n", repo.Name, repo.CurrentBranch())
		}
	}
}
