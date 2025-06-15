/*
Copyright © 2025 Oscar Marquez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/MoXcz/tasks/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var storage file.FileStorage

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "A task management CLI that tries to mimic the common TODO web application",
	Long: `tasks is a CLI tool to manage your tasks, that's it.
tasks add <task> to add a new task
tasks complete <task id> to mark a task as completed
tasks list to list all tasks
`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/tasks/tasks.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home + "/.config/tasks")
		viper.AddConfigPath(".")
		// Search config in home directory with name ".tasks" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("tasks")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			err = writeDefaultConfig()
			if err != nil {
				fmt.Fprintln(os.Stderr, "error writing default config", err)
				os.Exit(1)
			}
			fmt.Fprintln(os.Stdout, "config file created at:", viper.ConfigFileUsed())
		} else {
			// Config file was found but another error was produced
			fmt.Fprintln(os.Stderr, "error found in config:", viper.ConfigFileUsed())
			os.Exit(1)
		}
	}

	// Config file found and successfully parsed
	verbose := viper.GetBool("verbose")
	if verbose == true {
		fmt.Fprintln(os.Stdout, "using config file:", viper.ConfigFileUsed())
	}

	storageType := viper.GetString("storage")
	filepath := viper.GetString("filepath")
	s, err := file.SelectStorage(filepath, storageType)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error selecting storage:", err)
		os.Exit(1) // TODO: most probably an invalid storage type, verify it
	}
	storage = s
}

// writes default configuration to the user's config directory
func writeDefaultConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("could not find user config directory: %w", err)
	}
	filepath := configDir + "/tasks"

	// tasks.csv file for tasks is found at the same config directory by default,
	// this is a personal choice
	viper.SetDefault("filepath", filepath+"/tasks.csv")
	viper.SetDefault("verbose", false)
	viper.SetDefault("storage", "csv")

	err = os.Mkdir(filepath, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	viper.SetConfigFile(filepath + "/tasks.yaml")
	err = viper.WriteConfigAs(viper.ConfigFileUsed())
	if err != nil {
		return fmt.Errorf("could not write default config file: %w", err)
	}

	return nil
}
