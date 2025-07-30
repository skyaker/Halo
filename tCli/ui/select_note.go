package ui

import (
	"fmt"
	"halo/localstore"
	"halo/logger"
	"halo/models"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Notes     []models.NoteStruct
	cursor    int
	checked   map[string]bool
	paginator paginator.Model
	total     int
	quitting  bool
}

func newModel() model {
	notes := localstore.GetNotesLocally(1, 10)
	numOfNotes, err := localstore.GetNumberOfNotes()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("get number of notes")
	}

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).
		Render("•")
	p.InactiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).
		Render("•")
	p.SetTotalPages(len(notes))

	return model{
		Notes:     notes,
		cursor:    0,
		checked:   make(map[string]bool),
		paginator: p,
		total:     numOfNotes,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			if m.paginator.Page > 1 {
				m.paginator.PrevPage()
				m.Notes = localstore.GetNotesLocally(m.paginator.Page, m.paginator.PerPage)
			}
			m.paginator.PrevPage()
			m.cursor = 0
		case "right", "l":
			maxPage := (m.total + m.paginator.PerPage - 1) / m.paginator.PerPage
			if m.paginator.Page < maxPage {
				m.paginator.NextPage()
				m.Notes = localstore.GetNotesLocally(m.paginator.Page, m.paginator.PerPage)
				m.cursor = 0
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < m.paginator.ItemsOnPage(len(m.Notes))-1 {
				m.cursor++
			}
		case " ":
			start, _ := m.paginator.GetSliceBounds(len(m.Notes))
			id := m.Notes[start+m.cursor].Id.String()
			m.checked[id] = !m.checked[id]
		case "enter":
			// Удаление выбран
			var newNotes []models.NoteStruct
			for _, t := range m.Notes {
				id := t.Id.String()
				if !m.checked[id] {
					newNotes = append(newNotes, t)
				}
			}
			m.Notes = newNotes
			m.checked = make(map[string]bool)
			m.paginator.SetTotalPages(len(m.Notes))
			m.cursor = 0
		}
	}

	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "Exit...\n"
	}

	var b strings.Builder
	b.WriteString(
		"Manual (↑/↓ - navigation, Space - undo, ←/→ - page, Enter - delete, q - exit)\n\n",
	)

	start, end := m.paginator.GetSliceBounds(len(m.Notes))
	pageNotes := m.Notes[start:end]

	for i, t := range pageNotes {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		checked := " "
		if m.checked[t.Id.String()] {
			checked = "x"
		}
		done := "❌"
		if t.Completed {
			done = "✅"
		}
		line := fmt.Sprintf("%s [%s] %s %s\n", cursor, checked, t.Content, done)
		b.WriteString(line)
	}

	b.WriteString("\n" + m.paginator.View())
	b.WriteString("\n\n")

	return b.String()
}

func StartNoteSelector() error {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
