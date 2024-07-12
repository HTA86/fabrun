package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

var (
	listFlag           bool
	versionFlag        bool
	updateCommandsFlag bool
	version            = "0.5.0"
	commandsURL        = "https://api.github.com/repos/HTA86/fabrun/contents/commands"
	totalFiles         int
	progress           *progressbar.ProgressBar
)

func init() {
	flag.BoolVar(&listFlag, "list", false, "")
	flag.BoolVar(&listFlag, "l", false, "List all available commands")
	flag.BoolVar(&versionFlag, "version", false, "")
	flag.BoolVar(&versionFlag, "v", false, "Show the program's version")
	flag.BoolVar(&updateCommandsFlag, "update-commands", false, "Update commands from GitHub")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if versionFlag {
		printVersion()
		return
	}

	if listFlag {
		ListCommands()
		return
	}

	if updateCommandsFlag {
		updateCommands()
		return
	}

	if len(flag.Args()) < 1 {
		flag.Usage()
		return
	}

	commandID := flag.Arg(0)
	commandPath := GetCommandFilePath(commandID)

	cmd, err := ReadCommandFromFile(commandPath)
	if err != nil {
		LogError("Error reading command:", err)
		return
	}

	ExecuteCommand(cmd)
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: fabrun [options] <command_name>\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  -l, --list             List all available commands\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  -v, --version          Show the program's version\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  -h, --help             Show this help message\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  --update-commands      Update commands from GitHub\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\nExample usage:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun --list               # List all available commands\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun -l                   # List all available commands (short flag)\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun git_change_2weeks    # Run the git_change_2weeks command\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun --version            # Show the program's version\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun -v                   # Show the program's version (short flag)\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  fabrun --update-commands    # Update commands from GitHub\n")
}

func printVersion() {
	fmt.Println("fabrun version", version)
}

func updateCommands() {
	fmt.Println("Updating commands from GitHub...")

	// Create the destination directory if it doesn't exist
	destDir := filepath.Join(os.Getenv("HOME"), ".config", "fabrun", "commands")
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		LogError("Failed to create commands directory:", err)
		return
	}

	// Count the total number of files and directories
	totalFiles = 0
	if err := countFiles(commandsURL); err != nil {
		LogError("Failed to count files:", err)
		return
	}

	// Initialize the progress bar
	progress = progressbar.NewOptions(totalFiles,
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetRenderBlankState(true),
	)

	// Fetch the list of files and directories from the commands directory in the repository
	if err := fetchAndDownloadFiles(commandsURL, destDir); err != nil {
		LogError("Failed to update commands:", err)
		return
	}

	fmt.Println("\nCommands updated successfully!")
}

func countFiles(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch commands: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch commands: unexpected status code %d", resp.StatusCode)
	}

	var files []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	for _, file := range files {
		if file.Type == "dir" {
			if err := countFiles(file.URL); err != nil {
				return err
			}
		} else if file.Type == "file" {
			totalFiles++
		}
	}

	return nil
}

func fetchAndDownloadFiles(url string, destDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch commands: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch commands: unexpected status code %d", resp.StatusCode)
	}

	var files []struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		DownloadURL string `json:"download_url"`
		Type        string `json:"type"`
		URL         string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	for _, file := range files {
		if file.Type == "dir" {
			// Recursively fetch and download files in the directory
			newDestDir := filepath.Join(destDir, file.Name)
			if err := os.MkdirAll(newDestDir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", newDestDir, err)
			}
			if err := fetchAndDownloadFiles(file.URL, newDestDir); err != nil {
				return err
			}
		} else if file.Type == "file" {
			if file.DownloadURL == "" {
				LogError("No download URL found for file:", fmt.Errorf("%s", file.Name))
				continue
			}

			fileResp, err := http.Get(file.DownloadURL)
			if err != nil {
				return fmt.Errorf("failed to download file %s: %w", file.Name, err)
			}
			defer fileResp.Body.Close()

			outFile, err := os.Create(filepath.Join(destDir, file.Name))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", file.Name, err)
			}

			if _, err := io.Copy(outFile, fileResp.Body); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file %s: %w", file.Name, err)
			}

			outFile.Close()

			// Update the progress bar
			progress.Add(1)
		}
	}

	return nil
}
