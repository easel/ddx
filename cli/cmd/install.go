package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// runInstall implements the install command logic
func runInstall(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Installing DDx...")

	// Get installation parameters
	version, _ := cmd.Flags().GetString("version")
	installPath, _ := cmd.Flags().GetString("path")
	force, _ := cmd.Flags().GetBool("force")

	if version == "" {
		version = "latest"
	}

	if installPath == "" {
		var err error
		installPath, err = getDefaultInstallPath()
		if err != nil {
			return fmt.Errorf("failed to determine install path: %w", err)
		}
	}

	// Check if already installed
	if !force {
		existing := filepath.Join(installPath, getBinaryName())
		if _, err := os.Stat(existing); err == nil {
			fmt.Printf("⚠️  DDx is already installed at %s\n", existing)
			fmt.Println("Use --force to overwrite or --path to install elsewhere")
			return nil
		}
	}

	// Create install directory
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Download and install
	fmt.Printf("📦 Downloading DDx %s for %s/%s...\n", version, runtime.GOOS, runtime.GOARCH)

	downloadURL, err := getDownloadURL(version)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %w", err)
	}

	if err := downloadAndInstall(downloadURL, installPath); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Verify installation
	binaryPath := filepath.Join(installPath, getBinaryName())
	if _, err := os.Stat(binaryPath); err != nil {
		return fmt.Errorf("installation verification failed: binary not found at %s", binaryPath)
	}

	// Setup PATH if needed
	if err := setupPath(installPath); err != nil {
		fmt.Printf("⚠️  Warning: Could not setup PATH automatically: %v\n", err)
		fmt.Printf("Please add %s to your PATH manually\n", installPath)
	}

	fmt.Printf("✅ DDx installed successfully to %s\n", binaryPath)
	fmt.Println("💡 Run 'ddx version' to verify installation")

	return nil
}

// getDefaultInstallPath returns the default installation path for the current platform
func getDefaultInstallPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "bin"), nil
	default:
		return filepath.Join(homeDir, ".local", "bin"), nil
	}
}

// getBinaryName returns the binary name for the current platform
func getBinaryName() string {
	if runtime.GOOS == "windows" {
		return "ddx.exe"
	}
	return "ddx"
}

// getDownloadURL constructs the download URL for the specified version
func getDownloadURL(version string) (string, error) {
	if version == "latest" {
		version = "v0.0.1" // For now, use a fixed version
	}

	platform := runtime.GOOS
	arch := runtime.GOARCH

	// Map Go arch names to common names
	switch arch {
	case "amd64":
		arch = "x86_64"
	case "386":
		arch = "i386"
	}

	var ext string
	if platform == "windows" {
		ext = "zip"
	} else {
		ext = "tar.gz"
	}

	// Use GitHub releases URL pattern
	baseURL := "https://github.com/easel/ddx/releases/download"
	filename := fmt.Sprintf("ddx_%s_%s_%s.%s", version, platform, arch, ext)

	return fmt.Sprintf("%s/%s/%s", baseURL, version, filename), nil
}

// downloadAndInstall downloads and installs the binary
func downloadAndInstall(url, installPath string) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	// Extract based on file extension
	if strings.HasSuffix(url, ".zip") {
		return extractZip(resp.Body, installPath)
	} else {
		return extractTarGz(resp.Body, installPath)
	}
}

// extractZip extracts a ZIP archive
func extractZip(r io.Reader, installPath string) error {
	// For simplicity, we'll create a temp file then extract
	tempFile, err := os.CreateTemp("", "ddx-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, r); err != nil {
		return err
	}

	// Reopen for reading
	zipReader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if strings.Contains(file.Name, getBinaryName()) {
			return extractZipFile(file, installPath)
		}
	}

	return fmt.Errorf("binary not found in archive")
}

// extractZipFile extracts a single file from ZIP
func extractZipFile(file *zip.File, installPath string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outPath := filepath.Join(installPath, getBinaryName())
	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

// extractTarGz extracts a tar.gz archive
func extractTarGz(r io.Reader, installPath string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if strings.Contains(header.Name, getBinaryName()) && header.Typeflag == tar.TypeReg {
			outPath := filepath.Join(installPath, getBinaryName())
			outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tr)
			return err
		}
	}

	return fmt.Errorf("binary not found in archive")
}

// setupPath attempts to add the install path to PATH
func setupPath(installPath string) error {
	switch runtime.GOOS {
	case "windows":
		return setupWindowsPath(installPath)
	default:
		return setupUnixPath(installPath)
	}
}

// setupUnixPath adds to PATH via shell profile
func setupUnixPath(installPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Try to update shell profile
	profiles := []string{
		".bashrc",
		".zshrc",
		".profile",
	}

	pathExport := fmt.Sprintf("export PATH=\"%s:$PATH\"\n", installPath)

	for _, profile := range profiles {
		profilePath := filepath.Join(homeDir, profile)
		if _, err := os.Stat(profilePath); err == nil {
			return appendToFile(profilePath, pathExport)
		}
	}

	// Create .profile if none exist
	profilePath := filepath.Join(homeDir, ".profile")
	return appendToFile(profilePath, pathExport)
}

// setupWindowsPath adds to PATH via user environment
func setupWindowsPath(installPath string) error {
	// For Windows, we'd need to modify the registry or use environment commands
	// For now, just provide instructions
	fmt.Printf("💡 Add %s to your PATH environment variable\n", installPath)
	return nil
}

// appendToFile appends content to a file
func appendToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}