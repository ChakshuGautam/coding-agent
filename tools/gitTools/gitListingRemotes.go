package gitTools

import (
	"fmt"

	"google.golang.org/genai"
)


var GitListingRemotesDefination = &genai.FunctionDeclaration{
	Name:        "listingRemotes",
	Description: "List down all the remotes",
}

func GitListingRemotes(input *genai.FunctionCall) (string, error) {
	cmd := execCommand("git", "remote", "-v")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}

	return string(output), nil
}
