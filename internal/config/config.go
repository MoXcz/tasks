package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Filepath string
	Verbose  bool
	Storage  string
}

func Load(cfgFile string) (Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return Config{}, fmt.Errorf("cannot determine user home: %w", err)
		}
		viper.AddConfigPath(filepath.Join(home, ".config", "tasks"))
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("tasks.yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := writeDefaultConfig(); err != nil {
				return Config{}, fmt.Errorf("writing default config failed: %w", err)
			}
		} else {
			return Config{}, fmt.Errorf("error reading config: %w", err)
		}
	}

	return Config{
		Filepath: viper.GetString("filepath"),
		Verbose:  viper.GetBool("verbose"),
		Storage:  viper.GetString("storage"),
	}, nil
}

func writeDefaultConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("could not find user config directory: %w", err)
	}
	dir := filepath.Join(configDir, "tasks")

	// file for tasks is found at the same config directory by default,
	// this is a personal choice
	// TODO: When storage type changes "tasks" breaks, easy solution is to change
	// file name depending on the storage type
	// Maybe wait until import/export? it could be easy once that's implemented
	// to just export tasks to the new storage type and the import it in the same
	// file
	viper.SetDefault("filepath", filepath.Join(dir, "tasks.csv"))
	viper.SetDefault("verbose", false)
	viper.SetDefault("storage", "csv") // TODO: change default to sqlite

	if err := os.MkdirAll(dir, 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	viper.SetConfigFile(filepath.Join(dir, "tasks.yaml"))
	return viper.WriteConfigAs(viper.ConfigFileUsed())
}
