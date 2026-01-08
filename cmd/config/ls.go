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

var lsPinCommand = &cobra.Command{
	Use:     "pin",
	Aliases: []string{"pins"},
	Short:   "List pinned repositories",
	Long:    `List all pinned repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		showPins()
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
	lsCommand.AddCommand(lsPinCommand)
	lsCommand.AddCommand(lsAliasCommand)
}

var lsRun = func(cmd *cobra.Command, args []string) {
	// Show all config
	showGroups()
	log.Print("")
	showPins()
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

func showPins() {
	pins := config.GetPins()
	log.Print("Pinned Repositories:")
	if len(pins) == 0 {
		log.Print("  (none configured)")
		return
	}
	for i, pin := range pins {
		log.Print("  [{}] {}", i, pin)
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
