package gemini

import "google.golang.org/genai"

func GetFunctionSchemas() []*genai.Tool {
	return []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        "get_current_time",
					Description: "Gets the current time in EST",
				},
				{
					Name:        "analyze_image",
					Description: "Analyzes an image from a local file path and returns a detailed description",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path": {
								Type:        genai.TypeString,
								Description: "The local file path to the image to analyze",
							},
						},
						Required: []string{"path"},
					},
				},
				{
					Name:        "generate_image",
					Description: "Generates an image from a text prompt",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"prompt": {
								Type:        genai.TypeString,
								Description: "The text description of the image to generate",
							},
						},
						Required: []string{"prompt"},
					},
				},
				{
					Name:        "open_youtube_music",
					Description: "Opens YouTube Music in the Safari web browser",
				},
				{
					Name:        "open_github",
					Description: "Opens GitHub in the Safari web browser",
				},
				{
					Name:        "find_file",
					Description: "Executes a fuzzy search to find its file path",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"searchPrompt": {
								Type:        genai.TypeString,
								Description: "The string that will be used in the fuzzy search for the file path.",
							},
						},
						Required: []string{"searchPrompt"},
					},
				},
				{
					Name:        "empty_trash",
					Description: "Deletes all of the files in the recycle bin.",
				},
			},
		},
	}
}
