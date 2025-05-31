package app

import (
	"fmt"
	"strings"

	"notes-app/internal/index"
	"notes-app/internal/note"
	"notes-app/internal/storage"
)

// NotesApp represents the main application
type NotesApp struct {
	storage *storage.FileSystemStorage
	index   *index.Index
}

// NewNotesApp creates a new notes application
func NewNotesApp(rootPath string) *NotesApp {
	return &NotesApp{
		storage: storage.NewFileSystemStorage(rootPath),
		index:   index.NewIndex(),
	}
}

// Initialize initializes the application
func (app *NotesApp) Initialize() error {
	if err := app.storage.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	return app.RefreshIndex()
}

// RefreshIndex refreshes the note index
func (app *NotesApp) RefreshIndex() error {
	notes, err := app.storage.GetAllNotes()
	if err != nil {
		return fmt.Errorf("failed to load notes: %w", err)
	}

	app.index.RebuildIndex(notes)
	return nil
}

// CreateNote creates a new note
func (app *NotesApp) CreateNote(name, content string) error {
	_, err := app.storage.CreateNote(name, content)
	if err != nil {
		return err
	}

	return app.RefreshIndex()
}

// SearchNotes searches for notes
func (app *NotesApp) SearchNotes(query, searchType string) []*note.Note {
	switch searchType {
	case "tag":
		return app.index.SearchByTag(query)
	case "content":
		return app.index.SearchByContent(query)
	default:
		return app.index.SearchByContent(query)
	}
}

// UpdateNoteTags updates all tags for a note (replaces existing tags)
func (app *NotesApp) UpdateNoteTags(notePath string, tags []string) error {
	note, err := app.storage.GetNote(notePath)
	if err != nil {
		return err
	}

	note.Metadata.Tags = tags

	if err := note.Save(); err != nil {
		return err
	}

	return app.RefreshIndex()
}

// AddTagsToNote adds tags to a note without duplicates
func (app *NotesApp) AddTagsToNote(notePath string, newTags []string) error {
	note, err := app.storage.GetNote(notePath)
	if err != nil {
		return err
	}

	// Add new tags without duplicates
	tagSet := make(map[string]bool)
	for _, tag := range note.Metadata.Tags {
		tagSet[tag] = true
	}

	for _, tag := range newTags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tagSet[tag] = true
		}
	}

	// Convert back to slice
	note.Metadata.Tags = make([]string, 0, len(tagSet))
	for tag := range tagSet {
		note.Metadata.Tags = append(note.Metadata.Tags, tag)
	}

	if err := note.Save(); err != nil {
		return err
	}

	return app.RefreshIndex()
}

// RemoveTagsFromNote removes specific tags from a note
func (app *NotesApp) RemoveTagsFromNote(notePath string, tagsToRemove []string) error {
	note, err := app.storage.GetNote(notePath)
	if err != nil {
		return err
	}

	// Create set of tags to remove
	removeSet := make(map[string]bool)
	for _, tag := range tagsToRemove {
		removeSet[strings.TrimSpace(tag)] = true
	}

	// Filter out tags to remove
	var newTags []string
	for _, tag := range note.Metadata.Tags {
		if !removeSet[tag] {
			newTags = append(newTags, tag)
		}
	}

	note.Metadata.Tags = newTags

	if err := note.Save(); err != nil {
		return err
	}

	return app.RefreshIndex()
}

// ShowNoteTags displays tags for a specific note
func (app *NotesApp) ShowNoteTags(notePath string) error {
	note, err := app.storage.GetNote(notePath)
	if err != nil {
		return err
	}

	fmt.Printf("Tags for '%s': %v\n", note.Name, note.Metadata.Tags)
	return nil
}

// ListAllNotes returns all notes
func (app *NotesApp) ListAllNotes() []*note.Note {
	return app.index.GetAllNotes()
}

// GetAllTags returns all unique tags
func (app *NotesApp) GetAllTags() []string {
	return app.index.GetAllTags()
}

// GetStats prints application statistics
func (app *NotesApp) GetStats() {
	app.index.PrintStats()
}

// DeleteNote removes a note and its metadata
func (app *NotesApp) DeleteNote(notePath string) error {
	note, err := app.storage.GetNote(notePath)
	if err != nil {
		return err
	}

	if err := note.Delete(); err != nil {
		return err
	}

	return app.RefreshIndex()
}

// GetNote retrieves a specific note
func (app *NotesApp) GetNote(notePath string) (*note.Note, error) {
	return app.storage.GetNote(notePath)
}

// UpdateNoteContent updates the content of a note
func (app *NotesApp) UpdateNoteContent(notePath, content string) error {
	note, err := app.storage.GetNote(notePath)
	if err != nil {
		return err
	}

	note.Content = content

	if err := note.Save(); err != nil {
		return err
	}

	return app.RefreshIndex()
}
