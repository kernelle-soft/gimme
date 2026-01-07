package config

import (
	"fmt"
	"os"

	"github.com/kernelle-soft/gimmetool/internal/path"
	"github.com/kernelle-soft/gimmetool/internal/slice"
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
			fmt.Fprintln(os.Stderr, "Error reading gimme configuration:", err)
		}
	}
}

func GetSearchFolders() []string {
	rawPaths := viper.GetStringSlice("search-folders")

	return slice.Map(rawPaths, func (rawPath string) string {
		normalized, err := path.Normalize(rawPath)
		if err != nil {
			_ = fmt.Errorf("Error parsing search folder %s: %w", rawPath, err)
		}
		return normalized
	})
}
