package contract

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadSpec(ctx context.Context, path string) (*openapi3.T, error) {
	// Check if file exists first
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("openapi spec file not found: %s", path)
	}

	l := &openapi3.Loader{
		IsExternalRefsAllowed: true,
		Context:               ctx,
	}

	doc, err := l.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("openapi load from %s: %w", path, err)
	}

	if err := doc.Validate(ctx); err != nil {
		return nil, fmt.Errorf("openapi validate %s: %w", path, err)
	}

	return doc, nil
}

func SpecPath() string {
	// Allow override in CI via SPEC_PATH environment variable
	if p := os.Getenv("SPEC_PATH"); p != "" {
		// Check if the environment path exists
		if _, err := os.Stat(p); err == nil {
			fmt.Printf("Found OpenAPI spec via SPEC_PATH: %s\n", p)
			return p
		}
		fmt.Printf("SPEC_PATH set but file not found: %s\n", p)
	}

	// Get current working directory for debugging
	wd, _ := os.Getwd()

	// Try relative paths from different working directories
	candidates := []string{
		// From test/contract directory (when running locally)
		filepath.Join("..", "testdata", "openapi.yaml"),
		// From project root (when running in CI)
		filepath.Join("test", "testdata", "openapi.yaml"),
		// Fallback paths
		filepath.Join("..", "..", "test", "testdata", "openapi.yaml"),
		filepath.Join("testdata", "openapi.yaml"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			absPath, _ := filepath.Abs(candidate)
			fmt.Printf("Found OpenAPI spec at: %s\n", absPath)
			return candidate
		}
	}

	// Debug: print current working directory and what we tried
	fmt.Printf("Current working directory: %s\n", wd)
	fmt.Println("Tried the following paths:")
	for _, candidate := range candidates {
		absPath, _ := filepath.Abs(candidate)
		fmt.Printf("  - %s (absolute: %s)\n", candidate, absPath)
	}

	// Default fallback - prefer the path from test/contract directory
	return filepath.Join("..", "testdata", "openapi.yaml")
}

// Helper function to get absolute spec path for debugging
func AbsoluteSpecPath() (string, error) {
	specPath := SpecPath()
	absPath, err := filepath.Abs(specPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %w", specPath, err)
	}
	return absPath, nil
}

// Validate that the spec file exists and is readable
func ValidateSpecFile() error {
	specPath := SpecPath()

	// Check if file exists
	info, err := os.Stat(specPath)
	if os.IsNotExist(err) {
		// Print debug info
		wd, _ := os.Getwd()
		fmt.Printf("Working directory: %s\n", wd)
		fmt.Printf("Looking for spec at: %s\n", specPath)
		absPath, _ := filepath.Abs(specPath)
		fmt.Printf("Absolute path: %s\n", absPath)

		return fmt.Errorf("openapi spec file does not exist: %s", specPath)
	}
	if err != nil {
		return fmt.Errorf("error accessing spec file %s: %w", specPath, err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("spec path is not a regular file: %s", specPath)
	}

	// Check if file is readable
	file, err := os.Open(specPath)
	if err != nil {
		return fmt.Errorf("spec file is not readable: %s: %w", specPath, err)
	}
	file.Close()

	return nil
}
