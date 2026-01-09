package cmd

import (
	"os"

	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/spf13/cobra"
)

var pinBranchFlag bool

var pinCommand = &cobra.Command{
	Use:   "pin [repo|branch]",
	Short: "Pin a repository or branch",
	Long: `Pin a repository or branch.

Without -b flag: pins a repository so it appears at the top of search results.
  gimme pin           - pin current directory's repo
  gimme pin <path>    - pin repo at path

With -b flag: pins a branch in the current repo (protects from clean).
  gimme pin -b        - pin current branch
  gimme pin -b <name> - pin branch by name`,
	Run: pinRun,
}

func init() {
	pinCommand.Flags().BoolVarP(&pinBranchFlag, "branch", "b", false, "Pin a branch instead of a repository")
}

var pinRun = func(cmd *cobra.Command, args []string) {
	if pinBranchFlag {
		pinBranch(args)
	} else {
		pinRepo(args)
	}
}

func pinRepo(args []string) {
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

	config.AddPinnedRepo(repoPath)
}

func pinBranch(args []string) {
	// Get current repo
	cwd, err := os.Getwd()
	if err != nil {
		log.Error("Could not determine current directory: {}", err)
		return
	}

	currentRepo := findRepoForPath(cwd)
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

	// Check if branch exists
	branches := currentRepo.ListBranches()
	if !isInSlice(branchName, branches) {
		log.Print("Branch \"{}\" not found.", branchName)
		return
	}

	// Check if already globally pinned
	if config.IsBranchGloballyPinned(branchName) {
		log.Print("Branch \"{}\" is already globally protected.", branchName)
		return
	}

	config.AddRepoPinnedBranch(currentRepo.Identifier, branchName)
}
