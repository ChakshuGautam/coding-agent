package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

var ListFilesInput = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"path": {
			Type:        genai.TypeString,
			Description: "Optional relative path to list files from. Defaults to current directory if not provided.",
		},
	},
}

var ListFilesDefinition = &genai.FunctionDeclaration{
	Name:        "listFiles",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
	Parameters:  ListFilesInput,
}

func ListFiles(input *genai.FunctionCall) (string, error) {
	dir := "."
	if path, ok := input.Args["path"].(string); ok && path != "" {
		dir = path
	}

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to list files: %w", err)
	}

	// Join the files with newlines for better readability
	result := ""
	for i, file := range files {
		if i > 0 {
			result += "\n"
		}
		result += file
	}

	return result, nil
}
