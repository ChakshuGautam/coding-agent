package tools

import (
	"fmt"
	"os/exec"
	"os"
	"strings"
	"google.golang.org/genai"
)

var AddFileInput = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"name": {
			Type:        genai.TypeString,
			Description: "Optional name to add only specific file to staging area. If no file name is provided add all the files to staging area. ",
		},
	},
}

var AddFileDefination = &genai.FunctionDeclaration{
	Name:        "addFile",
	Description: "Adds the files in the staging area. If no name is provided, add all files to staging area.",
	Parameters:  AddFileInput,
}

// EnsureGitignore checks if .gitignore exists, and ensures ".env" is listed
func EnsureGitignore() error {
	const gitignorePath = ".gitignore"
	const entry = ".env"

	// Check if .gitignore exists
	content := ""
	if _, err := os.Stat(gitignorePath); err == nil {
		// File exists — read it
		data, err := os.ReadFile(gitignorePath)
		if err != nil {
			return fmt.Errorf("failed to read .gitignore: %v", err)
		}
		content = string(data)
	} else if os.IsNotExist(err) {
		// File doesn't exist — will create it
		content = ""
	} else {
		return fmt.Errorf("error checking .gitignore: %v", err)
	}

	// Check if ".env" is already listed
	if !strings.Contains(content, entry) {
		f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open .gitignore for writing: %v", err)
		}
		defer f.Close()

		// Add newline if needed
		if len(content) > 0 && content[len(content)-1] != '\n' {
			_, _ = f.WriteString("\n")
		}

		_, err = f.WriteString(entry + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to .gitignore: %v", err)
		}
	}

	return nil
}

func AddFile(input *genai.FunctionCall) (string, error) {

	if err := EnsureGitignore(); err != nil {
    	return "", fmt.Errorf("could not verify .gitignore: %v", err)
	}

	name, ok := input.Args["name"].(string)
	if !ok || name == "" {
		name = "."
	}

	cmd := exec.Command("git", "add", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}

	return string(output), nil
}
