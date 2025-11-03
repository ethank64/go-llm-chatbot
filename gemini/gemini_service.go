package gemini

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

const ModelName = "gemini-2.5-flash"

type GeminiService struct {
	client            *genai.Client
	conversation      []*genai.Content
	systemInstruction *genai.Content
}

// Method to create a new instance of the service
func NewService() (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	// If no errors from the client generation,
	// return the address of a new service instance
	return &GeminiService{client: client}, nil
}

// Methods for GeminiService

func (gs *GeminiService) Run() {
	greetUser()

	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("User: ")

		if !reader.Scan() {
			break
		}

		prompt := reader.Text()
		if prompt == "quit" {
			printGeminiMessage("Goodbye!")
			return
		}

		reply, err := gs.Ask(prompt)
		if err != nil {
			log.Println("Error: ", err)
			continue
		}

		if reply != "" {
			printGeminiMessage(reply)
		}
	}
}

func (gs *GeminiService) Ask(prompt string) (string, error) {
	ctx := context.Background()

	// Add user prompt to converation history
	gs.conversation = append(gs.conversation, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: prompt}},
	})

	// Get tool configs
	tools := GetFunctionSchemas()
	config := genai.GenerateContentConfig{
		Tools: tools,
	}

	// Add system instruction if set
	if gs.systemInstruction != nil {
		config.SystemInstruction = gs.systemInstruction
	}

	// Model invocation
	result, err := gs.client.Models.GenerateContent(ctx, ModelName, gs.conversation, &config)
	if err != nil {
		return "", err
	}

	reply := ""

	// Get the response
	candidate := result.Candidates[0]

	// If we have a function call, execute it
	if len(candidate.Content.Parts) > 0 && candidate.Content.Parts[0].FunctionCall != nil {
		fc := candidate.Content.Parts[0].FunctionCall
		gs.HandleFunctionCall(ctx, fc)
	} else {
		// Else, we have text response so print it out
		reply = candidate.Content.Parts[0].Text

		gs.conversation = append(gs.conversation, &genai.Content{
			Role:  "model",
			Parts: []*genai.Part{{Text: reply}},
		})
	}

	return reply, nil
}

// SetSystemInstruction sets the system instruction for the Gemini service
func (gs *GeminiService) SetSystemInstruction(instruction string) {
	gs.systemInstruction = &genai.Content{
		Parts: []*genai.Part{{Text: instruction}},
	}
}

// GetSystemInstruction returns the current system instruction text
func (gs *GeminiService) GetSystemInstruction() string {
	if gs.systemInstruction != nil && len(gs.systemInstruction.Parts) > 0 {
		return gs.systemInstruction.Parts[0].Text
	}
	return ""
}

func (gs *GeminiService) HandleFunctionCall(ctx context.Context, fc *genai.FunctionCall) {
	switch fc.Name {
	case "get_current_time":
		result := map[string]any{
			"time": getCurrentEST(),
		}
		gs.SendFunctionResult(ctx, fc, result)

	case "open_youtube_music":
		openYoutubeMusic()

		result := map[string]any{
			"successful": true,
		}

		gs.SendFunctionResult(ctx, fc, result)

	case "open_github":
		openGithub()

		result := map[string]any{
			"successful": true,
		}

		gs.SendFunctionResult(ctx, fc, result)

	case "find_file":
		// Extract search query
		searchPrompt, _ := fc.Args["searchPrompt"].(string)

		path := findFile(searchPrompt)

		result := map[string]any{
			"searchPrompt": searchPrompt,
			"path":         path,
		}

		gs.SendFunctionResult(ctx, fc, result)

	case "empty_trash":
		err := emptyTrash()

		successful := true
		errorMessage := ""

		if err != nil {
			successful = false
			errorMessage = err.Error()
		}

		result := map[string]any{
			"successful":   successful,
			"errorMessage": errorMessage,
		}

		gs.SendFunctionResult(ctx, fc, result)

	case "analyze_image":
		// Extract image path
		imagePath, _ := fc.Args["path"].(string)

		analysis := analyzeImage(imagePath)

		result := map[string]any{
			"path":     imagePath,
			"analysis": analysis,
		}

		gs.SendFunctionResult(ctx, fc, result)

	case "generate_image":
		// Extract prompt
		prompt, _ := fc.Args["prompt"].(string)

		filePath := generateImage(prompt)

		result := map[string]any{
			"prompt":   prompt,
			"filePath": filePath,
		}

		gs.SendFunctionResult(ctx, fc, result)

	default:
		fmt.Println("Unknown function:", fc.Name)
	}
}

func (gs *GeminiService) SendFunctionResult(ctx context.Context, fc *genai.FunctionCall, result map[string]any) {
	// Append the function call to the history
	gs.conversation = append(gs.conversation, &genai.Content{
		Role: "model",
		Parts: []*genai.Part{
			{FunctionCall: fc},
		},
	})

	// Append function response
	gs.conversation = append(gs.conversation, &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{FunctionResponse: &genai.FunctionResponse{
				Name:     fc.Name,
				Response: result,
			}},
		},
	})

	// Prepare config with system instruction if set
	var config *genai.GenerateContentConfig
	if gs.systemInstruction != nil {
		config = &genai.GenerateContentConfig{
			SystemInstruction: gs.systemInstruction,
		}
	}

	// Send full conversation to Gemini
	resp, err := gs.client.Models.GenerateContent(ctx, ModelName, gs.conversation, config)
	if err != nil {
		log.Fatal(err)
	}

	// Check if response has candidates and parts
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		// Get the text response if available
		textResponse := ""
		for _, part := range resp.Candidates[0].Content.Parts {
			if part.Text != "" {
				textResponse = part.Text
				break
			}
		}

		// Only append and print if we have a text response
		if textResponse != "" {
			gs.conversation = append(gs.conversation, &genai.Content{
				Role:  "model",
				Parts: []*genai.Part{{Text: textResponse}},
			})

			printGeminiMessage(textResponse)
		}
	}
}
