package config

import (
	"os"
	"path/filepath"
	"testing"
)

// changeDirAndRestore changes to the specified directory and returns a function
// that restores the original directory. It handles errors properly for linting.
func changeDirAndRestore(t *testing.T, newDir string) func() {
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	if err := os.Chdir(newDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	return func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("Failed to restore original directory: %v", err)
		}
	}
}

func TestLoadEnv(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "nannytracker-env-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock project structure with go.mod
	projectRoot := filepath.Join(tempDir, "project")
	if err := os.MkdirAll(projectRoot, 0755); err != nil {
		t.Fatalf("Failed to create project root: %v", err)
	}

	// Create go.mod file to mark as project root
	goModPath := filepath.Join(projectRoot, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create .env file in project root
	envPath := filepath.Join(projectRoot, ".env")
	envContent := "TEST_VAR=test_value\nGOOGLE_MAPS_API_KEY=test_api_key"
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	// Create a subdirectory to test from
	subDir := filepath.Join(projectRoot, "cmd", "test")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Change to subdirectory and restore on cleanup
	restore := changeDirAndRestore(t, subDir)
	defer restore()

	// Clear any existing environment variables
	os.Unsetenv("TEST_VAR")
	os.Unsetenv("GOOGLE_MAPS_API_KEY")

	// Test LoadEnv
	LoadEnv()

	// Verify that environment variables were loaded
	if os.Getenv("TEST_VAR") != "test_value" {
		t.Errorf("Expected TEST_VAR to be 'test_value', got '%s'", os.Getenv("TEST_VAR"))
	}

	if os.Getenv("GOOGLE_MAPS_API_KEY") != "test_api_key" {
		t.Errorf("Expected GOOGLE_MAPS_API_KEY to be 'test_api_key', got '%s'", os.Getenv("GOOGLE_MAPS_API_KEY"))
	}
}

func TestLoadEnvNoProjectRoot(t *testing.T) {
	// Create a temporary directory that's not a project root
	tempDir, err := os.MkdirTemp("", "nannytracker-env-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory and restore on cleanup
	restore := changeDirAndRestore(t, tempDir)
	defer restore()

	// Clear any existing environment variables
	os.Unsetenv("TEST_VAR")

	// Test LoadEnv - should not panic and should fall back gracefully
	LoadEnv()

	// The function should complete without error, even though no .env file exists
	// We can't easily test the fallback behavior without mocking, but we can ensure it doesn't panic
}

func TestFindProjectRoot(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "nannytracker-root-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock project structure
	projectRoot := filepath.Join(tempDir, "project")
	if err := os.MkdirAll(projectRoot, 0755); err != nil {
		t.Fatalf("Failed to create project root: %v", err)
	}

	// Create go.mod file to mark as project root
	goModPath := filepath.Join(projectRoot, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create a subdirectory
	subDir := filepath.Join(projectRoot, "cmd", "test")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Change to subdirectory and restore on cleanup
	restore := changeDirAndRestore(t, subDir)
	defer restore()

	// Test findProjectRoot
	root, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	if root != projectRoot {
		t.Errorf("Expected project root to be %s, got %s", projectRoot, root)
	}
}

func TestFindProjectRootNotFound(t *testing.T) {
	// Create a temporary directory that's not a project root
	tempDir, err := os.MkdirTemp("", "nannytracker-root-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory and restore on cleanup
	restore := changeDirAndRestore(t, tempDir)
	defer restore()

	// Test findProjectRoot - should return an error
	_, err = findProjectRoot()
	if err == nil {
		t.Error("Expected findProjectRoot to return an error when no go.mod is found")
	}
}
