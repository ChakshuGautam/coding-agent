package gitTools

import (
	"os"
	"os/exec"
	"testing"

	"google.golang.org/genai"
)

// Mock runner function
func TestGitStatusCommand(t *testing.T) {
	if os.Getenv("MOCK_GIT") != "1" {
		return
	}

	if os.Getenv("MOCK_ERROR") == "true" {
		os.Stderr.WriteString("simulated git error")
		os.Exit(1)
	}

	os.Stdout.WriteString(os.Getenv("MOCK_STDOUT"))
	os.Exit(0)
}

func mockGitStatusCommand(stdout string, simulateError bool) func(name string, args ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestGitStatusCommand", "--", name}
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

func TestGitFileStatus_Success(t *testing.T) {
	// Arrange
	execCommand = mockGitStatusCommand("On branch main\nnothing to commit", false)
	defer func() { execCommand = exec.Command }()

	input := &genai.FunctionCall{
		Args: map[string]any{},
	}

	// Act
	output, err := GitFileStatus(input)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if output != "On branch main\nnothing to commit" {
		t.Errorf("Unexpected output: %s", output)
	}
}

func TestGitFileStatus_Failure(t *testing.T) {
	// Arrange
	execCommand = mockGitStatusCommand("", true)
	defer func() { execCommand = exec.Command }()

	input := &genai.FunctionCall{
		Args: map[string]any{},
	}

	// Act
	_, err := GitFileStatus(input)

	// Assert
	if err == nil {
		t.Error("Expected error from simulated git status failure, got nil")
	}
}
