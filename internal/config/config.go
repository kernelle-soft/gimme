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
	keySearchFolders = "search-folders"
	keyPins          = "pins"
	keyAliases       = "aliases"
)

func Load() {
	// Set defaults
	viper.SetDefault(keySearchFolders, []string{"~/"})
	viper.SetDefault(keyPins, []string{})
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

// GetPins returns the list of pinned repositories
func GetPins() []string {
	rawPaths := viper.GetStringSlice(keyPins)

	return slice.Map(rawPaths, func(rawPath string) string {
		normalized, err := path.Normalize(rawPath)
		if err != nil {
			log.Error("Error parsing pinned path \"{}\". Error: {}", rawPath, err)
			return rawPath
		}
		return normalized
	})
}

// AddPin adds a pinned repository path
func AddPin(pinPath string) error {
	pins := viper.GetStringSlice(keyPins)

	// Normalize the path
	normalized, err := path.Normalize(pinPath)
	if err != nil {
		log.Error("Error parsing pinned path \"{}\". Error: {}", pinPath, err)
		return nil
	}

	// Check if already exists
	for _, p := range pins {
		existingNorm, _ := path.Normalize(p)
		if existingNorm == normalized {
			log.Error("Pin already exists: \"{}\".", pinPath)
			return nil
		}
	}

	pins = append(pins, pinPath)
	viper.Set(keyPins, pins)
	err = saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Added pinned repository \"{}\".", pinPath)
	return nil
}

// DeletePin removes a pinned repository by path
func DeletePin(pinPath string) error {
	pins := viper.GetStringSlice(keyPins)
	normalized, _ := path.Normalize(pinPath)

	newPins := []string{}
	found := false
	for _, p := range pins {
		existingNorm, _ := path.Normalize(p)
		if existingNorm == normalized {
			found = true
			continue
		}
		newPins = append(newPins, p)
	}

	if !found {
		log.Error("Pin not found: \"{}\".", pinPath)
		return nil
	}

	viper.Set(keyPins, newPins)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Deleted pinned repository \"{}\".", pinPath)
	return nil
}

// DeletePinByIndex removes a pinned repository by index
func DeletePinByIndex(index int) error {
	pins := viper.GetStringSlice(keyPins)

	if index < 0 || index >= len(pins) {
		log.Error("Index out of range: {} (have {} pins).", index, len(pins))
		return nil
	}

	pinPath := pins[index]
	pins = append(pins[:index], pins[index+1:]...)
	viper.Set(keyPins, pins)
	err := saveConfig()
	if err != nil {
		log.Error("Error saving config. Error: {}", err)
		return nil
	}
	log.Print("Deleted pinned repository \"{}\".", pinPath)
	return nil
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
