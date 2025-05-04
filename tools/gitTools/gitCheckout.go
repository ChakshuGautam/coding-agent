package gitTools

import (
	"fmt"
	"os/exec"
	"strings"
	"google.golang.org/genai"
)

var execCommand = exec.Command

var CheckoutInput = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"name": {
			Type:        genai.TypeString,
			Description: "The branch name or the new feature name where git has to checkout",
		},
	},
	Required: []string{"name"},
}

var GitCheckoutDefination = &genai.FunctionDeclaration{
	Name:        "checkout",
	Description: "Creates a new branch if no branch exists with given name or Switch to a Branch with given name. If no name is provided it will give an error",
	Parameters:  CheckoutInput,
}

func ExistingBranches() (string, error){
	cmd := execCommand("git", "branch")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}
	return string(output),nil
}

func GitCheckout(input *genai.FunctionCall) (string, error) {

	exisiting_branches,err := ExistingBranches()
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}

	name, ok := input.Args["name"].(string)
	if !ok || name == "" {
		return "Error! Check the name provided", fmt.Errorf("error : %v", ok)
	}

	branches := strings.Split(exisiting_branches,"\n")

	var cmd *exec.Cmd

	for _,branch := range branches{
		branch = strings.TrimSpace(strings.TrimPrefix(branch,"*"))
		if branch == name{
			cmd = execCommand("git", "checkout", name)
		}else{
			cmd = execCommand("git", "checkout","-b", name)
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}

	return string(output), nil
}
