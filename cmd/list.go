package cmd

import (
	"os"
	"strings"

	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/kernelle-soft/gimme/internal/repo"
	"github.com/kernelle-soft/gimme/internal/search"
	"github.com/kernelle-soft/gimme/internal/slice"
	"github.com/spf13/cobra"
)

var (
	listBranchFlag   bool
	listMergedFlag   bool
	listNoMergedFlag bool
)

var listDescription = Description{
	Short: "List this workstation's visible repositories",
	Long:  "List this workstations's visible repositories. Lists repositories under each search group recursively. Use -b to list branches instead.",
}

var listCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   listDescription.Short,
	Long:    listDescription.Long,
	Run:     listRun,
}

func init() {
	listCommand.Flags().BoolVarP(&listBranchFlag, "branch", "b", false, "List branches instead of repositories")
	listCommand.Flags().BoolVar(&listMergedFlag, "merged", false, "Show only merged branches (requires -b)")
	listCommand.Flags().BoolVar(&listNoMergedFlag, "no-merged", false, "Show only unmerged branches (requires -b)")
}

var listRun = func(cmd *cobra.Command, args []string) {
	if listBranchFlag {
		listBranches()
	} else {
		listRepos(args)
	}
}

func listRepos(args []string) {
	var query string
	if len(args) > 0 {
		query = args[0]
	}

	// Get all repos to find pinned ones
	allRepos := search.Repositories(search.ForRepo(query))

	// Show pinned repos first as their own group
	pinnedRepos := []repo.Repo{}
	for _, r := range allRepos {
		if r.Pinned {
			pinnedRepos = append(pinnedRepos, r)
		}
	}

	if len(pinnedRepos) > 0 {
		// Sort pinned repos by pin index
		search.SortByPins(pinnedRepos)
		log.Print("pinned repositories:")
		for _, r := range pinnedRepos {
			log.Print("- {} ({})", r.Name, r.CurrentBranch())
		}
		log.Print("") // Empty line for spacing
	}

	// Show repos by search folder
	for _, folder := range config.GetSearchFolders() {
		log.Print("{}/", folder)
		repos := search.Repositories(search.RepoSearchOptions{
			Query:         query,
			SearchFolders: []string{folder},
		})

		for _, r := range repos {
			if r.Pinned {
				log.Print("  {} ({}) [pinned]", r.Name, r.CurrentBranch())
			} else {
				log.Print("  {} ({})", r.Name, r.CurrentBranch())
			}
		}
	}
}

func listBranches() {
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

	globalPins := config.GetGlobalPinnedBranches()
	repoPins := config.GetRepoPinnedBranches()
	repoPinnedBranches := repoPins[currentRepo.Identifier]

	branches := currentRepo.ListBranches()
	currentBranch := currentRepo.CurrentBranch()

	log.Print("{}/", currentRepo.Name)

	for _, branch := range branches {
		// Check merged/unmerged filter
		isMerged := currentRepo.IsMerged(branch, globalPins)

		if listMergedFlag && !isMerged {
			continue
		}
		if listNoMergedFlag && isMerged {
			continue
		}

		// Current branch marker
		prefix := "  "
		if branch == currentBranch {
			prefix = "* "
		}

		// Pin status in square brackets
		pinStatus := ""
		if slice.Contains(globalPins, branch) {
			pinStatus = " [protected]"
		} else if slice.Contains(repoPinnedBranches, branch) {
			pinStatus = " [pinned]"
		}

		// Other status indicators in parentheses
		statusIndicators := []string{}
		if isMerged && !slice.Contains(globalPins, branch) {
			statusIndicators = append(statusIndicators, "merged")
		}

		if currentRepo.IsStale(branch) {
			statusIndicators = append(statusIndicators, "stale")
		}

		if currentRepo.HasWorktree(branch) && branch != currentBranch {
			statusIndicators = append(statusIndicators, "worktree")
		}

		// Format: "* branch (merged, stale) [protected]"
		statusPart := ""
		if len(statusIndicators) > 0 {
			statusPart = " (" + strings.Join(statusIndicators, ", ") + ")"
		}
		log.Print("{}{}{}{}", prefix, branch, statusPart, pinStatus)
	}
}
