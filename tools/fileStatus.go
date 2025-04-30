package tools

import (
	"fmt"
	"os/exec"

	"google.golang.org/genai"
)

var FileStatusInput = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"path": {
			Type:        genai.TypeString,
			Description: "Optional relative path to log file status from. Defaults to current directory if not provided.",
		},
	},
}

var FileStatusDefination = &genai.FunctionDeclaration{
	Name:	"fileStatus",
	Description: "Log status of files at a given path. If no path is provided, log file status of the current directory",
	Parameters:  FileStatusInput,
}

func FileStatus(input *genai.FunctionCall) (string,error){
	cmd := exec.Command("git", "status")
	output, err := cmd.CombinedOutput()
	if err!=nil{
		return "",fmt.Errorf("error : %v",err)
	}

	return string(output),nil
}