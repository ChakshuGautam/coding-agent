package tools

import (
	
	"fmt"
	"os"
	"os/exec"
	"google.golang.org/genai"

)

var AddRemoteAndPushInput = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"Name": {
			Type: genai.TypeString,
			Description: "Name of the repo where user wants to push code remotely",
		},
		"Branch" : {
			Type: genai.TypeString,
			Description: "Name of the branch where user wants to push code remotely,Defaults to current Branch",
		},
		"UserName":{
			Type: genai.TypeString,
			Description: "Username of the GitHub where user wants to connect",
		},
	},
	Required: []string{"Name"},
}

var AddRemoteAndPushDefination = &genai.FunctionDeclaration{
	Name: "addRemoteAndPush",
	Description: "Remotely connect to the specific repository and then push code to specific mentoined branch. If branch is not mentioned push code to the current branch",
	Parameters: AddRemoteAndPushInput,
}



func AddRemoteAndPush(input *genai.FunctionCall) (string, error) {

	// username := os.Getenv("GIT_USERNAME")
	username,ok := input.Args["UserName"].(string)
	if !ok || username==""{
		username=os.Getenv("GIT_USERNAME")
	}
	repoName := input.Args["Name"].(string)
	branch := input.Args["Branch"].(string)

	dir := "./"

	url := fmt.Sprintf("https://github.com/%s/%s.git", username, repoName)

	cmdCheckRemote := exec.Command("git", "remote", "get-url", "origin")
	cmdCheckRemote.Dir = dir 

	// Run the check command to see if remote exists
	_, err := cmdCheckRemote.CombinedOutput()
	if err == nil {
		// Remote exists, skip adding
		fmt.Println("Remote 'origin' already exists, skipping remote add.")
	} else {
		// Remote doesn't exist, add it
		cmdAddRemote := exec.Command("git", "remote", "add", "origin", url)
		cmdAddRemote.Dir = dir // Set the directory for the repo

		// Execute remote add command
		outputAddRemote, err := cmdAddRemote.CombinedOutput()
		if err != nil {
			return "Failed to add remote", fmt.Errorf("error adding remote: %v\noutput: %s", err, outputAddRemote)
		}
		fmt.Println("Remote added successfully")
	}

	cmdPush := exec.Command("git", "push", "-u", url, branch)
	cmdPush.Dir = dir 

	// Set up authentication using GitHub token for HTTPS
	cmdPush.Env = append(os.Environ(), "GIT_ASKPASS=echo", "GIT_USERNAME="+username, "GIT_PASSWORD="+os.Getenv("GITHUB_TOKEN"))

	// Execute the push command
	outputPush, err := cmdPush.CombinedOutput()
	if err != nil {
		return "Failed to push to GitHub", fmt.Errorf("error pushing to GitHub: %v\noutput: %s", err, outputPush)
	}

	// Return the result of both operations
	return fmt.Sprintf("Push to remote completed successfully. %s", outputPush), nil
}
