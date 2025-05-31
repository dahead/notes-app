package index

import (
	"fmt"
	"strings"

	"notes-app/internal/note"
)

// Index maintains an in-memory index of notes
type Index struct {
	notes    []*note.Note
	tagIndex map[string][]*note.Note
}

// NewIndex creates a new index
func NewIndex() *Index {
	return &Index{
		notes:    []*note.Note{},
		tagIndex: make(map[string][]*note.Note),
	}
}

// AddNote adds a note to the index
func (idx *Index) AddNote(n *note.Note) {
	idx.notes = append(idx.notes, n)
	idx.updateIndices(n)
}

// RemoveNote removes a note from the index
func (idx *Index) RemoveNote(notePath string) {
	for i, n := range idx.notes {
		if n.Path == notePath {
			idx.removeFromIndices(n)
			idx.notes = append(idx.notes[:i], idx.notes[i+1:]...)
			break
		}
	}
}

// UpdateNote updates a note in the index
func (idx *Index) UpdateNote(n *note.Note) {
	idx.RemoveNote(n.Path)
	idx.AddNote(n)
}

// SearchByTag searches notes by tag
func (idx *Index) SearchByTag(tag string) []*note.Note {
	return idx.tagIndex[strings.ToLower(tag)]
}

// SearchByContent searches notes by content (simple text search)
func (idx *Index) SearchByContent(query string) []*note.Note {
	var results []*note.Note
	queryLower := strings.ToLower(query)

	for _, n := range idx.notes {
		if strings.Contains(strings.ToLower(n.Content), queryLower) ||
			strings.Contains(strings.ToLower(n.Name), queryLower) {
			results = append(results, n)
		}
	}

	return results
}

// GetAllNotes returns all notes in the index
func (idx *Index) GetAllNotes() []*note.Note {
	return idx.notes
}

// GetAllTags returns all unique tags
func (idx *Index) GetAllTags() []string {
	var tags []string
	for tag := range idx.tagIndex {
		tags = append(tags, tag)
	}
	return tags
}

// RebuildIndex rebuilds the entire index
func (idx *Index) RebuildIndex(notes []*note.Note) {
	idx.notes = []*note.Note{}
	idx.tagIndex = make(map[string][]*note.Note)

	for _, n := range notes {
		idx.AddNote(n)
	}
}

// updateIndices updates all indices for a note
func (idx *Index) updateIndices(n *note.Note) {
	// Update tag index
	for _, tag := range n.Metadata.Tags {
		tagLower := strings.ToLower(tag)
		idx.tagIndex[tagLower] = append(idx.tagIndex[tagLower], n)
	}
}

// removeFromIndices removes a note from all indices
func (idx *Index) removeFromIndices(n *note.Note) {
	// Remove from tag index
	for _, tag := range n.Metadata.Tags {
		tagLower := strings.ToLower(tag)
		idx.removeNoteFromSlice(idx.tagIndex[tagLower], n)
	}
}

// removeNoteFromSlice removes a note from a slice
func (idx *Index) removeNoteFromSlice(slice []*note.Note, noteToRemove *note.Note) []*note.Note {
	for i, n := range slice {
		if n.Path == noteToRemove.Path {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// PrintStats prints index statistics
func (idx *Index) PrintStats() {
	fmt.Printf("Index Statistics:\n")
	fmt.Printf("  Total notes: %d\n", len(idx.notes))
	fmt.Printf("  Unique tags: %d\n", len(idx.tagIndex))
}
