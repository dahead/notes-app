package note

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Metadata represents the metadata for a note
type Metadata struct {
	Tags []string `json:"tags"`
}

// LoadMetadata loads metadata from a .meta file
func LoadMetadata(metaPath string) (*Metadata, error) {
	data, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Metadata{
				Tags: []string{},
			}, nil
		}
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var meta Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
	}

	// Ensure slice is not nil
	if meta.Tags == nil {
		meta.Tags = []string{}
	}

	return &meta, nil
}

// SaveMetadata saves metadata to a .meta file
func (m *Metadata) SaveMetadata(metaPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(metaPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metaPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}
