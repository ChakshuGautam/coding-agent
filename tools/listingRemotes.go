package tools

import (
	"fmt"
	"os/exec"

	"google.golang.org/genai"
)


var ListingRemotesDefination = &genai.FunctionDeclaration{
	Name:        "listingRemotes",
	Description: "List down all the remotes",
}

func ListingRemotes(input *genai.FunctionCall) (string, error) {
	cmd := exec.Command("git", "remote", "-v")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}

	return string(output), nil
}
