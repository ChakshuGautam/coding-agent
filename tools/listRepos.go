package tools

import (
	
	"os"

	"context"

	"github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"

	"google.golang.org/genai"
)

var ListResposInput = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"path": {
			Type: genai.TypeString,
			Description: "Optional relative path to list repos from. Defaults to your own repository if not provided",
		},
	},
}

var ListReposDefination = &genai.FunctionDeclaration{
	Name: "listRepos",
	Description: "List all the repos of user. If no user specified, list files of own user.",
	Parameters: ListResposInput,
}



func ListRepos(input *genai.FunctionCall) (string,error){
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
    )
    tc := oauth2.NewClient(ctx, ts)

    client := github.NewClient(tc)
    repos, _, _ := client.Repositories.List(ctx, "", nil)
	result :=""
    for i, repo := range repos {
		if( i >0 ){
			result +="\n"
		}
		result += *repo.Name
    }

	return result, nil
}
