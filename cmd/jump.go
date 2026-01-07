package cmd

import (
	"fmt"
	"os"

	"github.com/kernelle-soft/gimmetool/internal/search"
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
	found := search.Repositories(search.ForRepo(query))
	if len(found) == 0 {
		fmt.Fprintf(os.Stderr, "No repositories or aliases found beginning with '%s'", query)
		return
	}

	fmt.Printf("First match: %s", found[0])
}
