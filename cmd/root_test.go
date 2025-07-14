package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MoXcz/tasks/internal/config"
)

func TestCLIIntegration(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name     string
		cfg      config.Config
		Add      bool
		List     bool
		Delete   bool
		Complete bool
	}{
		{
			name: "CSV add and list",
			Add:  true,
			List: true,
			cfg: config.Config{
				Filepath: filepath.Join(dir, "tasks"),
				Storage:  "csv",
			},
		},
		{
			name: "JSON add",
			Add:  true,
			cfg: config.Config{
				Filepath: filepath.Join(dir, "tasks"),
				Storage:  "json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd := NewRootCmd(tt.cfg)
			if tt.Add {
				buf := &bytes.Buffer{}
				rootCmd.SetOut(buf)
				rootCmd.SetArgs([]string{
					"add", "Buy milk",
				})
				if err := rootCmd.Execute(); err != nil {
					t.Fatalf("add command failed: %v", err)
				}
			}
			if tt.List {
				buf := &bytes.Buffer{}
				rootCmd.SetOut(buf)
				rootCmd.SetArgs([]string{
					"list",
				})
				if err := rootCmd.Execute(); err != nil {
					t.Fatalf("list command failed: %v", err)
				}

				if !strings.Contains(buf.String(), "Buy milk") {
					t.Errorf("expected task in list, got %q", buf.String())
				}
			}
		})
	}

}
