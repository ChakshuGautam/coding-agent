package gitTools

import (
	"os"
	"os/exec"
	"testing"

	"google.golang.org/genai"
)

func TestMockedGitRemoteCommand(t *testing.T) {
	if os.Getenv("MOCK_GIT_REMOTE") != "1" {
		return
	}

	if os.Getenv("MOCK_ERROR") == "true" {
		os.Stderr.WriteString("simulated error")
		os.Exit(1)
	}

	os.Stdout.WriteString(os.Getenv("MOCK_STDOUT"))
	os.Exit(0)
}

func mockGitRemoteCommand(stdout string, simulateError bool) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestMockedGitRemoteCommand", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(os.Environ(),
			"MOCK_GIT_REMOTE=1",
			"MOCK_STDOUT="+stdout,
			"MOCK_ERROR="+boolToString(simulateError),
		)
		return cmd
	}
}

func TestGitListingRemotes_Success(t *testing.T) {
	mockOutput := `origin  https://github.com/user/repo.git (fetch)
origin  https://github.com/user/repo.git (push)`

	execCommand = mockGitRemoteCommand(mockOutput, false)
	defer func() { execCommand = exec.Command }()

	input := &genai.FunctionCall{Args: map[string]any{}}

	output, err := GitListingRemotes(input)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output != mockOutput {
		t.Errorf("Unexpected output.\nExpected:\n%s\nGot:\n%s", mockOutput, output)
	}
}

func TestGitListingRemotes_CommandFails(t *testing.T) {
	execCommand = mockGitRemoteCommand("", true)
	defer func() { execCommand = exec.Command }()

	input := &genai.FunctionCall{Args: map[string]any{}}

	_, err := GitListingRemotes(input)
	if err == nil {
		t.Fatal("Expected error from simulated git remote -v failure, got nil")
	}
}
