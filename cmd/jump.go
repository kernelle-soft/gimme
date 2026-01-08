package cmd

import (
	"os"

	"github.com/kernelle-soft/gimme/internal/config"
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/kernelle-soft/gimme/internal/path"
	"github.com/kernelle-soft/gimme/internal/search"
	"github.com/spf13/cobra"
)

var jumpToRepoCommand = &cobra.Command{
	Use:   "jump <repo>",
	Short: "Jump to a project's root directory",
	Long: `Jump to a project's root directory. 
	
	This command will change the current working directory to the root directory of the project.
	The specified string will be used to search recursively for the project's directory name and hit the first match, including partial matches.
	
	Example:
	gimme jump kernelle # jumps to the kernelle project's root directory
	gimme kernelle # same as above. 'jump' is optional and simply intended for disambiguation if ever necessary.
    gimme kern # jumps to the kernelle project's root directory because 'kern' is a partial match for 'kernelle'.
	`,
	Run: jumpRun,
}

var jumpRun = func(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	query := args[0]

	// Expand alias if one exists
	aliases := config.GetAliases()
	if expanded, ok := aliases[query]; ok {
		query = expanded
	}

	// Check if query is a direct path (alias may have expanded to a path)
	normalizedQuery, _ := path.Normalize(query)
	if info, err := os.Stat(normalizedQuery); err == nil && info.IsDir() {
		log.ToStdout(normalizedQuery)
		return
	}

	// Search for repositories matching the query
	found := search.Repositories(search.ForRepo(query))
	if len(found) == 0 {
		log.Print("No repositories found for \"{}\".", query)
		return
	}

	// Sort by pins so pinned repos are prioritized
	search.SortByPins(found)

	// Found match.
	log.ToStdout(found[0].Path)
}
