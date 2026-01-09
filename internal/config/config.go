package config

import (
	"os"
	"path/filepath"

	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/kernelle-soft/gimme/internal/path"
	"github.com/kernelle-soft/gimme/internal/slice"
	"github.com/spf13/viper"
)

const (
	keySearchFolders      = "search-folders"
	keyPinnedRepos        = "pinned-repos"
	keyPinnedBranches     = "pinned-branches"
	keyRepoPinnedBranches = "repo-pinned-branches"
	keyAliases            = "aliases"
)

func Load() {
	// Set defaults
	viper.SetDefault(keySearchFolders, []string{"~/"})
	viper.SetDefault(keyPinnedRepos, []string{})
	viper.SetDefault(keyPinnedBranches, []string{"main", "master"})
	viper.SetDefault(keyRepoPinnedBranches, map[string][]string{})
	viper.SetDefault(keyAliases, map[string]string{})

	// Config file location
	viper.SetConfigName(".gimme.config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Error("Error reading gimme configuration", "err", err)
		}
	}
}

// GetSearchFolders returns the list of search folders (groups)
func GetSearchFolders() []string {
	rawPaths := viper.GetStringSlice(keySearchFolders)

	return slice.Map(rawPaths, func(rawPath string) string {
		normalized, err := path.Normalize(rawPath)
		if err != nil {
			log.Error("Error parsing search folder \"{}\". Error: {}", rawPath, err)
		}
		return normalized
	})
}

// AddGroup adds a search group path
func AddGroup(groupPath string) error {
	groups := viper.GetStringSlice(keySearchFolders)

	// Normalize the path
	normalized, err := path.Normalize(groupPath)
	if err != nil {
		log.Error("Error parsing search folder \"{}\". Error: {}", groupPath, err)
	}

	// Check if already exists
	for _, g := range groups {
		existingNorm, _ := path.Normalize(g)
		if existingNorm == normalized {
			log.Error("Group already exists: \"{}\".", groupPath)
			return nil
		}
	}

	groups = append(groups, groupPath)
	viper.Set(keySearchFolders, groups)
	log.Print("Added search group \"{}\".", groupPath)
	return saveConfig()
}

// DeleteGroup removes a search group by path
func DeleteGroup(groupPath string) error {
	groups := viper.GetStringSlice(keySearchFolders)
	normalized, _ := path.Normalize(groupPath)

	newGroups := []string{}
	found := false
	for _, g := range groups {
		existingNorm, _ := path.Normalize(g)
		if existingNorm == normalized {
			found = true
			continue
		}
		newGroups = append(newGroups, g)
	}

	if !found {
		log.Error("Group not found: \"{}\".", groupPath)
		return nil
	}

	viper.Set(keySearchFolders, newGroups)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Deleted search group \"{}\".", groupPath)
	return nil
}

// DeleteGroupByIndex removes a search group by index
func DeleteGroupByIndex(index int) error {
	groups := viper.GetStringSlice(keySearchFolders)

	if index < 0 || index >= len(groups) {
		log.Error("Index out of range: {} (have {} groups).", index, len(groups))
		return nil
	}

	groupPath := groups[index]
	groups = append(groups[:index], groups[index+1:]...)
	viper.Set(keySearchFolders, groups)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}

	log.Print("Deleted search group \"{}\".", groupPath)
	return nil
}

// GetPinnedRepos returns the list of pinned repositories
func GetPinnedRepos() []string {
	rawPaths := viper.GetStringSlice(keyPinnedRepos)

	return slice.Map(rawPaths, func(rawPath string) string {
		normalized, err := path.Normalize(rawPath)
		if err != nil {
			log.Error("Error parsing pinned path \"{}\". Error: {}", rawPath, err)
			return rawPath
		}
		return normalized
	})
}

// AddPinnedRepo adds a pinned repository path
func AddPinnedRepo(repoPath string) error {
	repos := viper.GetStringSlice(keyPinnedRepos)

	// Normalize the path
	normalized, err := path.Normalize(repoPath)
	if err != nil {
		log.Error("Error parsing pinned path \"{}\". Error: {}", repoPath, err)
		return nil
	}

	// Check if already exists
	for _, r := range repos {
		existingNorm, _ := path.Normalize(r)
		if existingNorm == normalized {
			log.Error("Pinned repo already exists: \"{}\".", repoPath)
			return nil
		}
	}

	repos = append(repos, repoPath)
	viper.Set(keyPinnedRepos, repos)
	err = saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Added pinned repository \"{}\".", repoPath)
	return nil
}

// DeletePinnedRepo removes a pinned repository by path
func DeletePinnedRepo(repoPath string) error {
	repos := viper.GetStringSlice(keyPinnedRepos)
	normalized, _ := path.Normalize(repoPath)

	newRepos := []string{}
	found := false
	for _, r := range repos {
		existingNorm, _ := path.Normalize(r)
		if existingNorm == normalized {
			found = true
			continue
		}
		newRepos = append(newRepos, r)
	}

	if !found {
		log.Error("Pinned repo not found: \"{}\".", repoPath)
		return nil
	}

	viper.Set(keyPinnedRepos, newRepos)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Deleted pinned repository \"{}\".", repoPath)
	return nil
}

// DeletePinnedRepoByIndex removes a pinned repository by index
func DeletePinnedRepoByIndex(index int) error {
	repos := viper.GetStringSlice(keyPinnedRepos)

	if index < 0 || index >= len(repos) {
		log.Error("Index out of range: {} (have {} pinned repos).", index, len(repos))
		return nil
	}

	repoPath := repos[index]
	repos = append(repos[:index], repos[index+1:]...)
	viper.Set(keyPinnedRepos, repos)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Deleted pinned repository \"{}\".", repoPath)
	return nil
}

// GetPinnedBranches returns the list of globally pinned branches
func GetPinnedBranches() []string {
	return viper.GetStringSlice(keyPinnedBranches)
}

// GetRepoPinnedBranches returns the map of repo path to pinned branches
func GetRepoPinnedBranches() map[string][]string {
	result := make(map[string][]string)
	raw := viper.GetStringMap(keyRepoPinnedBranches)

	for repoPath, branches := range raw {
		// Normalize the repo path
		normalized, err := path.Normalize(repoPath)
		if err != nil {
			log.Error("Error parsing repo path \"{}\". Error: {}", repoPath, err)
			normalized = repoPath
		}

		// Convert branches to string slice
		switch v := branches.(type) {
		case []interface{}:
			branchList := make([]string, 0, len(v))
			for _, b := range v {
				if s, ok := b.(string); ok {
					branchList = append(branchList, s)
				}
			}
			result[normalized] = branchList
		case []string:
			result[normalized] = v
		}
	}

	return result
}

// GetPinnedBranchesForRepo returns the pinned branches for a specific repo
// This includes both globally pinned branches and repo-specific pinned branches
func GetPinnedBranchesForRepo(repoPath string) []string {
	normalized, _ := path.Normalize(repoPath)

	// Start with global pinned branches
	branches := GetPinnedBranches()

	// Add repo-specific pinned branches
	repoPins := GetRepoPinnedBranches()
	if repoBranches, ok := repoPins[normalized]; ok {
		branches = append(branches, repoBranches...)
	}

	return branches
}

// AddRepoPinnedBranch adds a pinned branch for a specific repository
func AddRepoPinnedBranch(repoPath, branch string) error {
	normalized, err := path.Normalize(repoPath)
	if err != nil {
		log.Error("Error parsing repo path \"{}\". Error: {}", repoPath, err)
		return nil
	}

	repoPins := GetRepoPinnedBranches()

	// Check if already exists
	if branches, ok := repoPins[normalized]; ok {
		for _, b := range branches {
			if b == branch {
				log.Error("Branch \"{}\" already pinned for repo \"{}\".", branch, repoPath)
				return nil
			}
		}
		repoPins[normalized] = append(branches, branch)
	} else {
		repoPins[normalized] = []string{branch}
	}

	viper.Set(keyRepoPinnedBranches, repoPins)
	err = saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Added pinned branch \"{}\" for repo \"{}\".", branch, repoPath)
	return nil
}

// DeleteRepoPinnedBranch removes a pinned branch for a specific repository
func DeleteRepoPinnedBranch(repoPath, branch string) error {
	normalized, _ := path.Normalize(repoPath)

	repoPins := GetRepoPinnedBranches()

	branches, ok := repoPins[normalized]
	if !ok {
		log.Error("No pinned branches found for repo \"{}\".", repoPath)
		return nil
	}

	newBranches := []string{}
	found := false
	for _, b := range branches {
		if b == branch {
			found = true
			continue
		}
		newBranches = append(newBranches, b)
	}

	if !found {
		log.Error("Branch \"{}\" not pinned for repo \"{}\".", branch, repoPath)
		return nil
	}

	if len(newBranches) == 0 {
		delete(repoPins, normalized)
	} else {
		repoPins[normalized] = newBranches
	}

	viper.Set(keyRepoPinnedBranches, repoPins)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Deleted pinned branch \"{}\" for repo \"{}\".", branch, repoPath)
	return nil
}

// IsBranchGloballyPinned checks if a branch is in the global pinned branches list
func IsBranchGloballyPinned(branch string) bool {
	for _, b := range GetPinnedBranches() {
		if b == branch {
			return true
		}
	}
	return false
}

// IsBranchPinnedForRepo checks if a branch is pinned for a specific repo (not global)
func IsBranchPinnedForRepo(repoPath, branch string) bool {
	normalized, _ := path.Normalize(repoPath)
	repoPins := GetRepoPinnedBranches()

	if branches, ok := repoPins[normalized]; ok {
		for _, b := range branches {
			if b == branch {
				return true
			}
		}
	}
	return false
}

// GetAliases returns the map of aliases
func GetAliases() map[string]string {
	return viper.GetStringMapString(keyAliases)
}

// AddAlias adds or updates an alias
func AddAlias(short, expanded string) error {
	aliases := viper.GetStringMapString(keyAliases)
	aliases[short] = expanded
	viper.Set(keyAliases, aliases)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Added alias \"{}\" -> \"{}\".", short, expanded)
	return nil
}

// DeleteAlias removes an alias by its short name
func DeleteAlias(short string) error {
	aliases := viper.GetStringMapString(keyAliases)

	if _, exists := aliases[short]; !exists {
		log.Error("Alias not found: \"{}\".", short)
		return nil
	}

	delete(aliases, short)
	viper.Set(keyAliases, aliases)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Deleted alias \"{}\".", short)
	return nil
}

// saveConfig writes the current configuration to disk
func saveConfig() error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// No config file found, create one in home directory
		home, err := os.UserHomeDir()
		if err != nil {
			log.Error("Could not determine home directory. Error: {}", err)
			return nil
		}
		configFile = filepath.Join(home, ".gimme.config.yaml")
		viper.SetConfigFile(configFile)
	}

	return viper.WriteConfig()
}
