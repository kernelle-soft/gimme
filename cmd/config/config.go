package config

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "config",
	Short: "Manage gimme configuration",
	Long:  `Manages gimme's configuration. This command allows you to set and display configuration values such as search groups, pinned repositories, and aliases.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	Command.AddCommand(addCommand)
	Command.AddCommand(deleteCommand)
	Command.AddCommand(lsCommand)
}
