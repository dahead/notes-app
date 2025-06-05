package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"notes-app/internal/app"
	"notes-app/internal/note"
)

type Model struct {
	notes       []*note.Note
	input       textinput.Model
	textarea    textarea.Model
	tagInput    textinput.Model
	notesApp    *app.NotesApp
	cursor      int
	selected    map[int]struct{}
	state       string // "list", "help", "create", "edit", "view", "create_name", "tags"
	tagEditMode string // "", "add", "remove"
	err         error
	newNoteName string
	showPreview bool // New field for preview toggle
}

func NewModel(notesApp *app.NotesApp) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter note name..."
	ti.Width = 58

	ta := textarea.New()
	ta.Placeholder = "Enter note content..."
	ta.ShowLineNumbers = false
	ta.SetWidth(78)
	ta.SetHeight(18)
	ta.KeyMap.InsertNewline = key.NewBinding(
		key.WithKeys("enter"),
	)

	tagInput := textinput.New()
	tagInput.Placeholder = "Enter tags (comma-separated)..."
	tagInput.Width = 58

	return Model{
		notes:       notesApp.ListAllNotes(),
		input:       ti,
		textarea:    ta,
		tagInput:    tagInput,
		notesApp:    notesApp,
		selected:    make(map[int]struct{}),
		state:       "list",
		showPreview: false, // Initialize preview as hidden
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle tag input updates first when in input mode
	if m.state == "tags" && m.tagEditMode != "" && m.tagInput.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc:
				m.tagInput.Reset()
				m.tagEditMode = ""
				return m, nil
			case tea.KeyEnter:
				tags := strings.Split(m.tagInput.Value(), ",")
				var validTags []string
				for _, tag := range tags {
					tag = strings.TrimSpace(tag)
					if tag != "" {
						validTags = append(validTags, tag)
					}
				}

				if len(validTags) > 0 {
					note := m.notes[m.cursor]
					var err error
					if m.tagEditMode == "add" {
						err = m.notesApp.AddTagsToNote(note.Path, validTags)
					} else {
						err = m.notesApp.RemoveTagsFromNote(note.Path, validTags)
					}
					if err != nil {
						m.err = err
					} else {
						m.notes = m.notesApp.ListAllNotes()
						m.state = "list"
						m.tagInput.Reset()
						m.tagEditMode = ""
					}
				}
				return m, nil
			default:
				m.tagInput, cmd = m.tagInput.Update(msg)
				return m, cmd
			}
		}
		m.tagInput, cmd = m.tagInput.Update(msg)
		return m, cmd
	}

	// Handle other messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == "create" || m.state == "edit" || m.state == "tags" {
				m.state = "list"
				return m, nil
			}
			return m, tea.Quit

		case "?", "h":
			if m.state == "list" || m.state == "view" {
				m.state = "help"
			}

		case " ":
			if m.state == "list" {
				m.showPreview = !m.showPreview
				return m, nil
			}

		case "esc":
			if m.state == "create" || m.state == "edit" || m.state == "view" ||
				m.state == "help" || m.state == "create_name" || m.state == "tags" {
				m.state = "list"
				m.textarea.Reset()
				m.input.Reset()
				m.tagInput.Reset()
				m.newNoteName = ""
				m.tagEditMode = ""
			}

		case "up", "k":
			if m.state == "list" && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.state == "list" && m.cursor < len(m.notes)-1 {
				m.cursor++
			}

		case "t":
			if (m.state == "list" || m.state == "view") && len(m.notes) > 0 {
				m.state = "tags"
				m.tagEditMode = ""
				return m, nil
			}

		case "a":
			if m.state == "tags" && m.tagEditMode == "" {
				m.tagEditMode = "add"
				m.tagInput.Focus()
				m.tagInput.SetValue("")
				return m, nil
			}

		case "r":
			if m.state == "tags" && m.tagEditMode == "" && len(m.notes[m.cursor].Metadata.Tags) > 0 {
				m.tagEditMode = "remove"
				m.tagInput.Focus()
				m.tagInput.SetValue("")
				return m, nil
			}

		case "enter":
			switch m.state {
			case "list":
				if len(m.notes) > 0 {
					m.state = "view"
				}
			case "create_name":
				m.newNoteName = m.input.Value()
				if m.newNoteName != "" {
					m.state = "create"
					m.textarea.Focus()
					m.input.Reset()
				}
			}

		case "ctrl+s":
			switch m.state {
			case "create":
				if m.newNoteName != "" {
					err := m.notesApp.CreateNote(m.newNoteName, m.textarea.Value())
					if err != nil {
						m.err = err
					} else {
						m.notes = m.notesApp.ListAllNotes()
						m.state = "list"
						m.textarea.Reset()
						m.newNoteName = ""
					}
				}
			case "edit":
				if len(m.notes) > 0 {
					err := m.notesApp.UpdateNoteContent(m.notes[m.cursor].Path, m.textarea.Value())
					if err != nil {
						m.err = err
					} else {
						m.notes = m.notesApp.ListAllNotes()
						m.state = "list"
						m.textarea.Reset()
					}
				}
			}

		case "n":
			if m.state == "list" {
				m.state = "create_name"
				m.input.Focus()
				return m, nil
			}

		case "e":
			if (m.state == "list" || m.state == "view") && len(m.notes) > 0 {
				m.state = "edit"
				m.textarea.SetValue(m.notes[m.cursor].Content)
				m.textarea.Focus()
				return m, nil
			}

		case "d":
			if (m.state == "list" || m.state == "view") && len(m.notes) > 0 {
				err := m.notesApp.DeleteNote(m.notes[m.cursor].Path)
				if err != nil {
					m.err = err
				} else {
					m.notes = m.notesApp.ListAllNotes()
					if m.cursor >= len(m.notes) {
						m.cursor = len(m.notes) - 1
					}
					if m.state == "view" {
						m.state = "list"
					}
				}
				return m, nil
			}
		}
	}

	// Handle other input states
	switch m.state {
	case "create_name":
		m.input, cmd = m.input.Update(msg)
	case "create", "edit":
		m.textarea, cmd = m.textarea.Update(msg)
	}

	return m, cmd
}

// View implements tea.Model
func (m Model) View() string {
	var s strings.Builder

	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v\n\n", m.err)))
	}

	switch m.state {
	case "help":
		s.WriteString(titleStyle.Render("Help") + "\n\n")
		s.WriteString(helpStyle.Render(`
Commands:
  h, ?      - Show this help
  j, ↓      - Move down
  k, ↑      - Move up
  n         - Create new note
  e         - Edit selected note
  d         - Delete selected note
  t         - Manage tags
  space     - Toggle preview
  enter     - View note
  ctrl+s    - Save (in edit/create mode)
  esc       - Back/cancel
  q, ctrl+c - Quit application

Tag Management:
  a         - Add tags
  r         - Remove tags
  enter     - Confirm tags
  esc       - Cancel
`))

	case "create_name":
		s.WriteString(titleStyle.Render("New Note") + "\n\n")
		s.WriteString("Enter note name (press enter when done):\n")
		s.WriteString(inputStyle.Render(m.input.View()))

	case "create":
		s.WriteString(titleStyle.Render("New Note: "+m.newNoteName) + "\n\n")
		s.WriteString("Enter note content (ctrl+s to save, esc to cancel):\n")
		s.WriteString(textareaStyle.Render(m.textarea.View()))

	case "edit":
		if len(m.notes) > 0 {
			s.WriteString(titleStyle.Render("Editing: "+m.notes[m.cursor].Name) + "\n\n")
			s.WriteString(textareaStyle.Render(m.textarea.View()))
			s.WriteString("\n" + helpStyle.Render("Press ctrl+s to save, esc to cancel"))
		}

	case "view":
		if len(m.notes) > 0 {
			note := m.notes[m.cursor]
			s.WriteString(titleStyle.Render(note.Name))
			if len(note.Metadata.Tags) > 0 {
				s.WriteString(" " + tagStyle.Render(fmt.Sprintf("[%s]", strings.Join(note.Metadata.Tags, ", "))))
			}
			s.WriteString("\n\n")
			s.WriteString(note.Content)
			s.WriteString("\n\n" + helpStyle.Render("Press esc to go back"))
		}

	case "tags":
		if len(m.notes) > 0 {
			note := m.notes[m.cursor]
			s.WriteString(titleStyle.Render("Tags for: "+note.Name) + "\n\n")

			if len(note.Metadata.Tags) > 0 {
				s.WriteString("Current tags: " + tagStyle.Render(strings.Join(note.Metadata.Tags, ", ")) + "\n\n")
				s.WriteString(helpStyle.Render("Press 'a' to add tags, 'r' to remove tags, esc to go back") + "\n\n")
			} else {
				s.WriteString("No tags\n\n")
				s.WriteString(helpStyle.Render("Press 'a' to add tags, esc to go back") + "\n\n")
			}

			if m.tagEditMode == "add" {
				s.WriteString("Add tags (comma-separated):\n")
				s.WriteString(inputStyle.Render(m.tagInput.View()) + "\n")
				s.WriteString(helpStyle.Render("Press enter to confirm"))
			} else if m.tagEditMode == "remove" {
				s.WriteString("Remove tags (comma-separated):\n")
				s.WriteString(inputStyle.Render(m.tagInput.View()) + "\n")
				s.WriteString(helpStyle.Render("Press enter to confirm"))
			}
		}

	case "list":
		s.WriteString(titleStyle.Render("📝 Notes") + "\n\n")

		if len(m.notes) == 0 {
			s.WriteString(listStyle.Render("No notes found. Press 'n' to create one."))
		} else {
			var listContent strings.Builder
			for i, note := range m.notes {
				cursor := " "
				if m.cursor == i {
					cursor = ">"
				}

				noteText := fmt.Sprintf("%s %s",
					cursor,
					note.Name)

				if len(note.Metadata.Tags) > 0 {
					noteText += " " + tagStyle.Render(
						fmt.Sprintf("[%s]", strings.Join(note.Metadata.Tags, ", ")),
					)
				}

				if m.cursor == i {
					listContent.WriteString(selectedNoteStyle.Render(noteText))
				} else {
					listContent.WriteString(noteStyle.Render(noteText))
				}
				listContent.WriteString("\n")
			}
			s.WriteString(listStyle.Render(listContent.String()))

			// Add preview of selected note at the bottom if enabled
			if m.showPreview && len(m.notes) > 0 {
				note := m.notes[m.cursor]
				s.WriteString("\n" + previewTitleStyle.Render("Preview"))

				content := note.Content
				if len(content) > 200 {
					content = content[:200] + "..."
				}

				s.WriteString("\n" + previewStyle.Render(content))
			}
		}

		s.WriteString("\n" + helpStyle.Render("Press '?' for help, space to toggle preview"))
	}

	return mainStyle.Render(s.String())
}
