package common

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetDefaultNotesPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./notes"
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "Documents", "Notes")
	case "darwin":
		return filepath.Join(homeDir, "Documents", "Notes")
	default: // linux and others
		return filepath.Join(homeDir, "Notes")
	}
}
