// Copyright 2022 Jan van den Berg
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/kmendell/gh-semver/internal/semver"
)

func main() {
	var (
		action     bool
		filterPath string
		prefix     string
		release    bool
		tag        bool
	)
	flag.BoolVar(&action, "action", false, "GitHub Action output format named 'version'")
	flag.StringVar(&filterPath, "filter-path", "", "The path to filter commits (in case of a mono-repo)")
	flag.StringVar(&prefix, "prefix", "", "The prefix of the tag (in case of a mono-repo)")
	flag.BoolVar(&release, "release", false, "Force release tag")
	flag.BoolVar(&tag, "tag", false, "Commit the tag")
	flag.Parse()

	// open current repository
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't open git repository: %v\n", err)
		os.Exit(1)
	}

	tagVersion := calculateSemVer(repo, filterPath, prefix, action, release)
	if tag {
		gitTag(repo, tagVersion)
	}

	if err := emitVersion(tagVersion, action); err != nil {
		fmt.Fprintf(os.Stderr, "error writing version output: %v\n", err)
		os.Exit(1)
	}
}

func emitVersion(tagVersion string, action bool) error {
	if action {
		if err := writeGitHubActionOutput("version", tagVersion); err != nil {
			return err
		}
	}

	fmt.Println(tagVersion)
	return nil
}

func writeGitHubActionOutput(name, value string) error {
	outputPath := os.Getenv("GITHUB_OUTPUT")
	if outputPath == "" {
		return nil
	}

	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("couldn't open GITHUB_OUTPUT: %w", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "%s=%s\n", name, value); err != nil {
		return fmt.Errorf("couldn't write GITHUB_OUTPUT: %w", err)
	}

	return nil
}

func calculateSemVer(repo *git.Repository, filterPath, prefix string, action, release bool) string {
	conventionalCommits := semver.NewConventionalCommits(repo, filterPath, prefix)
	nextVersion, err := conventionalCommits.SemVer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	nextVersion.Prefix = prefix

	return nextVersion.PrintTag(release)
}

func gitTag(repo *git.Repository, tagVersion string) {
	if _, err := repo.Tag(tagVersion); err != nil {
		headRef, err := repo.Head()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error determining tag: %v\n", err)
			os.Exit(1)
		}
		if _, err = repo.CreateTag(tagVersion, headRef.Hash(), &git.CreateTagOptions{
			Message: tagVersion,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "error creating tag: %v\v", err)
			os.Exit(1)
		}
	}
}
