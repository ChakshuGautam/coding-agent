package tools

import (
	
	"fmt"
	"os/exec"
	"os"
	"context"
	"google.golang.org/genai"
)


var CommintChangesInput = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"message": {
			Type:        genai.TypeString,
			Description: "message for committing changes",
		},
	},
}

var CommitChangesDefination = &genai.FunctionDeclaration{
	Name:	"commitChanges",
	Description: "Commits all the staged files with relevant message. If no message is provided try to figure what changes are done and based on that provide message to commit",
	Parameters:  CommintChangesInput,
}

func CommitChanges(input *genai.FunctionCall) (string, error) {

	message, ok := input.Args["message"].(string)

	if !ok || message == "" {
		// Step 1: Get staged diff
		diffCmd := exec.Command("git", "diff", "--cached")
		diffOutput, err := diffCmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to get diff: %v", err)
		}
		diffText := string(diffOutput)
		if diffText == "" {
			return "", fmt.Errorf("no staged changes found to commit")
		}

		// fmt.Print(functionCall,"is the function call output")
		resp, err := GenerateCommitMessage(diffText)
		if err != nil {
			return "", fmt.Errorf("AI generation failed: %v", err)
		}

		

		message = resp
	}

	cmd := exec.Command("git", "commit", "-m", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error committing: %v\n%s", err, string(output))
	}

	return string(output), nil
}


func GenerateCommitMessage(diff string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		// Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf("Generate a concise Git commit message for the following code changes:\n\n%s", diff)

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash-preview-04-17", []*genai.Content{
		{Parts: []*genai.Part{{Text: prompt}}},
	}, nil)

	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}






