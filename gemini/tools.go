package gemini

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/genai"
)

func getCurrentEST() string {
	loc, _ := time.LoadLocation(("America/New_York"))
	return time.Now().In(loc).Format("Mon Jan 2 15:04:05 MST 2006")
}

// Takes the path to an image and returns the analysis
func analyzeImage(path string) string {
	// Read the image file
	imageData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading image file: %v", err)
	}

	// Determine MIME type based on file extension
	mimeType := "image/jpeg"
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	}

	// Create Gemini client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return fmt.Sprintf("Error creating Gemini client: %v", err)
	}

	// Create content with image and analysis prompt
	content := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					InlineData: &genai.Blob{
						MIMEType: mimeType,
						Data:     imageData,
					},
				},
				{
					Text: "Analyze this image in detail. Describe what you see, including objects, colors, composition, and any notable features.",
				},
			},
		},
	}

	// Call Gemini API
	result, err := client.Models.GenerateContent(ctx, ModelName, content, nil)
	if err != nil {
		return fmt.Sprintf("Error analyzing image: %v", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "No analysis returned from Gemini"
	}

	return result.Candidates[0].Content.Parts[0].Text
}

// Takes in a prompt and generates an image
// Saves it locally and returns the file path
func generateImage(prompt string) string {
	// Create Gemini client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return fmt.Sprintf("Error creating Gemini client: %v", err)
	}

	// Use a model that supports image generation
	imageGenModel := "gemini-2.5-flash-image"

	// Create content with image generation prompt using genai.Text helper
	content := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: prompt,
				},
			},
		},
	}

	// Configure to request image generation
	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"IMAGE"},
		ImageConfig: &genai.ImageConfig{
			AspectRatio: "1:1",
		},
	}

	// Call Gemini API for image generation
	result, err := client.Models.GenerateContent(ctx, imageGenModel, content, config)
	if err != nil {
		// Log error for debugging
		log.Printf("Error generating image with model %s: %v", imageGenModel, err)
		return fmt.Sprintf("Error generating image with model %s: %v", imageGenModel, err)
	}

	// Check if the response contains image data
	var imageData []byte
	var mimeType string

	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		for _, part := range result.Candidates[0].Content.Parts {
			if part.InlineData != nil {
				// Data is already []byte
				imageData = part.InlineData.Data
				mimeType = part.InlineData.MIMEType
				break
			}
		}
	}

	// If no image data in response, return error
	if len(imageData) == 0 {
		return "Image generation returned no image data. Please check that the model supports image generation and your API key has the necessary permissions."
	}

	// Create generated_images directory if it doesn't exist
	outputDir := "generated_images"
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Sprintf("Error creating output directory: %v", err)
	}

	// Determine file extension from MIME type
	ext := ".jpg"
	if mimeType == "image/png" {
		ext = ".png"
	} else if mimeType == "image/gif" {
		ext = ".gif"
	} else if mimeType == "image/webp" {
		ext = ".webp"
	}

	// Generate unique filename with timestamp
	filename := fmt.Sprintf("generated_%d%s", time.Now().Unix(), ext)
	filePath := filepath.Join(outputDir, filename)

	// Save the image
	err = os.WriteFile(filePath, imageData, 0644)
	if err != nil {
		return fmt.Sprintf("Error saving image: %v", err)
	}

	return filePath
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
