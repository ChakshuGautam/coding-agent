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

// mockGitChechoutCommand returns a mock exec.Command function that returns specified output and error
func mockGitChechoutCommand(stdout string, simulateError bool) func(name string, args ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestGitChechoutCommand", "--", name}
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
func TestGitChechoutCommand(t *testing.T) {
	if os.Getenv("MOCK_GIT") != "1" {
		return
	}

	if os.Getenv("MOCK_ERROR") == "true" {
		os.Stderr.WriteString("simulated git error")
		os.Exit(1) // Simulates failure
	}

	// Simulate success
	fmt.Fprint(os.Stdout, os.Getenv("MOCK_STDOUT"))
	os.Exit(0) // Critical: ensures successful exit
}


func TestGitCheckout_ExistingBranch(t *testing.T) {
	execCommand = mockGitChechoutCommand("* main\nfeature-a\n", false)
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
	execCommand = mockGitChechoutCommand("* main\nfeature-a\n", false)
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
	execCommand = mockGitChechoutCommand("* main\n", false)
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
	execCommand = mockGitChechoutCommand("", true)
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
