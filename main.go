package main

import (
	"log"

	"github.com/ethank64/go-llm-chatbot/gemini"
)

func main() {
	geminiService, err := gemini.NewService()
	if err != nil {
		log.Fatalf("Failed to create Gemini service: %v: ", err)
	}

	geminiService.Run()
}
