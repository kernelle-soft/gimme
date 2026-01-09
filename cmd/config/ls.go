package config

import (
	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/spf13/cobra"
)

var lsCommand = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List configuration values",
	Long:    `List configuration values. Use subcommands to show specific sections, or run without subcommands to show all.`,
	Run:     lsRun,
}

var lsGroupCommand = &cobra.Command{
	Use:     "group",
	Aliases: []string{"groups"},
	Short:   "List search groups",
	Long:    `List all configured search groups.`,
	Run: func(cmd *cobra.Command, args []string) {
		showGroups()
	},
}

var lsPinnedRepoCommand = &cobra.Command{
	Use:     "repo",
	Aliases: []string{"repos"},
	Short:   "List pinned repositories",
	Long:    `List all pinned repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		showPinnedRepos()
	},
}

var lsPinnedBranchCommand = &cobra.Command{
	Use:     "branch",
	Aliases: []string{"branches"},
	Short:   "List pinned branches",
	Long:    `List all pinned branches (global and per-repo).`,
	Run: func(cmd *cobra.Command, args []string) {
		showPinnedBranches()
	},
}

var lsAliasCommand = &cobra.Command{
	Use:     "alias",
	Aliases: []string{"aliases"},
	Short:   "List aliases",
	Long:    `List all configured aliases.`,
	Run: func(cmd *cobra.Command, args []string) {
		showAliases()
	},
}

func init() {
	lsCommand.AddCommand(lsGroupCommand)
	lsCommand.AddCommand(lsPinnedRepoCommand)
	lsCommand.AddCommand(lsPinnedBranchCommand)
	lsCommand.AddCommand(lsAliasCommand)
}

var lsRun = func(cmd *cobra.Command, args []string) {
	// Show all config
	showGroups()
	log.Print("")
	showPinnedRepos()
	log.Print("")
	showPinnedBranches()
	log.Print("")
	showAliases()
}

func showGroups() {
	groups := config.GetSearchFolders()
	log.Print("Search Groups:")
	if len(groups) == 0 {
		log.Print("  (none configured)")
		return
	}
	for i, group := range groups {
		log.Print("  [{}] {}", i, group)
	}
}

func showPinnedRepos() {
	repos := config.GetPinnedRepos()
	log.Print("Pinned Repositories:")
	if len(repos) == 0 {
		log.Print("  (none configured)")
		return
	}
	for i, repo := range repos {
		log.Print("  [{}] {}", i, repo)
	}
}

func showPinnedBranches() {
	branches := config.GetPinnedBranches()
	log.Print("Pinned Branches (global):")
	if len(branches) == 0 {
		log.Print("  (none configured)")
	} else {
		for _, branch := range branches {
			log.Print("  🛡️  {}", branch)
		}
	}

	repoBranches := config.GetRepoPinnedBranches()
	if len(repoBranches) > 0 {
		log.Print("")
		log.Print("Pinned Branches (per-repo):")
		for repo, branches := range repoBranches {
			log.Print("  {}:", repo)
			for _, branch := range branches {
				log.Print("    📌 {}", branch)
			}
		}
	}
}

func showAliases() {
	aliases := config.GetAliases()
	log.Print("Aliases:")
	if len(aliases) == 0 {
		log.Print("  (none configured)")
		return
	}
	for short, expanded := range aliases {
		log.Print("  {} -> {}", short, expanded)
	}
}
