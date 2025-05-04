package gitTools

import (
	"os"
	"path/filepath"
	"testing"
	"os/exec"
	"strings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"
)

// Helper to create a temporary .gitignore file
func createTempGitignore(dir string, content string) (string, error) {
	path := filepath.Join(dir, ".gitignore")
	err := os.WriteFile(path, []byte(content), 0644)
	return path, err
}

func TestEnsureGitignore_CreatesFileWithDotEnv(t *testing.T) {
    // Change working directory to root (one level up)
    oldWd, _ := os.Getwd()
    defer os.Chdir(oldWd)
    os.Chdir("../../") // Go to the root directory where .gitignore should exist

    err := EnsureGitignore()
    require.NoError(t, err)

    // Read the .gitignore file and check if it contains .env
    data, err := os.ReadFile(".gitignore")
    require.NoError(t, err) // Ensure there was no error reading the file
    assert.Contains(t, string(data), ".env") // Ensure .env is in the file
}


func TestEnsureGitignore_AppendsDotEnv(t *testing.T) {
    // Same setup to change to the root directory
    oldWd, _ := os.Getwd()
    defer os.Chdir(oldWd)
    os.Chdir("../../")

    err := EnsureGitignore()
    require.NoError(t, err)

    // Ensure .gitignore was created or exists
    if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
        t.Fatalf(".gitignore does not exist, and couldn't create it!")
    }

    data, err := os.ReadFile(".gitignore")
    require.NoError(t, err)
    assert.Contains(t, string(data), ".env")
}



func TestEnsureGitignore_SkipsIfDotEnvExists(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	_ = os.WriteFile(".gitignore", []byte(".env\n"), 0644)

	err := EnsureGitignore()
	assert.NoError(t, err)

	data, _ := os.ReadFile(".gitignore")
	assert.Equal(t, ".env\n", string(data)) // Should not duplicate
}

func TestGitAddFile_AddsAllFiles(t *testing.T) {
    input := &genai.FunctionCall{
        Args: map[string]interface{}{
            "name": "",
        },
    }

    output, err := GitAddFile(input)

    // Handle possible warning or clean it up
    output = strings.TrimSpace(output)
    if strings.Contains(output, "LF will be replaced by CRLF") {
        output = "" // Ignore the line ending warning
    }

    require.NoError(t, err)
    require.Equal(t, "", output) // Check that no errors occurred and output is as expected
}

func TestGitAddFile_AddsSpecificFile(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	_ = os.WriteFile("example.txt", []byte("hi"), 0644)
	_ = exec.Command("git", "init").Run()

	input := &genai.FunctionCall{
		Args: map[string]any{
			"name": "example.txt",
		},
	}

	out, err := GitAddFile(input)
	assert.NoError(t, err)
	assert.Equal(t, "", out)
}
