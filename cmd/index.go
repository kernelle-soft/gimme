package cmd

import (
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
	"github.com/spf13/cobra"

	"github.com/kernelle-soft/gimme/internal/config"
)

type Description struct {
	Short string
	Long  string
}

var root = &cobra.Command{
	Use:   "gimme <repo>",
	Short: "The multi-repo manager for professional developers",
	Long:  indent.String(wordwrap.String(`The multi-repo manager for professional developers. Gimme is a tool that helps you streamline the process of hopping from project to project, branch to branch, and worktree to worktree.`, 80), 2),
	Run:   jumpRun,
	Args:  cobra.MaximumNArgs(1),
}

func init() {
	config.Load()

	root.AddCommand(jumpToRepoCommand)
	root.AddCommand(listCommand)
}

func Execute() {
	root.Execute()
}
