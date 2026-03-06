package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteGitHubActionOutputWritesNamedValue(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "github_output.txt")
	t.Setenv("GITHUB_OUTPUT", outputPath)

	if err := writeGitHubActionOutput("version", "v1.2.3"); err != nil {
		t.Fatalf("writeGitHubActionOutput() error = %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	if got, want := string(content), "version=v1.2.3\n"; got != want {
		t.Fatalf("writeGitHubActionOutput() = %q, want %q", got, want)
	}
}

func TestWriteGitHubActionOutputNoopWithoutEnvironment(t *testing.T) {
	t.Setenv("GITHUB_OUTPUT", "")

	if err := writeGitHubActionOutput("version", "v1.2.3"); err != nil {
		t.Fatalf("writeGitHubActionOutput() error = %v", err)
	}
}