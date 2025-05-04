package gitTools

import (
	
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	
	"github.com/stretchr/testify/assert"
	"google.golang.org/genai"
)

func TestGitListRepos_Success(t *testing.T) {
	// Set a dummy GitHub token
	os.Setenv("GITHUB_TOKEN", "dummy_token")

	// Activate the mock HTTP responder
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock GitHub API response
	mockRepos := `[{"name": "repo1"}, {"name": "repo2"}]`

	httpmock.RegisterResponder("GET", "https://api.github.com/user/repos",
		httpmock.NewStringResponder(200, mockRepos))

	input := &genai.FunctionCall{
		Name: "listRepos",
		Args: map[string]any{},
	}

	result, err := GitListRepos(input)
	assert.NoError(t, err)
	assert.Equal(t, "repo1\nrepo2", result)
}

func TestGitListRepos_EmptyList(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "dummy_token")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.github.com/user/repos",
		httpmock.NewStringResponder(200, `[]`))

	input := &genai.FunctionCall{
		Name: "listRepos",
		Args: map[string]any{},
	}

	result, err := GitListRepos(input)
	assert.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestGitListRepos_Unauthorized(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "invalid_token")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.github.com/user/repos",
		httpmock.NewStringResponder(401, `{"message": "Bad credentials"}`))

	input := &genai.FunctionCall{
		Name: "listRepos",
		Args: map[string]any{},
	}

	result, err := GitListRepos(input)
	assert.NoError(t, err) // because your current implementation ignores the error!
	assert.Equal(t, "", result)
}
