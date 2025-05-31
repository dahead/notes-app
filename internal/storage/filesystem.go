package storage

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "notes-app/internal/note"
)

// FileSystemStorage handles file system operations for notes
type FileSystemStorage struct {
    rootPath string
}

// NewFileSystemStorage creates a new filesystem storage
func NewFileSystemStorage(rootPath string) *FileSystemStorage {
    return &FileSystemStorage{
        rootPath: rootPath,
    }
}

// Initialize creates the root directory if it doesn't exist
func (fs *FileSystemStorage) Initialize() error {
    return os.MkdirAll(fs.rootPath, 0755)
}

// GetAllNotes returns all notes in the storage
func (fs *FileSystemStorage) GetAllNotes() ([]*note.Note, error) {
    var notes []*note.Note

    err := filepath.Walk(fs.rootPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() && strings.HasSuffix(path, ".note") {
            note, err := note.LoadNote(path)
            if err != nil {
                fmt.Printf("Warning: failed to load note %s: %v\n", path, err)
                return nil // Continue walking
            }
            notes = append(notes, note)
        }

        return nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to walk directory: %w", err)
    }

    return notes, nil
}

// GetNote loads a specific note by path
func (fs *FileSystemStorage) GetNote(notePath string) (*note.Note, error) {
    fullPath := filepath.Join(fs.rootPath, notePath)
    if !strings.HasSuffix(fullPath, ".note") {
        fullPath += ".note"
    }

    return note.LoadNote(fullPath)
}

// CreateNote creates a new note
func (fs *FileSystemStorage) CreateNote(notePath, content string) (*note.Note, error) {
    fullPath := filepath.Join(fs.rootPath, notePath)
    if !strings.HasSuffix(fullPath, ".note") {
        fullPath += ".note"
    }

    // Check if note already exists
    if _, err := os.Stat(fullPath); err == nil {
        return nil, fmt.Errorf("note already exists: %s", notePath)
    }

    newNote := note.NewNote(fullPath)
    newNote.Content = content

    if err := newNote.Save(); err != nil {
        return nil, fmt.Errorf("failed to save new note: %w", err)
    }

    return newNote, nil
}

// GetRootPath returns the root path of the storage
func (fs *FileSystemStorage) GetRootPath() string {
    return fs.rootPath
}
