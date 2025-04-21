package tools

import (
	"fmt"
	"os"
	"path"
	"strings"

	"google.golang.org/genai"
)

var EditFileInput = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"path": {
			Type:        genai.TypeString,
			Description: "The path to the file",
		},
		"old_str": {
			Type:        genai.TypeString,
			Description: "Text to search for - must match exactly and must only have one match exactly",
		},
		"new_str": {
			Type:        genai.TypeString,
			Description: "Text to replace old_str with",
		},
	},
	Required: []string{"path", "old_str", "new_str"},
}

var EditFileDefinition = &genai.FunctionDeclaration{
	Name: "editFile",
	Description: `Make edits to a text file.

Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str' MUST be different from each other.

If the file specified with path doesn't exist, it will be created.`,
	Parameters: EditFileInput,
}

func EditFile(input *genai.FunctionCall) (string, error) {
	path, ok := input.Args["path"].(string)
	if !ok || path == "" {
		return "", fmt.Errorf("path is required")
	}

	oldStr, ok := input.Args["old_str"].(string)
	if !ok {
		return "", fmt.Errorf("old_str is required")
	}

	newStr, ok := input.Args["new_str"].(string)
	if !ok {
		return "", fmt.Errorf("new_str is required")
	}

	if oldStr == newStr {
		return "", fmt.Errorf("old_str and new_str must be different")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) && oldStr == "" {
			return createNewFile(path, newStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.Replace(oldContent, oldStr, newStr, -1)

	if oldContent == newContent && oldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	err = os.WriteFile(path, []byte(newContent), 0644)
	if err != nil {
		return "", err
	}

	return "File edited successfully", nil
}

func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	return fmt.Sprintf("Successfully created file %s", filePath), nil
}
