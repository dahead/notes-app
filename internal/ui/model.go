package ui

import (
	"fmt"
	"strings"

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
	showPreview bool
}

func NewModel(notesApp *app.NotesApp) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter note name..."
	ti.Width = StandardWidth - StandardTextInputPadding

	ta := textarea.New()
	ta.Placeholder = "Enter note content..."
	ta.ShowLineNumbers = false
	ta.SetWidth(StandardWidth - StandardTextInputPadding)
	ta.SetHeight(18)

	tagInput := textinput.New()
	tagInput.Placeholder = "Enter tags (comma-separated)..."
	tagInput.Width = StandardWidth - StandardTextInputPadding

	return Model{
		notes:       notesApp.ListAllNotes(),
		input:       ti,
		textarea:    ta,
		tagInput:    tagInput,
		notesApp:    notesApp,
		selected:    make(map[int]struct{}),
		state:       "list",
		showPreview: false,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global shortcuts that work in any state
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			if m.state == "list" {
				return m, tea.Quit
			}
			m.state = "list"
			return m, nil
		}

		// State-specific shortcuts
		switch m.state {
		case "list":
			switch msg.String() {
			case "ctrl+h", "?":
				m.state = "help"
			case "ctrl+t", "t":
				if len(m.notes) > 0 {
					m.state = "tags"
					m.tagEditMode = ""
				}
			case "ctrl+n", "n":
				m.state = "create_name"
				m.input.Focus()
			case "ctrl+e", "e":
				if len(m.notes) > 0 {
					m.state = "edit"
					m.textarea.SetValue(m.notes[m.cursor].Content)
					m.textarea.Focus()
				}
			case "ctrl+d":
				if len(m.notes) > 0 {
					m.state = "confirm_delete"
				}

			case "enter":
				if len(m.notes) > 0 {
					m.state = "view"
				}
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.notes)-1 {
					m.cursor++
				}
			case " ":
				m.showPreview = !m.showPreview
			}

		case "view":
			switch msg.String() {
			case "ctrl+h", "?":
				m.state = "help"
			case "ctrl+t", "t":
				if len(m.notes) > 0 {
					m.state = "tags"
					m.tagEditMode = ""
				}
			case "ctrl+e", "e":
				if len(m.notes) > 0 {
					m.state = "edit"
					m.textarea.SetValue(m.notes[m.cursor].Content)
					m.textarea.Focus()
				}
			case "ctrl+d", "d":
				if len(m.notes) > 0 {
					m.state = "confirm_delete"
				}

			case "esc":
				m.state = "list"
			}

		case "create_name":
			switch msg.String() {
			case "enter":
				m.newNoteName = m.input.Value()
				if m.newNoteName != "" {
					if m.notesApp.NoteExists(m.newNoteName) {
						m.err = fmt.Errorf("note with name '%s' already exists", m.newNoteName)
						m.input.SetValue("")
						return m, nil
					}
					m.state = "create"
					m.textarea.Reset()
					m.textarea.Focus()
					m.input.Reset()
				}
			case "esc":
				m.state = "list"
				m.input.Reset()
				m.newNoteName = ""
			}
			m.input, cmd = m.input.Update(msg)

		case "create":
			switch msg.String() {
			case "ctrl+s":
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
			case "esc":
				m.state = "list"
				m.textarea.Reset()
				m.newNoteName = ""
			}
			m.textarea, cmd = m.textarea.Update(msg)

		case "edit":
			switch msg.String() {
			case "ctrl+s":
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
			case "esc":
				m.state = "list"
				m.textarea.Reset()
			}
			m.textarea, cmd = m.textarea.Update(msg)

		case "confirm_delete":
			switch msg.String() {
			case "y":
				err := m.notesApp.DeleteNote(m.notes[m.cursor].Path)
				if err != nil {
					m.err = err
				} else {
					m.notes = m.notesApp.ListAllNotes()
					if m.cursor >= len(m.notes) {
						m.cursor = len(m.notes) - 1
					}
					m.state = "list"
				}
			case "n", "esc":
				m.state = "list"
			}

		case "tags":
			if m.tagEditMode != "" {
				switch msg.String() {
				case "esc":
					m.tagInput.Reset()
					m.tagEditMode = ""
				case "enter":
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
				default:
					m.tagInput, cmd = m.tagInput.Update(msg)
				}
			} else {
				switch msg.String() {
				case "ctrl+a":
					m.tagEditMode = "add"
					m.tagInput.Focus()
					m.tagInput.SetValue("")
				case "ctrl+d":
					if len(m.notes[m.cursor].Metadata.Tags) > 0 {
						m.tagEditMode = "delete"
						m.tagInput.Focus()
						m.tagInput.SetValue("")
					}
				case "esc":
					m.state = "list"
					m.tagInput.Reset()
				}
			}

		case "help":
			switch msg.String() {
			case "esc":
				m.state = "list"
			}
		}

	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width - 4)
		return m, nil
	}

	return m, cmd
}

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
  ctrl+h, ? - Show this help
  j, ↓      - Move down
  k, ↑      - Move up
  ctrl+n    - Create new note
  ctrl+e    - Edit selected note
  ctrl+d    - Delete selected note
  ctrl+t    - Manage tags
  space     - Toggle preview
  enter     - View note
  ctrl+s    - Save (in edit/create mode)
  esc       - Back/cancel
  ctrl+q    - Quit application

Tag Management:
  ctrl+a    - Add tags
  ctrl+d    - Delete tags
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
	case "confirm_delete":
		if len(m.notes) > 0 {
			note := m.notes[m.cursor]
			s.WriteString(titleStyle.Render("Confirm Delete") + "\n\n")
			s.WriteString(fmt.Sprintf("Are you sure you want to delete note '%s'?\n\n", note.Name))
			s.WriteString(helpStyle.Render("Press 'y' to confirm, 'n' or 'esc' to cancel"))
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
				s.WriteString(helpStyle.Render("Press 'ctrl+a' to add tags, 'ctrl+d' to delete tags, esc to go back") + "\n\n")
			} else {
				s.WriteString("No tags\n\n")
				s.WriteString(helpStyle.Render("Press 'ctrl+a' to add tags, esc to go back") + "\n\n")
			}

			if m.tagEditMode == "add" {
				s.WriteString("Add tags (comma-separated):\n")
				s.WriteString(inputStyle.Render(m.tagInput.View()) + "\n")
				s.WriteString(helpStyle.Render("Press enter to confirm"))
			} else if m.tagEditMode == "delete" {
				s.WriteString("Delete tags (comma-separated):\n")
				s.WriteString(inputStyle.Render(m.tagInput.View()) + "\n")
				s.WriteString(helpStyle.Render("Press enter to confirm"))
			}
		}

	case "list":
		s.WriteString(titleStyle.Render("📝 Notes") + "\n\n")

		if len(m.notes) == 0 {
			s.WriteString(listStyle.Render("No notes found. Press 'ctrl+n' to create one."))
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
