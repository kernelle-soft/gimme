package cmd

import (
	"os"

	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/kernelle-soft/gimme/internal/search"
	"github.com/kernelle-soft/gimme/internal/slice"
	"github.com/spf13/cobra"
)

var (
	cleanBranchFlag  bool
	cleanAllFlag     bool
	cleanDryRunFlag  bool
	cleanForceFlag   bool
	cleanVerboseFlag bool
)

var cleanCommand = &cobra.Command{
	Use:   "clean",
	Short: "Clean up branches",
	Long: `Clean up branches in the current repository.

Requires -b flag for branch cleaning (future: may support other clean operations).

With -b flag: bulk delete branches with protection awareness.
  gimme clean -b            - delete merged branches (default)
  gimme clean -b --all      - delete all non-pinned branches
  gimme clean -b --dry-run  - preview without deleting

Protection hierarchy (branches that won't be deleted):
  1. Current branch — always protected
  2. Global pins (main, master, etc.) — always protected (even with --force)
  3. Per-repo pinned branches — protected unless --force
  4. Branches with active worktrees — always skipped`,
	Run: cleanRun,
}

func init() {
	cleanCommand.Flags().BoolVarP(&cleanBranchFlag, "branch", "b", false, "Clean branches (required)")
	cleanCommand.Flags().BoolVar(&cleanAllFlag, "all", false, "Delete all non-pinned branches (default: merged only)")
	cleanCommand.Flags().BoolVar(&cleanDryRunFlag, "dry-run", false, "Preview without deleting")
	cleanCommand.Flags().BoolVar(&cleanForceFlag, "force", false, "Include per-repo pinned branches")
	cleanCommand.Flags().BoolVarP(&cleanVerboseFlag, "verbose", "v", false, "Show each deleted branch")
}

var cleanRun = func(cmd *cobra.Command, args []string) {
	if !cleanBranchFlag {
		log.Print("Please specify -b flag to clean branches.")
		log.Print("Usage: gimme clean -b [--all] [--dry-run] [--force] [-v]")
		return
	}

	cleanBranches()
}

func cleanBranches() {
	// Get current working directory to determine which repo we're in
	cwd, err := os.Getwd()
	if err != nil {
		log.Error("Could not determine current directory: {}", err)
		return
	}

	// Find the repo for the current directory
	currentRepo := search.FindRepoForPath(cwd)
	if currentRepo == nil {
		log.Print("Not in a git repository.")
		return
	}

	// Get protection lists
	globalPins := config.GetGlobalPinnedBranches()
	repoPins := config.GetRepoPinnedBranches()
	repoPinnedBranches := repoPins[currentRepo.Identifier]

	// Get all branches and current branch
	branches := currentRepo.ListBranches()
	currentBranch := currentRepo.CurrentBranch()

	// Determine which branches to delete
	var toDelete []string
	var skipped []string

	for _, branch := range branches {
		// Protection check 1: Current branch — always protected
		if branch == currentBranch {
			continue
		}

		// Protection check 2: Global pins — always protected (even with --force)
		if slice.Contains(globalPins, branch) {
			continue
		}

		// Protection check 3: Per-repo pins — protected unless --force
		if slice.Contains(repoPinnedBranches, branch) && !cleanForceFlag {
			continue
		}

		// Protection check 4: Branches with active worktrees — always skipped
		if currentRepo.HasWorktree(branch) {
			skipped = append(skipped, branch)
			continue
		}

		// Apply filter: default is merged-only, --all skips this check
		if !cleanAllFlag {
			// Only delete if merged into any global pin
			if !currentRepo.IsMerged(branch, globalPins) {
				continue
			}
		}

		toDelete = append(toDelete, branch)
	}

	// Handle dry-run
	if cleanDryRunFlag {
		if len(toDelete) == 0 {
			log.Print("No branches to delete.")
		} else {
			log.Print("Would delete {} branches:", len(toDelete))
			for _, branch := range toDelete {
				log.Print("  {}", branch)
			}
		}
		if len(skipped) > 0 {
			log.Print("")
			log.Print("Skipped {} branches with active worktrees:", len(skipped))
			for _, branch := range skipped {
				log.Print("  {}", branch)
			}
		}
		return
	}

	// Delete branches
	deletedCount := 0
	for _, branch := range toDelete {
		err := currentRepo.DeleteBranch(branch)
		if err != nil {
			log.Warning("Failed to delete branch \"{}\": {}", branch, err)
			continue
		}
		deletedCount++
		if cleanVerboseFlag {
			log.Print("Deleted branch \"{}\".", branch)
		}
	}

	// Output summary
	switch deletedCount {
		case 0:
			log.Print("No branches deleted.")
		case 1:
			log.Print("Deleted 1 branch.")
		default:
			log.Print("Deleted {} branches.", deletedCount)
	}

	// Show skipped worktrees if any
	if len(skipped) > 0 && cleanVerboseFlag {
		log.Print("")
		log.Print("Skipped {} branches with active worktrees:", len(skipped))
		for _, branch := range skipped {
			log.Print("  {}", branch)
		}
	}
}
