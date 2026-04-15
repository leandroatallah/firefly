//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Story struct {
	ID          string
	Slug        string
	Status      string
	LastAgent   string
	LastModel   string
	LastAction  string
	CurrentStep string
}

func main() {
	activeDir := ".agents/work/active"
	entries, err := os.ReadDir(activeDir)
	if err != nil {
		fmt.Printf("Error reading active stories: %v\n", err)
		return
	}

	// Get current worktree if applicable
	worktree, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	fmt.Printf("\033[1;34m=== SDD PIPELINE DASHBOARD ===\033[0m\n")
	fmt.Printf("\033[1;30mWorktree: %s\033[0m\n\n", strings.TrimSpace(string(worktree)))

	fmt.Printf("%-10s | %-25s | %-15s | %-15s | %s\n", "ID", "STORY", "STEP", "AGENT [MODEL]", "LAST ACTION")
	fmt.Println(strings.Repeat("-", 100))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		story := parseStory(filepath.Join(activeDir, entry.Name()))

		color := "\033[0m" // Default
		if story.CurrentStep == "Feature Implementer" {
			color = "\033[1;32m" // Green for implementation
		} else if story.CurrentStep == "TDD Specialist" {
			color = "\033[1;33m" // Yellow for tests
		} else if strings.Contains(story.CurrentStep, "[/]") || strings.Contains(story.LastAction, "[STARTED]") {
			color = "\033[1;36m" // Cyan for IN PROGRESS
			story.CurrentStep = "⌛ " + strings.Trim(story.CurrentStep, "[/ ]")
		}

		fmt.Printf("%-10s | %-25s | %s%-15s\033[0m | %-15s | %s\n",
			story.ID,
			truncate(story.Slug, 25),
			color,
			story.CurrentStep,
			fmt.Sprintf("%s [%s]", story.LastAgent, story.LastModel),
			truncate(story.LastAction, 40))
	}
}

func parseStory(path string) Story {
	id := filepath.Base(path)
	story := Story{ID: id, Slug: id}

	// Parse PROGRESS.md
	f, err := os.Open(filepath.Join(path, "PROGRESS.md"))
	if err != nil {
		return story
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lastLog string
	pipelineStarted := false

	reLog := regexp.MustCompile(`^- (?:\[?(\d{4}-\d{2}-\d{2})\]?:? )?(?:\[?([^\]]+)\]? )?(?:\[?([^\]]+)\]? )?(.*)`)

	for scanner.Scan() {
		line := scanner.Text()

		// Parse Current Step from Pipeline State checklist
		if strings.Contains(line, "## Pipeline State") {
			pipelineStarted = true
		}
		if pipelineStarted && strings.HasPrefix(line, "- [x]") {
			story.CurrentStep = strings.TrimSpace(strings.TrimPrefix(line, "- [x]"))
		}
		if pipelineStarted && strings.HasPrefix(line, "##") && line != "## Pipeline State" {
			pipelineStarted = false
		}

		// Parse Last Log entry
		if strings.HasPrefix(line, "- ") && !strings.Contains(line, "- [") { // Not a checklist item
			lastLog = line
		}
	}

	if lastLog != "" {
		matches := reLog.FindStringSubmatch(lastLog)
		if len(matches) > 4 {
			// Extracting based on format: - [Model] [Agent] [date]: Action
			// Or standard format: - date: Agent — Action
			if strings.Contains(lastLog, " — ") {
				parts := strings.SplitN(strings.TrimPrefix(lastLog, "- "), " — ", 2)
				header := parts[0]
				story.LastAction = parts[1]

				headerParts := strings.Split(header, ": ")
				if len(headerParts) > 1 {
					story.LastAgent = headerParts[1]
				} else {
					story.LastAgent = headerParts[0]
				}
			} else {
				story.LastAction = matches[4]
				story.LastAgent = matches[2]
				story.LastModel = matches[3]
			}
		}
	}

	return story
}

func truncate(s string, l int) string {
	if len(s) > l {
		return s[:l-3] + "..."
	}
	return s
}

func (s Story) GetModel() string {
	if s.LastModel == "" {
		return "Unknown"
	}
	return s.LastModel
}
