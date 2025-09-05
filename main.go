package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Prompt user for target directory
	fmt.Print("Enter the directory to download Gooo into (e.g., ~/projects/gooo or C:\\projects\\gooo): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	targetDir := strings.TrimSpace(scanner.Text())

	if targetDir == "" {
		log.Fatalf("Error: Directory cannot be empty.")
	}

	// Expand ~ to home directory if used
	if strings.HasPrefix(targetDir, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Error: Could not resolve home directory: %v\n", err)
		}
		targetDir = filepath.Join(home, targetDir[1:])
	}

	// Convert to absolute path for consistency
	targetDir, err := filepath.Abs(targetDir)
	if err != nil {
		log.Fatalf("Error: Could not resolve absolute path for %s: %v\n", targetDir, err)
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Fatalf("Error: Could not create directory %s: %v\n", targetDir, err)
	}

	// Download the ZIP file
	zipURL := "https://github.com/Tmmcmasters/Gooo/archive/refs/heads/main.zip"
	zipPath := filepath.Join(targetDir, "gooo.zip")
	fmt.Printf("Downloading Gooo repository to %s...\n", zipPath)
	if err := downloadFile(zipURL, zipPath); err != nil {
		log.Fatalf("Error: Failed to download repository: %v\n", err)
	}

	// Extract the ZIP file
	fmt.Printf("Extracting %s to %s...\n", zipPath, targetDir)
	if err := extractZip(zipPath, targetDir); err != nil {
		log.Fatalf("Error: Failed to extract ZIP file: %v\n", err)
	}

	// Remove the ZIP file
	if err := os.Remove(zipPath); err != nil {
		fmt.Printf("Warning: Could not remove ZIP file: %v\n", err)
	}

	fmt.Printf("Success! Gooo project is scaffolded in %s\n", targetDir)
}
