package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	workspaceDir = ".bashlog-workspaces"
	configFile   = "config.txt"
)

// Workspace represents a bash logging workspace
type Workspace struct {
	Name      string
	CreatedAt time.Time
	Path      string
	CommandCount int
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not determine home directory: %v\n", err)
		os.Exit(1)
	}

	basePath := filepath.Join(homeDir, workspaceDir)

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "list":
		handleList(basePath)
	case "create":
		handleCreate(basePath, args)
	case "delete":
		handleDelete(basePath, args)
	case "view":
		handleView(basePath, args)
	case "stats":
		handleStats(basePath, args)
	case "history":
		handleHistory(basePath, args)
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

// handleList displays all available workspaces
func handleList(basePath string) {
	workspaces, err := getWorkspaces(basePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading workspaces: %v\n", err)
		os.Exit(1)
	}

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found. Create one with: bashlog-mgr create <name>")
		return
	}

	fmt.Printf("%-20s %-19s %-10s %s\n", "NAME", "CREATED", "COMMANDS", "PATH")
	fmt.Println(strings.Repeat("-", 70))

	for _, ws := range workspaces {
		fmt.Printf("%-20s %-19s %-10d %s\n",
			ws.Name,
			ws.CreatedAt.Format("2006-01-02 15:04:05"),
			ws.CommandCount,
			ws.Path)
	}
}

// handleCreate creates a new workspace
func handleCreate(basePath string, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: workspace name required\n")
		fmt.Fprintf(os.Stderr, "Usage: bashlog-mgr create <name>\n")
		os.Exit(1)
	}

	name := args[0]

	// Validate workspace name
	if !isValidName(name) {
		fmt.Fprintf(os.Stderr, "Error: invalid workspace name '%s'\n", name)
		fmt.Fprintf(os.Stderr, "Names must contain only alphanumeric characters, hyphens, and underscores\n")
		os.Exit(1)
	}

	wsPath := filepath.Join(basePath, name)

	// Check if workspace already exists
	if _, err := os.Stat(wsPath); err == nil {
		fmt.Fprintf(os.Stderr, "Error: workspace '%s' already exists\n", name)
		os.Exit(1)
	}

	// Create workspace directory structure
	if err := os.MkdirAll(wsPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating workspace: %v\n", err)
		os.Exit(1)
	}

	// Create config file
	configPath := filepath.Join(wsPath, configFile)
	config := fmt.Sprintf("name=%s\ncreated=%s\ncommands=0\n",
		name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config file: %v\n", err)
		os.Exit(1)
	}

	// Create history file
	historyPath := filepath.Join(wsPath, "history.log")
	if err := os.WriteFile(historyPath, []byte(""), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating history file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Workspace '%s' created successfully at %s\n", name, wsPath)
}

// handleDelete removes a workspace
func handleDelete(basePath string, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: workspace name required\n")
		fmt.Fprintf(os.Stderr, "Usage: bashlog-mgr delete <name>\n")
		os.Exit(1)
	}

	name := args[0]
	wsPath := filepath.Join(basePath, name)

	// Check if workspace exists
	if _, err := os.Stat(wsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: workspace '%s' not found\n", name)
		os.Exit(1)
	}

	// Confirm deletion
	fmt.Printf("Are you sure you want to delete workspace '%s'? (yes/no): ", name)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "yes" && response != "y" {
		fmt.Println("Deletion cancelled")
		return
	}

	if err := os.RemoveAll(wsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting workspace: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Workspace '%s' deleted successfully\n", name)
}

// handleView displays workspace details
func handleView(basePath string, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: workspace name required\n")
		fmt.Fprintf(os.Stderr, "Usage: bashlog-mgr view <name>\n")
		os.Exit(1)
	}

	name := args[0]
	wsPath := filepath.Join(basePath, name)

	// Check if workspace exists
	if _, err := os.Stat(wsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: workspace '%s' not found\n", name)
		os.Exit(1)
	}

	// Read config
	config := readConfig(filepath.Join(wsPath, configFile))

	fmt.Printf("\n=== Workspace: %s ===\n", name)
	fmt.Printf("Path: %s\n", wsPath)
	fmt.Printf("Created: %s\n", config["created"])
	fmt.Printf("Commands Logged: %s\n", config["commands"])

	// Show recent history
	historyPath := filepath.Join(wsPath, "history.log")
	if data, err := os.ReadFile(historyPath); err == nil {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		count := len(lines)
		if count > 0 && lines[0] != "" {
			fmt.Printf("\nRecent Commands (last 5):\n")
			start := count - 5
			if start < 0 {
				start = 0
			}
			for i := start; i < count; i++ {
				if lines[i] != "" {
					fmt.Printf("  %s\n", lines[i])
				}
			}
		}
	}
	fmt.Println()
}

// handleStats displays workspace statistics
func handleStats(basePath string, args []string) {
	workspaces, err := getWorkspaces(basePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading workspaces: %v\n", err)
		os.Exit(1)
	}

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found")
		return
	}

	totalCommands := 0
	oldestWorkspace := workspaces[0]
	newestWorkspace := workspaces[0]

	for _, ws := range workspaces {
		totalCommands += ws.CommandCount
		if ws.CreatedAt.Before(oldestWorkspace.CreatedAt) {
			oldestWorkspace = ws
		}
		if ws.CreatedAt.After(newestWorkspace.CreatedAt) {
			newestWorkspace = ws
		}
	}

	fmt.Println("\n=== Workspace Statistics ===")
	fmt.Printf("Total Workspaces: %d\n", len(workspaces))
	fmt.Printf("Total Commands Logged: %d\n", totalCommands)
	if len(workspaces) > 0 {
		fmt.Printf("Average Commands per Workspace: %.2f\n", float64(totalCommands)/float64(len(workspaces)))
	}
	fmt.Printf("Oldest Workspace: %s (created %s)\n", oldestWorkspace.Name, oldestWorkspace.CreatedAt.Format("2006-01-02"))
	fmt.Printf("Newest Workspace: %s (created %s)\n", newestWorkspace.Name, newestWorkspace.CreatedAt.Format("2006-01-02"))
	fmt.Println()
}

// handleHistory displays command history for a workspace
func handleHistory(basePath string, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: workspace name required\n")
		fmt.Fprintf(os.Stderr, "Usage: bashlog-mgr history <name> [lines]\n")
		os.Exit(1)
	}

	name := args[0]
	wsPath := filepath.Join(basePath, name)

	// Check if workspace exists
	if _, err := os.Stat(wsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: workspace '%s' not found\n", name)
		os.Exit(1)
	}

	// Parse number of lines to display (default: 20)
	lines := 20
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &lines)
	}

	historyPath := filepath.Join(wsPath, "history.log")
	data, err := os.ReadFile(historyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading history: %v\n", err)
		os.Exit(1)
	}

	if len(data) == 0 {
		fmt.Printf("No command history for workspace '%s'\n", name)
		return
	}

	historyLines := strings.Split(strings.TrimSpace(string(data)), "\n")
	
	// Display last N lines
	fmt.Printf("\n=== Command History for '%s' (last %d commands) ===\n", name, lines)
	fmt.Println(strings.Repeat("-", 80))

	start := len(historyLines) - lines
	if start < 0 {
		start = 0
	}

	for i, line := range historyLines[start:] {
		if line != "" {
			fmt.Printf("%3d. %s\n", i+1, line)
		}
	}
	fmt.Println()
}

// Helper functions

func getWorkspaces(basePath string) ([]Workspace, error) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Workspace{}, nil
		}
		return nil, err
	}

	var workspaces []Workspace

	for _, entry := range entries {
		if entry.IsDir() {
			wsPath := filepath.Join(basePath, entry.Name())
			config := readConfig(filepath.Join(wsPath, configFile))

			createdTime, _ := time.Parse(time.RFC3339, config["created"])
			commandCount := 0
			fmt.Sscanf(config["commands"], "%d", &commandCount)

			workspaces = append(workspaces, Workspace{
				Name:         entry.Name(),
				CreatedAt:    createdTime,
				Path:         wsPath,
				CommandCount: commandCount,
			})
		}
	}

	// Sort by creation time (newest first)
	sort.Slice(workspaces, func(i, j int) bool {
		return workspaces[i].CreatedAt.After(workspaces[j].CreatedAt)
	})

	return workspaces, nil
}

func readConfig(configPath string) map[string]string {
	config := make(map[string]string)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return config
}

func isValidName(name string) bool {
	if len(name) == 0 || len(name) > 255 {
		return false
	}

	for _, ch := range name {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '_') {
			return false
		}
	}

	return true
}

func printUsage() {
	fmt.Println(`bashlog-mgr - Bash Command Logging Workspace Manager

Usage:
  bashlog-mgr <command> [options]

Commands:
  list              List all workspaces with statistics
  create <name>     Create a new workspace
  delete <name>     Delete a workspace (with confirmation)
  view <name>       View detailed information about a workspace
  stats             Display overall statistics across all workspaces
  history <name> [lines]  Show command history for a workspace (default: last 20 lines)
  help              Show this help message

Examples:
  bashlog-mgr list
  bashlog-mgr create my-project
  bashlog-mgr delete old-workspace
  bashlog-mgr view my-project
  bashlog-mgr stats
  bashlog-mgr history my-project 50

Workspaces are stored in: ~/.bashlog-workspaces/
`)
}
