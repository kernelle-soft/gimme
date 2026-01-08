package config

import (
	"github.com/kernelle-soft/gimme/internal/log"
	"github.com/kernelle-soft/gimme/internal/path"
	"github.com/kernelle-soft/gimme/internal/slice"
	"github.com/spf13/viper"
)

func Load() {
	// Set defaults
	viper.SetDefault("search-folders", []string{"~/"})

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

func GetSearchFolders() []string {
	rawPaths := viper.GetStringSlice("search-folders")

	return slice.Map(rawPaths, func(rawPath string) string {
		normalized, err := path.Normalize(rawPath)
		if err != nil {
			log.Error("Error parsing search folder", "path", rawPath, "error", err)
		}
		return normalized
	})
}
