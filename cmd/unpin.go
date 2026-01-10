package cmd

import (
	"os"

	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/kernelle-soft/gimme/internal/search"
	"github.com/spf13/cobra"
)

var unpinBranchFlag bool

var unpinCommand = &cobra.Command{
	Use:   "unpin [repo|branch]",
	Short: "Unpin a repository or branch",
	Long: `Unpin a repository or branch. This will deprioritize the repository in searching and jumping, and will remove protection on the branch.

Without -b flag: unpins a repository.
  gimme unpin           - unpin current directory's repo
  gimme unpin <path>    - unpin repo at path

With -b flag: unpins a branch in the current repo.
  gimme unpin -b        - unpin current branch
  gimme unpin -b <name> - unpin branch by name

Note: Cannot unpin globally protected branches (main, master, etc.) - use config to modify those.`,
	Run: unpinRun,
}

func init() {
	unpinCommand.Flags().BoolVarP(&unpinBranchFlag, "branch", "b", false, "Unpin a branch instead of a repository")
}

var unpinRun = func(cmd *cobra.Command, args []string) {
	if unpinBranchFlag {
		unpinBranch(args)
	} else {
		unpinRepo(args)
	}
}

func unpinRepo(args []string) {
	var repoPath string

	if len(args) > 0 {
		repoPath = args[0]
	} else {
		// Use current directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Error("Could not determine current directory: {}", err)
			return
		}
		repoPath = cwd
	}

	config.DeletePinnedRepo(repoPath)
}

func unpinBranch(args []string) {
	// Get current repo
	cwd, err := os.Getwd()
	if err != nil {
		log.Error("Could not determine current directory: {}", err)
		return
	}

	currentRepo := search.FindRepoForPath(cwd)
	if currentRepo == nil {
		log.Print("Not in a git repository.")
		return
	}

	var branchName string
	if len(args) > 0 {
		branchName = args[0]
	} else {
		// Use current branch
		branchName = currentRepo.CurrentBranch()
	}

	// Check if globally pinned - can't unpin those
	if config.IsBranchGloballyPinned(branchName) {
		log.Print("Branch \"{}\" is globally protected. Use 'gimme config' to modify global settings.", branchName)
		return
	}

	// Check if branch is pinned for this repo
	if !config.IsBranchPinnedForRepo(currentRepo.Identifier, branchName) {
		log.Print("Branch \"{}\" is not pinned for this repo.", branchName)
		return
	}

	config.DeleteRepoPinnedBranch(currentRepo.Identifier, branchName)
}
