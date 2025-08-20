package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
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
		fmt.Println("Error: Directory cannot be empty.")
		os.Exit(1)
	}

	// Expand ~ to home directory if used
	if strings.HasPrefix(targetDir, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error: Could not resolve home directory: %v\n", err)
			os.Exit(1)
		}
		targetDir = filepath.Join(home, targetDir[1:])
	}

	// Convert to absolute path for consistency
	targetDir, err := filepath.Abs(targetDir)
	if err != nil {
		fmt.Printf("Error: Could not resolve absolute path for %s: %v\n", targetDir, err)
		os.Exit(1)
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error: Could not create directory %s: %v\n", targetDir, err)
		os.Exit(1)
	}

	// Download the ZIP file
	zipURL := "https://github.com/Tmmcmasters/Gooo/archive/refs/heads/main.zip"
	zipPath := filepath.Join(targetDir, "gooo.zip")
	fmt.Printf("Downloading Gooo repository to %s...\n", zipPath)
	if err := downloadFile(zipURL, zipPath); err != nil {
		fmt.Printf("Error: Failed to download repository: %v\n", err)
		os.Exit(1)
	}

	// Extract the ZIP file
	fmt.Printf("Extracting %s to %s...\n", zipPath, targetDir)
	if err := extractZip(zipPath, targetDir); err != nil {
		fmt.Printf("Error: Failed to extract ZIP file: %v\n", err)
		os.Exit(1)
	}

	// Remove the ZIP file
	if err := os.Remove(zipPath); err != nil {
		fmt.Printf("Warning: Could not remove ZIP file: %v\n", err)
	}

	fmt.Printf("Success! Gooo project is scaffolded in %s\n", targetDir)
}

// downloadFile downloads a file from the given URL to the specified path
func downloadFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractZip extracts a ZIP file to the target directory, stripping the top-level directory
func extractZip(zipPath, targetDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Strip the top-level directory (e.g., "Gooo-main/")
		parts := strings.SplitN(file.Name, string(os.PathSeparator), 2)
		var filePath string
		if len(parts) > 1 {
			filePath = filepath.Join(targetDir, parts[1])
		} else {
			filePath = filepath.Join(targetDir, file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
