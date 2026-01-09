package config

import (
	"strconv"

	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/spf13/cobra"
)

var deleteCommand = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete configuration values",
	Long:    `Delete configuration values such as search groups or aliases.`,
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

var deleteAliasCommand = &cobra.Command{
	Use:   "alias <short>",
	Short: "Delete an alias",
	Long:  `Remove an alias by its short name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.DeleteAlias(args[0])
	},
}

var deleteProtectedCommand = &cobra.Command{
	Use:   "protected <branch>",
	Short: "Delete a protected branch",
	Long:  `Remove a branch from the global protected branches list.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.DeleteGlobalPinnedBranch(args[0]); err != nil {
			log.Error("Failed to delete protected branch: {}", err)
		}
	},
}

func init() {
	deleteCommand.AddCommand(deleteGroupCommand)
	deleteCommand.AddCommand(deleteAliasCommand)
	deleteCommand.AddCommand(deleteProtectedCommand)
}
