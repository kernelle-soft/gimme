package config

import (
	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add configuration values",
	Long:  `Add configuration values such as search groups, pinned repositories, or aliases.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var addGroupCommand = &cobra.Command{
	Use:   "group <path>",
	Short: "Add a search group",
	Long:  `Add a folder to search for git repositories.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.AddGroup(args[0])
	},
}

var addPinCommand = &cobra.Command{
	Use:   "pin <path>",
	Short: "Add a pinned repository",
	Long:  `Pin a repository so it appears at the top of search results.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.AddPin(args[0])
	},
}

var addAliasCommand = &cobra.Command{
	Use:   "alias <short> <expanded>",
	Short: "Add an alias",
	Long:  `Add an alias to use a short name for a repository path.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.AddAlias(args[0], args[1]); err != nil {
			log.Error("Failed to add alias: {}", err)
		}
	},
}

func init() {
	addCommand.AddCommand(addGroupCommand)
	addCommand.AddCommand(addPinCommand)
	addCommand.AddCommand(addAliasCommand)
}
