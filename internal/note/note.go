package note

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Note represents a single note with its content and metadata
type Note struct {
	Path     string
	Name     string
	Content  string
	Metadata *Metadata
	ModTime  time.Time
}

// NewNote creates a new note
func NewNote(notePath string) *Note {
	name := strings.TrimSuffix(filepath.Base(notePath), ".note")
	return &Note{
		Path:    notePath,
		Name:    name,
		Content: "",
		Metadata: &Metadata{
			Tags: []string{},
		},
	}
}

// LoadNote loads a note from the filesystem
func LoadNote(notePath string) (*Note, error) {
	// Check if note file exists
	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("note file does not exist: %s", notePath)
	}

	note := NewNote(notePath)

	// Load content
	content, err := os.ReadFile(notePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read note content: %w", err)
	}
	note.Content = string(content)

	// Get modification time
	info, err := os.Stat(notePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	note.ModTime = info.ModTime()

	// Load metadata
	metaPath := strings.TrimSuffix(notePath, ".note") + ".meta"
	metadata, err := LoadMetadata(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}
	note.Metadata = metadata

	return note, nil
}

// Save saves the note and its metadata to the filesystem
func (n *Note) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(n.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Save note content
	if err := os.WriteFile(n.Path, []byte(n.Content), 0644); err != nil {
		return fmt.Errorf("failed to save note content: %w", err)
	}

	// Save metadata
	metaPath := strings.TrimSuffix(n.Path, ".note") + ".meta"
	if err := n.Metadata.SaveMetadata(metaPath); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// Delete removes the note and its metadata from the filesystem
func (n *Note) Delete() error {
	// Delete note file
	if err := os.Remove(n.Path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete note file: %w", err)
	}

	// Delete metadata file
	metaPath := strings.TrimSuffix(n.Path, ".note") + ".meta"
	if err := os.Remove(metaPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata file: %w", err)
	}

	return nil
}

// GetMetaPath returns the path to the metadata file
func (n *Note) GetMetaPath() string {
	return strings.TrimSuffix(n.Path, ".note") + ".meta"
}
