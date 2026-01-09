package config

import (
	"strconv"

	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/spf13/cobra"
)

var deleteCommand = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete configuration values",
	Long:    `Delete configuration values such as search groups, pinned repositories, or aliases.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var deleteGroupCommand = &cobra.Command{
	Use:   "group <path|index>",
	Short: "Delete a search group",
	Long:  `Remove a search group by path or index.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if it's an index
		if idx, err := strconv.Atoi(args[0]); err == nil {
			config.DeleteGroupByIndex(idx)
		} else {
			config.DeleteGroup(args[0])
		}
	},
}

var deletePinnedRepoCommand = &cobra.Command{
	Use:   "repo <path|index>",
	Short: "Delete a pinned repository",
	Long:  `Remove a pinned repository by path or index.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if it's an index
		if idx, err := strconv.Atoi(args[0]); err == nil {
			config.DeletePinnedRepoByIndex(idx)
		} else {
			config.DeletePinnedRepo(args[0])
		}
	},
}

var deleteAliasCommand = &cobra.Command{
	Use:   "alias <short>",
	Short: "Delete an alias",
	Long:  `Remove an alias by its short name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.DeleteAlias(args[0])
	},
}

func init() {
	deleteCommand.AddCommand(deleteGroupCommand)
	deleteCommand.AddCommand(deletePinnedRepoCommand)
	deleteCommand.AddCommand(deleteAliasCommand)
}
