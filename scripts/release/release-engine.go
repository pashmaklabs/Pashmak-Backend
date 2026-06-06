package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SemVer struct {
	Major      int
	Minor      int
	Patch      int
	HasVPrefix bool
}

func parseVersion(tag string) SemVer {
	hasV := strings.HasPrefix(tag, "v")
	versionStr := tag
	if hasV {
		versionStr = tag[1:]
	}

	parts := strings.Split(versionStr, ".")
	major, minor, patch := 0, 0, 0
	if len(parts) >= 1 {
		major, _ = strconv.Atoi(parts[0])
	}
	if len(parts) >= 2 {
		minor, _ = strconv.Atoi(parts[1])
	}
	if len(parts) >= 3 {
		patchStr := parts[2]
		if idx := strings.IndexAny(patchStr, "-+"); idx != -1 {
			patchStr = patchStr[:idx]
		}
		patch, _ = strconv.Atoi(patchStr)
	}

	return SemVer{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		HasVPrefix: hasV,
	}
}

func (s SemVer) String() string {
	prefix := ""
	if s.HasVPrefix {
		prefix = "v"
	}
	return fmt.Sprintf("%s%d.%d.%d", prefix, s.Major, s.Minor, s.Patch)
}

func (s SemVer) RawString() string {
	return fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
}

func (s SemVer) Bump(bumpType string) SemVer {
	bumped := s
	switch strings.ToLower(bumpType) {
	case "major":
		bumped.Major++
		bumped.Minor = 0
		bumped.Patch = 0
	case "minor":
		bumped.Minor++
		bumped.Patch = 0
	case "patch":
		bumped.Patch++
	default:
		bumped.Patch++
	}
	return bumped
}

func getLastTag() (string, bool) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "0.0.0", false
	}
	return strings.TrimSpace(string(output)), true
}

func getCommits(lastTag string, hasTag bool) ([]string, error) {
	var cmd *exec.Cmd
	if hasTag {
		cmd = exec.Command("git", "log", fmt.Sprintf("%s..HEAD", lastTag), "--oneline")
	} else {
		cmd = exec.Command("git", "log", "--oneline")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var commits []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			commits = append(commits, trimmed)
		}
	}
	return commits, nil
}

func main() {
	bumpTypeOpt := flag.String("bump", "patch", "SemVer bump type (patch, minor, major)")
	customNotesOpt := flag.String("notes", "", "Optional custom release notes")
	flag.Parse()

	// Fallback to positional argument for bump type if provided
	bumpType := *bumpTypeOpt
	if flag.NArg() > 0 {
		bumpType = flag.Arg(0)
	}

	// Validate bump type
	bumpType = strings.ToLower(strings.TrimSpace(bumpType))
	if bumpType != "patch" && bumpType != "minor" && bumpType != "major" {
		fmt.Printf("Warning: invalid bump type '%s'. Defaulting to 'patch'.\n", bumpType)
		bumpType = "patch"
	}

	customNotes := strings.TrimSpace(*customNotesOpt)

	// Get last tag
	lastTag, hasTag := getLastTag()
	fmt.Printf("Current Tag: %s (Has Tag: %t)\n", lastTag, hasTag)

	// Parse SemVer and Bump version
	currentVersion := parseVersion(lastTag)
	newVersion := currentVersion.Bump(bumpType)
	newTag := newVersion.String()

	fmt.Printf("New Version: %s\n", newVersion.RawString())
	fmt.Printf("New Tag: %s\n", newTag)

	// Get commit history
	commits, err := getCommits(lastTag, hasTag)
	if err != nil {
		fmt.Printf("Error running git log: %v. Continuing with empty commit list.\n", err)
		commits = []string{}
	}

	// Parse commits
	categories := map[string][]string{
		"### Added":      {},
		"### Fixed":      {},
		"### Changed":    {},
		"### Refactored": {},
		"### Other":      {},
	}

	commitHashRegex := regexp.MustCompile(`^([a-fA-F0-9]+)\s+(.*)$`)

	// Conventional commit matching regexes
	// We extract: scope (optional, group 1), breaking indicator (!, optional, group 2), subject/message (group 3)
	addedRegex := regexp.MustCompile(`^(?i)(?:feat|feature)(?:\(([^)]+)\))?(!?)(?::|\s+)\s*(.*)$`)
	fixedRegex := regexp.MustCompile(`^(?i)(?:fix|bugfix)(?:\(([^)]+)\))?(!?)(?::|\s+)\s*(.*)$`)
	changedRegex := regexp.MustCompile(`^(?i)(?:change|refactor)(?:\(([^)]+)\))?(!?)(?::|\s+)\s*(.*)$`)

	// Normalize trailing issues/PRs (e.g. #12) and strip #0 / (#0)
	stripZeroRegex := regexp.MustCompile(`\s*\(?#0\)?\s*$`)
	formatIdRegex := regexp.MustCompile(`\s+#([1-9][0-9]*)\s*$`)

	for _, commit := range commits {
		match := commitHashRegex.FindStringSubmatch(commit)
		if len(match) < 3 {
			// If it doesn't match standard git log oneline output hash + space + subject, classify as other
			categories["### Other"] = append(categories["### Other"], "- "+commit)
			continue
		}
		// commitHash := match[1] // Unused, we only need the subject
		subject := strings.TrimSpace(match[2])

		// Apply trailing #id normalizations
		subject = stripZeroRegex.ReplaceAllString(subject, "")
		subject = formatIdRegex.ReplaceAllString(subject, " (#$1)")

		// Categorize
		if m := addedRegex.FindStringSubmatch(subject); len(m) >= 4 {
			scope, breaking, msg := m[1], m[2], m[3]
			categories["### Added"] = append(categories["### Added"], formatItem(scope, breaking, msg))
		} else if m := fixedRegex.FindStringSubmatch(subject); len(m) >= 4 {
			scope, breaking, msg := m[1], m[2], m[3]
			categories["### Fixed"] = append(categories["### Fixed"], formatItem(scope, breaking, msg))
		} else if m := changedRegex.FindStringSubmatch(subject); len(m) >= 4 {
			scope, breaking, msg := m[1], m[2], m[3]
			categories["### Changed"] = append(categories["### Changed"], formatItem(scope, breaking, msg))
		} else {
			categories["### Other"] = append(categories["### Other"], "- "+subject)
		}
	}

	// Build release notes body (without version header for GitHub Releases)
	var bodySb strings.Builder
	if customNotes != "" {
		bodySb.WriteString(customNotes)
		bodySb.WriteString("\n")
	}

	order := []string{"### Added", "### Fixed", "### Changed", "### Other"}
	for _, sec := range order {
		items := categories[sec]
		if len(items) > 0 {
			if bodySb.Len() > 0 {
				bodySb.WriteString("\n")
			}
			bodySb.WriteString(sec)
			bodySb.WriteString("\n\n")
			for _, item := range items {
				bodySb.WriteString(item)
				bodySb.WriteString("\n")
			}
		}
	}
	releaseNotesBody := bodySb.String()

	// Build complete changelog block (including version header)
	dateStr := time.Now().UTC().Format("2006-01-02")
	var changelogBlock string
	if strings.TrimSpace(releaseNotesBody) != "" {
		changelogBlock = fmt.Sprintf("## [%s] - %s\n\n%s", newVersion.RawString(), dateStr, strings.TrimSpace(releaseNotesBody))
	} else {
		changelogBlock = fmt.Sprintf("## [%s] - %s\n", newVersion.RawString(), dateStr)
	}

	// Update CHANGELOG.md
	if err := updateChangelog(changelogBlock); err != nil {
		fmt.Printf("Error updating CHANGELOG.md: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("CHANGELOG.md updated successfully.")

	// Output to GITHUB_OUTPUT if applicable
	if err := writeGithubOutputs(newVersion.RawString(), newTag, releaseNotesBody); err != nil {
		fmt.Printf("Error writing GITHUB_OUTPUT: %v\n", err)
		os.Exit(1)
	}
}

func formatItem(scope, breaking, msg string) string {
	res := ""
	if scope != "" {
		res += fmt.Sprintf("**%s**: ", scope)
	}
	if breaking == "!" {
		res += "[BREAKING] "
	}
	res += msg
	return "- " + res
}

func updateChangelog(newReleaseBlock string) error {
	changelogPath := "CHANGELOG.md"
	header := "# Changelog\n\nAll notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).\n"

	var content string
	if _, err := os.Stat(changelogPath); err == nil {
		data, err := os.ReadFile(changelogPath)
		if err != nil {
			return err
		}
		content = string(data)
	} else {
		content = header
	}

	re := regexp.MustCompile(`(?m)^##\s+\[`)
	loc := re.FindStringIndex(content)

	var newContent string
	if loc != nil {
		insertPos := loc[0]
		before := strings.TrimRight(content[:insertPos], "\r\n\t ")
		after := content[insertPos:]
		newContent = fmt.Sprintf("%s\n\n%s\n\n---\n\n%s", before, strings.TrimSpace(newReleaseBlock), after)
	} else {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		newContent = fmt.Sprintf("%s\n%s\n", content, strings.TrimSpace(newReleaseBlock))
	}

	return os.WriteFile(changelogPath, []byte(newContent), 0644)
}

func writeGithubOutputs(version, tag, releaseNotes string) error {
	githubOutput := os.Getenv("GITHUB_OUTPUT")
	if githubOutput == "" {
		return nil
	}

	file, err := os.OpenFile(githubOutput, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "version=%s\n", version); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(file, "tag=%s\n", tag); err != nil {
		return err
	}

	if _, err := fmt.Fprint(file, "release_notes<<EOF\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(file, releaseNotes); err != nil {
		return err
	}
	if _, err := fmt.Fprint(file, "\nEOF\n"); err != nil {
		return err
	}

	return nil
}
