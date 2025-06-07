package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"notes-app/internal/app"
	"notes-app/internal/common"
	"notes-app/internal/ui"
	"os"
)

func main() {
	notesPath := os.Getenv("NOTES_PATH")
	if notesPath == "" {
		notesPath = common.GetDefaultNotesPath()
	}

	notesApp := app.NewNotesApp(notesPath)
	if err := notesApp.Initialize(); err != nil {
		fmt.Printf("Error initializing app: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		ui.NewModel(notesApp),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
