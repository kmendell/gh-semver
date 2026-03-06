package semver

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestIsRelevantCommitUsesFilterPath(t *testing.T) {
	repo, repoDir := initTestRepo(t)
	cc := NewConventionalCommits(repo, "pkg/", "")

	outsideCommit := commitFile(t, repo, repoDir, "docs/readme.md", "hello", "docs: update readme")
	insideCommit := commitFile(t, repo, repoDir, "pkg/service.txt", "hello", "feat: add service")

	if cc.isRelevantCommit(outsideCommit) {
		t.Fatal("expected commit outside filter-path to be ignored")
	}

	if !cc.isRelevantCommit(insideCommit) {
		t.Fatal("expected commit inside filter-path to be included")
	}
}

func initTestRepo(t *testing.T) (*git.Repository, string) {
	t.Helper()

	repoDir := t.TempDir()
	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatalf("git.PlainInit() error = %v", err)
	}

	return repo, repoDir
}

func commitFile(t *testing.T, repo *git.Repository, repoDir, name, content, message string) *object.Commit {
	t.Helper()

	fullPath := filepath.Join(repoDir, name)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("repo.Worktree() error = %v", err)
	}
	if _, err := worktree.Add(name); err != nil {
		t.Fatalf("worktree.Add() error = %v", err)
	}

	hash, err := worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("worktree.Commit() error = %v", err)
	}

	commit, err := repo.CommitObject(hash)
	if err != nil {
		t.Fatalf("repo.CommitObject() error = %v", err)
	}

	return commit
}