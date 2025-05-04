package gitTools

import (
	"os"
	"os/exec"
	"testing"
	"fmt"

	"google.golang.org/genai"
)

func TestMockedGitPushCommand(t *testing.T) {
	if os.Getenv("MOCK_GIT_PUSH") != "1" {
		return
	}

	if os.Getenv("MOCK_ERROR") == "true" {
		os.Stderr.WriteString("simulated error")
		os.Exit(1)
	}

	os.Stdout.WriteString(os.Getenv("MOCK_STDOUT"))
	os.Exit(0)
}

func mockGitPushCommand(stdout string, simulateError bool) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestMockedGitPushCommand", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(os.Environ(),
			"MOCK_GIT_PUSH=1",
			"MOCK_STDOUT="+stdout,
			"MOCK_ERROR="+boolToString(simulateError),
		)
		return cmd
	}
}


func TestGitAddRemoteAndPush_Success(t *testing.T) {
	mockOutput := `Enumerating objects: 10, done.
Counting objects: 100% (10/10), done.
Delta compression using up to 8 threads
Compressing objects: 100% (7/7), done.
Writing objects: 100% (10/10), 1.38 KiB | 1.38 MiB/s, done.
Total 10 (delta 2), reused 0 (delta 0), pack-reused 0
To https://github.com/username/reponame.git
 * [new branch]      main -> main`

	execCommand = mockGitPushCommand(mockOutput, false)
	defer func() { execCommand = exec.Command }()

	input := &genai.FunctionCall{
		Args: map[string]any{
			"Name":     "reponame",
			"Branch":   "main",
			"UserName": "username",
		},
	}

	output, err := GitAddRemoteAndPush(input)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := fmt.Sprintf("Push to remote completed successfully. %s", mockOutput)
	if output != expected {
		t.Errorf("Unexpected output.\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestGitAddRemoteAndPush_CommandFails(t *testing.T) {
	execCommand = mockGitPushCommand("", true)
	defer func() { execCommand = exec.Command }()

	input := &genai.FunctionCall{
		Args: map[string]any{
			"Name":     "reponame",
			"Branch":   "main",
			"UserName": "username",
		},
	}

	_, err := GitAddRemoteAndPush(input)
	if err == nil {
		t.Fatal("Expected error from simulated git push failure, got nil")
	}
}

func TestGitAddRemoteAndPush_ExistingRemote(t *testing.T) {
	mockOutput := "Remote 'origin' already exists, skipping remote add."

	execCommand = mockGitPushCommand(mockOutput, false)
	defer func() { execCommand = exec.Command }()

	input := &genai.FunctionCall{
		Args: map[string]any{
			"Name":     "reponame",
			"Branch":   "main",
			"UserName": "username",
		},
	}

	output, err := GitAddRemoteAndPush(input)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := fmt.Sprintf("Push to remote completed successfully. %s", mockOutput)
	if output != expected {
		t.Errorf("Unexpected output.\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestGitAddRemoteAndPush_MissingRepoName(t *testing.T) {
	input := &genai.FunctionCall{
		Args: map[string]any{
			"Branch":   "main",
			"UserName": "username",
		},
	}

	_, err := GitAddRemoteAndPush(input)
	if err == nil {
		t.Fatal("Expected error for missing repo name, got nil")
	}
}
