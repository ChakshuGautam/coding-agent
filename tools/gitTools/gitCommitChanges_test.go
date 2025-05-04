package gitTools

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"google.golang.org/genai"
)

// Mock command test runner
func TestMockedGitCommand(t *testing.T) {
	if os.Getenv("MOCK_GIT") != "1" {
		return
	}

	if os.Getenv("MOCK_ERROR") == "true" {
		os.Stderr.WriteString("simulated error")
		os.Exit(1)
	}

	os.Stdout.WriteString(os.Getenv("MOCK_STDOUT"))
	os.Exit(0)
}

func mockGitCommand(stdout string, simulateError bool) func(name string, args ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestMockedGitCommand", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(os.Environ(),
			"MOCK_GIT=1",
			"MOCK_STDOUT="+stdout,
			"MOCK_ERROR="+boolToString(simulateError),
		)
		return cmd
	}
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// === TESTS ===

func TestGitCommitChanges_WithMessage(t *testing.T) {
	execCommand = mockGitCommand("[main 123abc] user msg\n", false)
	defer func() { execCommand = exec.Command }()

	input := &genai.FunctionCall{
		Args: map[string]any{
			"message": "user msg",
		},
	}

	output, err := GitCommitChanges(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if output != "[main 123abc] user msg\n" {
		t.Errorf("Unexpected commit output: %s", output)
	}
}



func TestGitCommitChanges_NoMessage_AICommit(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		if args[0] == "diff" {
			return mockGitCommand("++ added new feature\n", false)(name, args...)
		}
		if args[0] == "commit" {
			return mockGitCommand("[main abc123] auto-gen msg\n", false)(name, args...)
		}
		return exec.Command(name, args...)
	}
	defer func() { execCommand = exec.Command }()

	// Override GenerateCommitMessage
	originalGen := GenerateCommitMessage
	GenerateCommitMessage = func(diff string, optClient*genai.Client) (string, error) {
		return "auto-gen msg", nil
	}
	defer func() { GenerateCommitMessage = originalGen }()

	input := &genai.FunctionCall{
		Args: map[string]any{},
	}

	output, err := GitCommitChanges(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if output != "[main abc123] auto-gen msg\n" {
		t.Errorf("Unexpected commit output: %s", output)
	}
}

func TestGitCommitChanges_NoStagedChanges(t *testing.T) {
	execCommand = mockGitCommand("", false)
	defer func() { execCommand = exec.Command }()

	input := &genai.FunctionCall{
		Args: map[string]any{},
	}

	_, err := GitCommitChanges(input)
	if err == nil {
		t.Error("Expected error due to no staged changes, got nil")
	}
}

func TestGitCommitChanges_GitCommitFails(t *testing.T) {
	execCommand = mockGitCommand("", true)
	defer func() { execCommand = exec.Command }()

	input := &genai.FunctionCall{
		Args: map[string]any{
			"message": "test commit",
		},
	}

	_, err := GitCommitChanges(input)
	if err == nil {
		t.Error("Expected error from failed git commit, got nil")
	}
}

func TestGitCommitChanges_AIGenerationFails(t *testing.T) {
	execCommand = mockGitCommand("++ added line\n", false)
	defer func() { execCommand = exec.Command }()

	// Override GenerateCommitMessage
	originalGen := GenerateCommitMessage
	GenerateCommitMessage = func(diff string,optClient *genai.Client) (string, error) {
		return "", fmt.Errorf("AI failed")
	}
	defer func() { GenerateCommitMessage = originalGen }()

	input := &genai.FunctionCall{
		Args: map[string]any{},
	}

	_, err := GitCommitChanges(input)
	if err == nil || err.Error() != "AI generation failed: AI failed" {
		t.Errorf("Expected AI generation error, got: %v", err)
	}
}
