package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

func main() {
	// Initial AI message
	printGeminiMessage("Hello good sir! How can I be of assistance today?")

	// Get user input
	reader := bufio.NewScanner(os.Stdin)

	for reader.Scan() {
		fmt.Print("User: ")

		// Read user input
		prompt := reader.Text()

		// Handle quit input
		if prompt == "quit" {
			printGeminiMessage("Goodbye!")
			break
		}

		// Print Gemini's response
		promptGemini(prompt)
	}

}

func promptGemini(prompt string) {
	ctx := context.Background()
	// The client gets the API key from the environment variable `GEMINI_API_KEY`.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	printGeminiMessage(result.Text())
}

func printGeminiMessage(msg string) {
	fmt.Println("Chap GPT: " + msg)
}
