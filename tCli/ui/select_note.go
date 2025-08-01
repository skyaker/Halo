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
	Notes       []models.NoteStruct
	PickedNotes []models.NoteStruct
	cursor      int
	checked     map[string]bool
	paginator   paginator.Model
	total       int
	quitting    bool
}

func newModel() model {
	notes := localstore.GetNotesLocally(0, 10)
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
	p.SetTotalPages(numOfNotes)

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
			if m.paginator.Page > 0 {
				// -1 because of lib auto page change only after switch case
				m.Notes = localstore.GetNotesLocally(m.paginator.Page-1, m.paginator.PerPage)
				m.cursor = 0
			}
		case "right", "l":
			if m.paginator.Page < m.paginator.TotalPages-1 {
				// +1 because of lib auto page change only after switch case
				m.Notes = localstore.GetNotesLocally(m.paginator.Page+1, m.paginator.PerPage)
				m.cursor = 0
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.Notes)-1 {
				m.cursor++
			}
		case " ":
			id := m.Notes[m.cursor].Id.String()
			m.checked[id] = !m.checked[id]
		case "enter":
			for id := range m.checked {
				if m.checked[id] {
					err := localstore.DeleteNoteLocally(id)
					if err != nil {
						logger.Logger.Error().Err(err).Msg("local delete note")
					}
					m.total--
				}
			}

			m.paginator.SetTotalPages(m.total)
			if m.paginator.Page <= m.paginator.TotalPages-1 {
				m.Notes = localstore.GetNotesLocally(m.paginator.Page, m.paginator.PerPage)
			} else {
				m.paginator.OnLastPage()
				m.Notes = localstore.GetNotesLocally(m.paginator.Page, m.paginator.PerPage)
			}
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
		"Manual\n· ↑(k)/↓(j) - navigation\n· Space - select\n· ←(h)/→(l) - page\n· Enter - delete\n· q - exit\n\n",
	)

	pageNotes := m.Notes

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
