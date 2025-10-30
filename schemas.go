package main

import "google.golang.org/genai"

func GetFunctionSchemas() []*genai.Tool {
	return []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        "get_current_time",
					Description: "Gets the current time in EST",
				},
			},
		},
	}
}
