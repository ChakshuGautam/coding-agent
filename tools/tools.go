package tools

import "google.golang.org/genai"
import "agent/tools/gitTools"

// ToolDefinition represents a tool that can be used by the agent
type ToolDefinition struct {
	Definition *genai.FunctionDeclaration
	Function   func(input *genai.FunctionCall) (string, error)
}

// GetTools returns all available tools
func GetTools() []ToolDefinition {
	return []ToolDefinition{
		{
			Definition: ReadFileDefinition,
			Function:   ReadFile,
		},
		{
			Definition: ListFilesDefinition,
			Function:   ListFiles,
		},
		{
			Definition: EditFileDefinition,
			Function:   EditFile,
		},
		{
			Definition: ListReposDefination,
			Function: ListRepos,
		},
		{
			Definition: gitTools.GitAddFileDefination,
			Function: gitTools.GitAddFile,
		},
		{
			Definition: gitTools.GitAddFileDefination,
			Function: gitTools.GitAddFile,
		},
		{
			Definition: gitTools.GitCheckoutDefination,
			Function: gitTools.GitCheckout,
		},
		{
			Definition: gitTools.GitAddRemoteAndPushDefination,
			Function: gitTools.GitAddRemoteAndPush,
		},
		{
			Definition: gitTools.GitCommitChangesDefination,
			Function: gitTools.GitCommitChanges,
		},
		{
			Definition: gitTools.GitListingRemotesDefination,
			Function: gitTools.GitListingRemotes,
		},
	}
}
