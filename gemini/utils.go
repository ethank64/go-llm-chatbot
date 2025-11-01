package gemini

import (
	"fmt"
)

func printGeminiMessage(msg string) {
	fmt.Println("Chap GPT: " + msg)
}

func greetUser() {
	printGeminiMessage("Hello good sir! How can I be of assistance today?")
}
