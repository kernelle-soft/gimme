package config

import (
	"os"
	"path/filepath"

	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/kernelle-soft/gimme/internal/path"
	"github.com/kernelle-soft/gimme/internal/slice"
	"github.com/spf13/viper"
)

// Config keys using nested structure:
//
//	search-folders: [...]
//	pins:
//	  repositories: [...]
//	  branches:
//	    global: [main, master]
//	    repositories:
//	      github.com/user/repo: [branch1, branch2]
//	aliases: {...}
const (
	keySearchFolders = "search-folders"
	keyAliases       = "aliases"

	// Nested pins keys
	keyPinsRepositories        = "pins.repositories"
	keyPinsBranchesGlobal      = "pins.branches.global"
	keyPinsBranchesRepositores = "pins.branches.repositories"
)

// Defaults
var defaultSearchFolder = "~/"
var defaultPinnedGlobalBranches = []string{"main", "master"}

// isDefaultSearchFolder checks if the given groups match the default
func isDefaultSearchFolder(groups []string) bool {
	return groups[0] == defaultSearchFolder && len(groups) == 1
}

func Load() {
	// Set defaults
	viper.SetDefault(keySearchFolders, []string{defaultSearchFolder})
	viper.SetDefault(keyPinsRepositories, []string{})
	viper.SetDefault(keyPinsBranchesGlobal, defaultPinnedGlobalBranches)
	viper.SetDefault(keyPinsBranchesRepositores, map[string][]string{})
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
			log.Print("Group already exists: \"{}\".", groupPath)
			return nil
		}
	}

	// If the only group is the default, replace it instead of appending
	if isDefaultSearchFolder(groups) {
		groups = []string{groupPath}
	} else {
		groups = append(groups, groupPath)
	}

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
		log.Print("Group not found: \"{}\".", groupPath)
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
		log.Print("Index out of range: {} (have {} groups).", index, len(groups))
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

// =============================================================================
// Pinned Repositories (pins.repositories)
// =============================================================================

// GetPinnedRepos returns the list of pinned repository paths
func GetPinnedRepos() []string {
	rawPaths := viper.GetStringSlice(keyPinsRepositories)

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
	repos := viper.GetStringSlice(keyPinsRepositories)

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
			log.Print("Pinned repo already exists: \"{}\".", repoPath)
			return nil
		}
	}

	repos = append(repos, repoPath)
	viper.Set(keyPinsRepositories, repos)
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
	repos := viper.GetStringSlice(keyPinsRepositories)
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
		log.Print("Pinned repo not found: \"{}\".", repoPath)
		return nil
	}

	viper.Set(keyPinsRepositories, newRepos)
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
	repos := viper.GetStringSlice(keyPinsRepositories)

	if index < 0 || index >= len(repos) {
		log.Print("Index out of range: {} (have {} pinned repos).", index, len(repos))
		return nil
	}

	repoPath := repos[index]
	repos = append(repos[:index], repos[index+1:]...)
	viper.Set(keyPinsRepositories, repos)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Deleted pinned repository \"{}\".", repoPath)
	return nil
}

// =============================================================================
// Pinned Branches - Global (pins.branches.global)
// =============================================================================

// GetGlobalPinnedBranches returns the list of globally pinned branches
func GetGlobalPinnedBranches() []string {
	return viper.GetStringSlice(keyPinsBranchesGlobal)
}

// IsBranchGloballyPinned checks if a branch is in the global pinned branches list
func IsBranchGloballyPinned(branch string) bool {
	for _, b := range GetGlobalPinnedBranches() {
		if b == branch {
			return true
		}
	}
	return false
}

// =============================================================================
// Pinned Branches - Per-Repo (pins.branches.repositories)
// Uses repo identifier (e.g. "github.com/user/repo") as key
// =============================================================================

// GetRepoPinnedBranches returns the map of repo identifier to pinned branches
func GetRepoPinnedBranches() map[string][]string {
	result := make(map[string][]string)
	raw := viper.GetStringMap(keyPinsBranchesRepositores)

	for repoID, branches := range raw {
		// Convert branches to string slice
		switch v := branches.(type) {
		case []interface{}:
			branchList := make([]string, 0, len(v))
			for _, b := range v {
				if s, ok := b.(string); ok {
					branchList = append(branchList, s)
				}
			}
			result[repoID] = branchList
		case []string:
			result[repoID] = v
		}
	}

	return result
}

// GetPinnedBranchesForRepo returns all pinned branches for a specific repo
// (both global and repo-specific)
func GetPinnedBranchesForRepo(repoIdentifier string) []string {
	// Start with global pinned branches
	branches := GetGlobalPinnedBranches()

	// Add repo-specific pinned branches
	repoPins := GetRepoPinnedBranches()
	if repoBranches, ok := repoPins[repoIdentifier]; ok {
		branches = append(branches, repoBranches...)
	}

	return branches
}

// AddRepoPinnedBranch adds a pinned branch for a specific repository
func AddRepoPinnedBranch(repoIdentifier, branch string) error {
	repoPins := GetRepoPinnedBranches()

	// Check if already exists
	if branches, ok := repoPins[repoIdentifier]; ok {
		for _, b := range branches {
			if b == branch {
				log.Print("Branch \"{}\" already pinned for repo \"{}\".", branch, repoIdentifier)
				return nil
			}
		}
		repoPins[repoIdentifier] = append(branches, branch)
	} else {
		repoPins[repoIdentifier] = []string{branch}
	}

	viper.Set(keyPinsBranchesRepositores, repoPins)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Added pinned branch \"{}\" for repo \"{}\".", branch, repoIdentifier)
	return nil
}

// DeleteRepoPinnedBranch removes a pinned branch for a specific repository
func DeleteRepoPinnedBranch(repoIdentifier, branch string) error {
	repoPins := GetRepoPinnedBranches()

	branches, ok := repoPins[repoIdentifier]
	if !ok {
		log.Print("No pinned branches found for repo \"{}\".", repoIdentifier)
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
		log.Print("Branch \"{}\" not pinned for repo \"{}\".", branch, repoIdentifier)
		return nil
	}

	if len(newBranches) == 0 {
		delete(repoPins, repoIdentifier)
	} else {
		repoPins[repoIdentifier] = newBranches
	}

	viper.Set(keyPinsBranchesRepositores, repoPins)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Deleted pinned branch \"{}\" for repo \"{}\".", branch, repoIdentifier)
	return nil
}

// IsBranchPinnedForRepo checks if a branch is pinned for a specific repo (not global)
func IsBranchPinnedForRepo(repoIdentifier, branch string) bool {
	repoPins := GetRepoPinnedBranches()

	if branches, ok := repoPins[repoIdentifier]; ok {
		for _, b := range branches {
			if b == branch {
				return true
			}
		}
	}
	return false
}

// =============================================================================
// Aliases
// =============================================================================

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
		log.Print("Alias not found: \"{}\".", short)
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

// =============================================================================
// Config persistence
// =============================================================================

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
