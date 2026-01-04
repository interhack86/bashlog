package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Config holds the configuration for bashlog
type Config struct {
	Timezone  string
	Date      string
	Time      string
	LogDir    string
	RCFile    string
	SessionID string
}

func main() {
	// Define flags
	tzFlag := flag.String("tz", "", "Timezone for logging (e.g., UTC, America/New_York)")
	dateFlag := flag.String("date", "", "Date for logging (YYYY-MM-DD format)")
	timeFlag := flag.String("time", "", "Time for logging (HH:MM:SS format)")

	flag.Parse()

	// Setup configuration
	config, err := setupConfig(*tzFlag, *dateFlag, *timeFlag)
	if err != nil {
		log.Fatalf("Failed to setup configuration: %v", err)
	}

	// Show session information
	showSessionInfo(config)

	// Create RC file
	if err := createRCFile(config); err != nil {
		log.Fatalf("Failed to create RC file: %v", err)
	}

	// Run shell with logging
	if err := runShell(config); err != nil {
		log.Fatalf("Failed to run shell: %v", err)
	}
}

// setupConfig initializes the configuration for bashlog
func setupConfig(tz, date, timeStr string) (*Config, error) {
	config := &Config{
		Timezone: tz,
		Date:     date,
		Time:     timeStr,
	}

	// Set default timezone if not provided
	if config.Timezone == "" {
		config.Timezone = "UTC"
	}

	// Set current date if not provided
	if config.Date == "" {
		loc, err := time.LoadLocation(config.Timezone)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone: %w", err)
		}
		config.Date = time.Now().In(loc).Format("2006-01-02")
	}

	// Set current time if not provided
	if config.Time == "" {
		loc, err := time.LoadLocation(config.Timezone)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone: %w", err)
		}
		config.Time = time.Now().In(loc).Format("15:04:05")
	}

	// Generate session ID
	config.SessionID = fmt.Sprintf("session_%s_%s", config.Date, config.Time)

	// Setup log directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	config.LogDir = filepath.Join(homeDir, ".bashlog", "logs", config.Date)

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Set RC file path
	config.RCFile = filepath.Join(homeDir, ".bashlog", "bashlog.rc")

	return config, nil
}

// showSessionInfo displays information about the current session
func showSessionInfo(config *Config) {
	fmt.Println("====================================")
	fmt.Println("         Bashlog Session Info")
	fmt.Println("====================================")
	fmt.Printf("Timezone:    %s\n", config.Timezone)
	fmt.Printf("Date:        %s\n", config.Date)
	fmt.Printf("Time:        %s\n", config.Time)
	fmt.Printf("Session ID:  %s\n", config.SessionID)
	fmt.Printf("Log Dir:     %s\n", config.LogDir)
	fmt.Printf("RC File:     %s\n", config.RCFile)
	fmt.Println("====================================")
}

// createRCFile creates the bashlog RC configuration file
func createRCFile(config *Config) error {
	rcDir := filepath.Dir(config.RCFile)

	// Create RC directory if it doesn't exist
	if err := os.MkdirAll(rcDir, 0755); err != nil {
		return fmt.Errorf("failed to create RC directory: %w", err)
	}

	// RC file content
	content := fmt.Sprintf(`# Bashlog RC Configuration
# Generated at %s

# Timezone setting
export BASHLOG_TIMEZONE="%s"

# Logging directory
export BASHLOG_LOG_DIR="%s"

# Session ID
export BASHLOG_SESSION_ID="%s"

# Enable logging
export BASHLOG_ENABLED=1

# Log history
export HISTFILE="%s"
export HISTSIZE=10000
export HISTFILESIZE=10000

# Log command execution
PROMPT_COMMAND="history -a; $PROMPT_COMMAND"
`, time.Now().UTC().Format("2006-01-02 15:04:05"), config.Timezone, config.LogDir, config.SessionID, filepath.Join(config.LogDir, ".bash_history"))

	// Write RC file
	if err := os.WriteFile(config.RCFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write RC file: %w", err)
	}

	log.Printf("RC file created at: %s", config.RCFile)
	return nil
}

// runShell executes an interactive shell with logging enabled
func runShell(config *Config) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	// Create log file path
	logFile := filepath.Join(config.LogDir, fmt.Sprintf("session_%s.log", config.Time))

	// Setup command
	cmd := exec.Command(shell, "-i")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment variables for the shell
	env := os.Environ()
	env = append(env, fmt.Sprintf("BASHLOG_SESSION_ID=%s", config.SessionID))
	env = append(env, fmt.Sprintf("BASHLOG_LOG_FILE=%s", logFile))
	env = append(env, fmt.Sprintf("BASHLOG_TIMEZONE=%s", config.Timezone))
	cmd.Env = env

	log.Printf("Starting shell: %s", shell)
	log.Printf("Logging to: %s", logFile)

	// Execute shell
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run shell: %w", err)
	}

	return nil
}
