package gitTools

import (
	"fmt"
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

var GitFileStatusDefinition = &genai.FunctionDeclaration{
	Name:	"gitFileStatus",
	Description: "Log status of files at a given path. If no path is provided, log file status of the current directory",
	Parameters:  FileStatusInput,
}

func GitFileStatus(input *genai.FunctionCall) (string,error){
	cmd := execCommand("git", "status")
	output, err := cmd.CombinedOutput()
	if err!=nil{
		return "",fmt.Errorf("error : %v",err)
	}

	return string(output),nil
}