package model

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mhilmyh/git-cz/pkg/list2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var state = 0

type Model struct {
	TypeOfChangeList      list.Model
	SelectedTypeOfChange  list2.Item
	ScopeOfChangeList     list.Model
	SelectedScopeOfChange list2.Item
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if state == 0 {
				i, ok := m.TypeOfChangeList.SelectedItem().(list2.Item)
				if ok {
					m.SelectedTypeOfChange = i
				}
			} else if state == 1 {
				i, ok := m.ScopeOfChangeList.SelectedItem().(list2.Item)
				if ok {
					m.SelectedScopeOfChange = i
				}
			}
			state += 1
			if state == 1 {
				m.ScopeOfChangeList.SetSize(m.TypeOfChangeList.Width(), m.TypeOfChangeList.Height())
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		if state == 0 {
			m.TypeOfChangeList.SetSize(msg.Width-h, msg.Height-v)
		} else if state == 1 {
			m.ScopeOfChangeList.SetSize(msg.Width-h, msg.Height-v)
		}
	}

	var cmd tea.Cmd
	if state == 0 {
		m.TypeOfChangeList, cmd = m.TypeOfChangeList.Update(msg)
	} else if state == 1 {
		m.ScopeOfChangeList, cmd = m.ScopeOfChangeList.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if state == 0 {
		return docStyle.Render(m.TypeOfChangeList.View())
	} else if state == 1 {
		return docStyle.Render(m.ScopeOfChangeList.View())
	}
	return "---"
}
