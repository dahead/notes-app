package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"notes-app/internal/app"
	"notes-app/internal/note"
)

func getDefaultNotesPath() string {
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

func printHelp() {
	fmt.Println("Notes App - Commands:")
	fmt.Println("  help                          	- Show this help")
	fmt.Println("  create <name>                 	- Create a new note")
	fmt.Println("  list                          	- List all notes")
	fmt.Println("  search <query>                	- Search notes by content")
	fmt.Println("  search-tag <tag>              	- Search notes by tag")
	fmt.Println("  show-tags [note]              	- Show all tags or tags for specific note")
	fmt.Println("  set-tags <note> <tag1,tag2>   	- Set note tags (replace all)")
	fmt.Println("  add-tags <note> <tag1,tag2>   	- Add tags to a note")
	fmt.Println("  remove-tags <note> <tag1,tag2> 	- Remove tags from a note")
	fmt.Println("  list-tags                     	- List all unique tags")
	fmt.Println("  delete <note>                 	- Delete a note")
	fmt.Println("  edit <note>                   	- Edit note content")
	fmt.Println("  stats                         	- Show index statistics")
	fmt.Println("  refresh                       	- Refresh the index")
	fmt.Println("  clear                         	- Clear the screen")
	fmt.Println("  quit                          	- Exit the application")
}

func printNotes(notes []*note.Note) {
	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return
	}

	fmt.Printf("Found %d note(s):\n\n", len(notes))
	for _, n := range notes {
		fmt.Printf("Name: %s\n", n.Name)
		fmt.Printf("Path: %s\n", n.Path)
		fmt.Printf("Modified: %s\n", n.ModTime.Format("2006-01-02 15:04:05"))

		if len(n.Metadata.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(n.Metadata.Tags, ", "))
		}

		// Show first few lines of content
		lines := strings.Split(n.Content, "\n")
		preview := ""
		for i, line := range lines {
			if i >= 3 || len(preview) > 100 {
				preview += "..."
				break
			}
			if line != "" {
				preview += line + " "
			}
		}
		if preview != "" {
			fmt.Printf("Preview: %s\n", preview)
		}

		fmt.Println("---")
	}
}

func readMultilineInput(scanner *bufio.Scanner, prompt string) string {
	fmt.Printf("%s (type '.' on a new line to finish):\n", prompt)

	var content strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if line == "." {
			break
		}
		content.WriteString(line + "\n")
	}

	return content.String()
}

func main() {
	// Get notes directory from environment variable or use default
	notesPath := os.Getenv("NOTES_PATH")
	if notesPath == "" {
		notesPath = getDefaultNotesPath()
	}

	fmt.Printf("Notes App - Using directory: %s\n", notesPath)
	fmt.Println("Type 'help' for available commands")

	notesApp := app.NewNotesApp(notesPath)
	if err := notesApp.Initialize(); err != nil {
		fmt.Printf("Error initializing app: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded notes from: %s\n", notesPath)
	notesApp.GetStats()
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "help", "h":
			printHelp()

		case "create", "c":
			if len(parts) < 2 {
				fmt.Println("Usage: create <name>")
				continue
			}
			name := strings.Join(parts[1:], " ")
			content := readMultilineInput(scanner, "Enter note content")

			if err := notesApp.CreateNote(name, content); err != nil {
				fmt.Printf("Error creating note: %v\n", err)
			} else {
				fmt.Printf("Note '%s' created successfully\n", name)
			}

		case "list", "l":
			notes := notesApp.ListAllNotes()
			printNotes(notes)

		case "search", "?":
			if len(parts) < 2 {
				fmt.Println("Usage: search <query>")
				continue
			}
			query := strings.Join(parts[1:], " ")
			notes := notesApp.SearchNotes(query, "content")
			printNotes(notes)

		case "search-tag":
			if len(parts) < 2 {
				fmt.Println("Usage: search-tag <tag>")
				continue
			}
			tag := parts[1]
			notes := notesApp.SearchNotes(tag, "tag")
			printNotes(notes)

		case "show-tags":
			if len(parts) >= 2 {
				// Original behavior: show tags for a specific note
				noteName := parts[1]
				if err := notesApp.ShowNoteTags(noteName); err != nil {
					fmt.Printf("Error showing tags: %v\n", err)
				}
			} else {
				// New behavior: show all tags from all notes
				tags := notesApp.GetAllTags()
				if len(tags) == 0 {
					fmt.Println("No tags found.")
				} else {
					fmt.Printf("All tags (%d):\n", len(tags))
					for _, tag := range tags {
						fmt.Printf("  %s\n", tag)
					}
				}
			}

		case "set-tags", "st":
			if len(parts) < 2 {
				fmt.Println("Usage: set-tags <note-name> [tag1,tag2,tag3]")
				continue
			}
			noteName := parts[1]
			var tags []string
			if len(parts) > 2 {
				tagsStr := strings.Join(parts[2:], " ")
				tags = strings.Split(tagsStr, ",")
				for i, tag := range tags {
					tags[i] = strings.TrimSpace(tag)
				}
			}

			if err := notesApp.UpdateNoteTags(noteName, tags); err != nil {
				fmt.Printf("Error setting tags: %v\n", err)
			} else {
				fmt.Println("Tags updated successfully")
			}

		case "add-tags", "at":
			if len(parts) < 3 {
				fmt.Println("Usage: add-tags <note-name> <tag1,tag2,tag3>")
				continue
			}
			noteName := parts[1]
			tagsStr := strings.Join(parts[2:], " ")
			newTags := strings.Split(tagsStr, ",")

			if err := notesApp.AddTagsToNote(noteName, newTags); err != nil {
				fmt.Printf("Error adding tags: %v\n", err)
			} else {
				fmt.Println("Tags added successfully")
			}

		case "remove-tags", "rt":
			if len(parts) < 3 {
				fmt.Println("Usage: remove-tags <note-name> <tag1,tag2,tag3>")
				continue
			}
			noteName := parts[1]
			tagsStr := strings.Join(parts[2:], " ")
			tagsToRemove := strings.Split(tagsStr, ",")

			if err := notesApp.RemoveTagsFromNote(noteName, tagsToRemove); err != nil {
				fmt.Printf("Error removing tags: %v\n", err)
			} else {
				fmt.Println("Tags removed successfully")
			}

		case "list-tags", "lt":
			tags := notesApp.GetAllTags()
			if len(tags) == 0 {
				fmt.Println("No tags found.")
			} else {
				fmt.Printf("All tags (%d):\n", len(tags))
				for _, tag := range tags {
					fmt.Printf("  %s\n", tag)
				}
			}

		case "delete", "d":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <note-name>")
				continue
			}
			noteName := parts[1]

			fmt.Printf("Are you sure you want to delete '%s'? (y/N): ", noteName)
			scanner.Scan()
			confirmation := strings.ToLower(strings.TrimSpace(scanner.Text()))

			if confirmation == "y" || confirmation == "yes" {
				if err := notesApp.DeleteNote(noteName); err != nil {
					fmt.Printf("Error deleting note: %v\n", err)
				} else {
					fmt.Printf("Note '%s' deleted successfully\n", noteName)
				}
			} else {
				fmt.Println("Delete cancelled")
			}

		case "edit", "e":
			if len(parts) < 2 {
				fmt.Println("Usage: edit <note-name>")
				continue
			}
			noteName := parts[1]

			// Show current content first
			currentNote, err := notesApp.GetNote(noteName)
			if err != nil {
				fmt.Printf("Error loading note: %v\n", err)
				continue
			}

			fmt.Printf("Current content of '%s':\n", noteName)
			fmt.Println("---")
			fmt.Print(currentNote.Content)
			fmt.Println("---")

			newContent := readMultilineInput(scanner, "Enter new content")

			if err := notesApp.UpdateNoteContent(noteName, newContent); err != nil {
				fmt.Printf("Error updating note: %v\n", err)
			} else {
				fmt.Printf("Note '%s' updated successfully\n", noteName)
			}

		case "stats":
			notesApp.GetStats()

		case "refresh":
			if err := notesApp.RefreshIndex(); err != nil {
				fmt.Printf("Error refreshing index: %v\n", err)
			} else {
				fmt.Println("Index refreshed successfully")
				notesApp.GetStats()
			}

		case "clear", "cls":
			// Clear screen command
			switch runtime.GOOS {
			case "windows":
				// Windows
				fmt.Print("\033[H\033[2J")
			default:
				// Unix-like systems (Linux, macOS, etc.)
				fmt.Print("\033[2J\033[H")
			}

		case "quit", "exit", "q":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Printf("Unknown command: %s\n", command)
			fmt.Println("Type 'help' for available commands")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
