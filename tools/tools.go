package tools

import "google.golang.org/genai"

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
			Definition: FileStatusDefination,
			Function: FileStatus,
		},
		{
			Definition: AddFileDefination,
			Function: AddFile,
		},
		{
			Definition: CheckoutDefination,
			Function: Checkout,
		},
		{
			Definition: AddRemoteAndPushDefination,
			Function: AddRemoteAndPush,
		},
		{
			Definition: CommitChangesDefination,
			Function: CommitChanges,
		},
		{
			Definition: ListingRemotesDefination,
			Function: ListingRemotes,
		},
	}
}
