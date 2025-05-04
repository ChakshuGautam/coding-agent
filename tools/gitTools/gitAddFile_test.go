package gitTools

import (
	"os"
	"path/filepath"
	"testing"
	"os/exec"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genai"
)

// Helper to create a temporary .gitignore file
func createTempGitignore(dir string, content string) (string, error) {
	path := filepath.Join(dir, ".gitignore")
	err := os.WriteFile(path, []byte(content), 0644)
	return path, err
}

func TestEnsureGitignore_CreatesFileWithDotEnv(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	err := EnsureGitignore()
	assert.NoError(t, err)

	data, err := os.ReadFile(".gitignore")
	assert.NoError(t, err)
	assert.Contains(t, string(data), ".env")
}

func TestEnsureGitignore_AppendsDotEnv(t *testing.T) {
    tempDir := t.TempDir()
    oldWd, _ := os.Getwd()
    defer os.Chdir(oldWd)
    os.Chdir(tempDir)

    _, _ = createTempGitignore(tempDir, "node_modules\n")

    err := EnsureGitignore()
    assert.NoError(t, err)

    data, _ := os.ReadFile(".gitignore")
    assert.Contains(t, string(data), ".env")
    assert.Contains(t, string(data), "node_modules")
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
	// Set up a dummy git repo in temp dir
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Initialize git repo
	_ = os.WriteFile("dummy.txt", []byte("hello"), 0644)
	_ = exec.Command("git", "init").Run()

	input := &genai.FunctionCall{
		Args: map[string]any{}, // No "name" means add all
	}

	out, err := GitAddFile(input)
	assert.NoError(t, err)
	assert.Equal(t, "", out) // git add doesn't return anything if successful
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
