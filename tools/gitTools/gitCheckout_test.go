package gitTools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"google.golang.org/genai"
)

// backup and restore execCommand
var originalExecCommand = execCommand

func restoreExecCommand() {
	execCommand = originalExecCommand
}

// mockGitCommand returns a mock exec.Command function that returns specified output and error
func mockGitCommand(stdout string, simulateError bool) func(name string, args ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestMockedGitCommand", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(cmd.Env,
			"MOCK_GIT=1",
			"MOCK_STDOUT="+stdout,
			"MOCK_ERROR="+fmt.Sprintf("%v", simulateError),
		)
		return cmd
	}
}

// This function acts as a helper subprocess to simulate git command output
func TestMockedGitCommand(t *testing.T) {
	if os.Getenv("MOCK_GIT") != "1" {
		return
	}

	if os.Getenv("MOCK_ERROR") == "true" {
		os.Stderr.WriteString("simulated git error")
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, os.Getenv("MOCK_STDOUT"))
	os.Exit(0)
}

func TestGitCheckout_ExistingBranch(t *testing.T) {
	execCommand = mockGitCommand("* main\nfeature-a\n", false)
	defer restoreExecCommand()

	input := &genai.FunctionCall{
		Name: "checkout",
		Args: map[string]any{
			"name": "feature-a",
		},
	}

	output, err := GitCheckout(input)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !strings.Contains(output, "simulated") {
		t.Log("Checkout to existing branch simulated successfully.")
	}
}

func TestGitCheckout_NewBranch(t *testing.T) {
	execCommand = mockGitCommand("* main\nfeature-a\n", false)
	defer restoreExecCommand()

	input := &genai.FunctionCall{
		Name: "checkout",
		Args: map[string]any{
			"name": "feature-b",
		},
	}

	output, err := GitCheckout(input)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !strings.Contains(output, "simulated") {
		t.Log("New branch creation simulated successfully.")
	}
}

func TestGitCheckout_NoBranchName(t *testing.T) {
	execCommand = mockGitCommand("* main\n", false)
	defer restoreExecCommand()

	input := &genai.FunctionCall{
		Name: "checkout",
		Args: map[string]any{},
	}

	_, err := GitCheckout(input)
	if err == nil {
		t.Fatal("Expected error due to missing branch name, got none")
	}
}

func TestGitCheckout_CommandFails(t *testing.T) {
	execCommand = mockGitCommand("", true)
	defer restoreExecCommand()

	input := &genai.FunctionCall{
		Name: "checkout",
		Args: map[string]any{
			"name": "feature-x",
		},
	}

	_, err := GitCheckout(input)
	if err == nil {
		t.Fatal("Expected error from simulated git command failure, got none")
	}
}
