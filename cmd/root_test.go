package cmd

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MoXcz/tasks/internal/config"
)

func TestCLIIntegration(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name  string
		cfg   config.Config
		DoAll bool
	}{
		{
			name:  "CSV cmd",
			DoAll: true,
			cfg: config.Config{
				Filepath: filepath.Join(dir, "tasks"),
				Storage:  "csv",
			},
		},
		{
			name:  "JSON cmd",
			DoAll: true,
			cfg: config.Config{
				Filepath: filepath.Join(dir, "tasks"),
				Storage:  "json",
			},
		},
	}

	for _, tt := range tests {
		task1 := "Rice linux"
		task2 := "Configure Neovim"
		t.Run(tt.name, func(t *testing.T) {
			if tt.DoAll {
				// Add task
				buf := &bytes.Buffer{}
				if err := runCommand(tt.cfg, buf, "add", "Rice linux"); err != nil {
					t.Fatalf("add command failed: %v", err)
				}

				if err := runCommand(tt.cfg, buf, "add", "Configure Neovim"); err != nil {
					t.Fatalf("add command failed: %v", err)
				}

				// List tasks
				buf.Reset()
				if err := runCommand(tt.cfg, buf, "list"); err != nil {
					t.Fatalf("list command failed: %v", err)
				}
				if !strings.Contains(buf.String(), task1) {
					t.Errorf("expected task in list %q, got %q", task1, buf.String())
				}
				if !strings.Contains(buf.String(), task2) {
					t.Errorf("expected task in list %q, got %q", task2, buf.String())
				}

				// Complete task
				buf.Reset()
				if err := runCommand(tt.cfg, buf, "complete", "1"); err != nil {
					t.Fatalf("complete command failed: %v", err)
				}
				if !strings.Contains(buf.String(), "Completing task: Rice linux") {
					t.Errorf("expected completed task %q, got %q", "Rice linux", buf.String())
				}

				// Do not list completed task
				buf.Reset()
				if err := runCommand(tt.cfg, buf, "list"); err != nil {
					t.Fatalf("list command failed: %v", err)
				}
				if strings.Contains(buf.String(), task1) {
					t.Errorf("not expected task in list %q, got %q", task1, buf.String())
				}

				// List completed task
				buf.Reset()
				if err := runCommand(tt.cfg, buf, "list", "--all"); err != nil {
					t.Fatalf("list command failed: %v", err)
				}
				if !strings.Contains(buf.String(), task1) {
					t.Errorf("not expected task in list %q, got %q", task1, buf.String())
				}

				// Delete task
				buf.Reset()
				if err := runCommand(tt.cfg, buf, "delete", "1", "--force"); err != nil {
					t.Fatalf("delete command failed: %v", err)
				}
				if !strings.Contains(buf.String(), "Deleting task: Rice linux") {
					t.Errorf("expected deleted task %q, got %q", "Rice linux", buf.String())
				}

				// Do not list deleted task
				buf.Reset()
				if err := runCommand(tt.cfg, buf, "list", "--all"); err != nil {
					t.Fatalf("list command failed: %v", err)
				}
				if strings.Contains(buf.String(), task1) {
					t.Errorf("not expected task in list %q, got %q", task1, buf.String())
				}
			}
		})
	}

}

func runCommand(cfg config.Config, out io.Writer, args ...string) error {
	cmd := NewRootCmd(cfg)
	cmd.SetOut(out)
	cmd.SetArgs(args)
	return cmd.Execute()
}
