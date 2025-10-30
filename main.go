package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/genai"
)

var conversation []*genai.Content

func main() {
	// Initial AI message
	printGeminiMessage("Hello good sir! How can I be of assistance today?")

	// Get user input
	fmt.Print("User: ")
	reader := bufio.NewScanner(os.Stdin)

	for reader.Scan() {
		// Read user input
		prompt := reader.Text()

		// Handle quit input
		if prompt == "quit" {
			printGeminiMessage("Goodbye!")
			break
		}

		// Print Gemini's response
		promptGemini(prompt)

		fmt.Print("User: ")
	}

}

func promptGemini(prompt string) {
	ctx := context.Background()
	// The client gets the API key from the environment variable `GEMINI_API_KEY`.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get tool configs
	tools := GetFunctionSchemas()
	config := genai.GenerateContentConfig{
		Tools: tools,
	}

	// Append user message to history
	conversation = append(conversation, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: prompt}},
	})

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		conversation,
		&config,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get the response
	candidate := result.Candidates[0]

	// If we have a function call, execute it
	if len(candidate.Content.Parts) > 0 && candidate.Content.Parts[0].FunctionCall != nil {
		fc := candidate.Content.Parts[0].FunctionCall
		handleFunctionCall(ctx, client, fc)
	} else {
		// Else, we have text response so print it out
		reply := candidate.Content.Parts[0].Text
		printGeminiMessage(reply)

		conversation = append(conversation, &genai.Content{
			Role:  "model",
			Parts: []*genai.Part{{Text: reply}},
		})
	}
}

func printGeminiMessage(msg string) {
	fmt.Println("Chap GPT: " + msg)
}

func handleFunctionCall(ctx context.Context, client *genai.Client, fc *genai.FunctionCall) {
	if fc.Name == "get_current_time" {
		result := map[string]any{
			"time": getCurrentEST(),
		}
		sendFunctionResult(ctx, client, fc, result)
	} else {
		fmt.Println("Unknown function: ", fc.Name)
	}
}

func getCurrentEST() string {
	loc, _ := time.LoadLocation(("America/New_York"))
	return time.Now().In(loc).Format("Mon Jan 2 15:04:05 MST 2006")
}

func sendFunctionResult(ctx context.Context, client *genai.Client, fc *genai.FunctionCall, result map[string]any) {
	// Append the function call to the history
	conversation = append(conversation, &genai.Content{
		Role: "model",
		Parts: []*genai.Part{
			{FunctionCall: fc},
		},
	})

	// Append function response
	conversation = append(conversation, &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{FunctionResponse: &genai.FunctionResponse{
				Name:     fc.Name,
				Response: result,
			}},
		},
	})

	// Send full conversation to Gemini
	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", conversation, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Append Gemini's text reply
	conversation = append(conversation, &genai.Content{
		Role:  "model",
		Parts: []*genai.Part{{Text: resp.Candidates[0].Content.Parts[0].Text}},
	})

	printGeminiMessage(resp.Candidates[0].Content.Parts[0].Text)
}
