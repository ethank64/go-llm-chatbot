package gemini

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func getCurrentEST() string {
	loc, _ := time.LoadLocation(("America/New_York"))
	return time.Now().In(loc).Format("Mon Jan 2 15:04:05 MST 2006")
}

// Takes the path to an image and returns the analysis
func analyzeImage(path string) string {
	return ""
}

// Takes in a prompt and generates an image
// Saves it locally and returns the file path
func generateImage(prompt string) string {
	return ""
}

func openYoutubeMusic() {
	cmd := exec.Command("open", "-a", "Safari", "https://music.youtube.com")
	cmd.Run()
}

// Executes a fuzzy search and returns the file path of the matching file
func findFile(searchQuery string) string {
	// Run a find + grep pipeline
	cmd := exec.Command("bash", "-c", fmt.Sprintf(
		`find ~/Desktop -type f 2>/dev/null | grep -i "%s" | head -n 1`, searchQuery,
	))

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Error while finding while"
	}

	filePath := strings.TrimSpace(out.String())
	if filePath == "" {
		return "No match found."
	} else {
		return "Found: " + filePath
	}
}

func emptyTrash() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	trashPath := filepath.Join(home, ".Trash")
	entries, err := os.ReadDir(trashPath)
	if err != nil {
		return fmt.Errorf("could not read Trash direcotry: %w", err)
	}

	for _, entry := range entries {
		itemPath := filepath.Join(trashPath, entry.Name())
		err := os.RemoveAll(itemPath)
		if err != nil {
			fmt.Printf("Failed to delete %s: %v\n", itemPath, err)
		}
	}

	// No errors; return nothing
	return nil
}

// Tool ideas:
// Search wikipedia for information
// Report generation
// Empty trash
