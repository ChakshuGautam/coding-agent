package tools

import (
	"fmt"
	"os"

	"google.golang.org/genai"
)

var ReadFileInput = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"path": {Type: genai.TypeString},
	},
	Required: []string{"path"},
}

var ReadFileDefinition = &genai.FunctionDeclaration{
	Name:        "readFile",
	Description: "Read a file from the local filesystem",
	Parameters:  ReadFileInput,
}

func ReadFile(input *genai.FunctionCall) (string, error) {
	content, err := os.ReadFile(input.Args["path"].(string))
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s': %w", input.Args["path"], err)
	}

	return string(content), nil
}
