package model

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mhilmyh/git-cz/pkg/list2"
)

var selectionStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	PaddingTop(1)
var titleLineStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("62")).
	Foreground(lipgloss.Color("230")).
	Padding(0, 1)
var grayLineStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})
var textInputStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	PaddingLeft(2).
	PaddingRight(2).
	PaddingTop(1).
	PaddingBottom(2)

var state = 0

type Model struct {
	TypeOfChangeList      list.Model
	SelectedTypeOfChange  list2.Item
	ScopeOfChangeList     list.Model
	SelectedScopeOfChange list2.Item
	TitleCommitInput      textinput.Model
	WrittenTitleCommit    string
	FinalPromptInput      textinput.Model
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
			} else if state == 2 {
				m.WrittenTitleCommit = m.TitleCommitInput.Value()
			}
			state += 1
			if state == 1 {
				m.ScopeOfChangeList.SetSize(m.TypeOfChangeList.Width(), m.TypeOfChangeList.Height())
			} else if state == 2 {
				m.TitleCommitInput.Width = 64
			} else if state == 4 {
				cmd := exec.Command("git", "commit", "-m", m.FullCommitMessage())
				if cmd != nil {
					cmd.Run()
				}
				return m, tea.Quit
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := selectionStyle.GetFrameSize()
		if state == 0 {
			m.TypeOfChangeList.SetSize(msg.Width-h, msg.Height-v)
		} else if state == 1 {
			m.ScopeOfChangeList.SetSize(msg.Width-h, msg.Height-v)
		} else if state == 2 {
			m.TitleCommitInput.Width = min(64, max(16, msg.Width*50/100))
		}
	}

	var cmd tea.Cmd
	if state == 0 {
		m.TypeOfChangeList, cmd = m.TypeOfChangeList.Update(msg)
	} else if state == 1 {
		m.ScopeOfChangeList, cmd = m.ScopeOfChangeList.Update(msg)
	} else if state == 2 {
		m.TitleCommitInput, cmd = m.TitleCommitInput.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if state == 0 {
		return selectionStyle.Render(m.TypeOfChangeList.View())
	} else if state == 1 {
		return selectionStyle.Render(m.ScopeOfChangeList.View())
	} else if state == 2 {
		renderedTitle := titleLineStyle.Render("Title of commit " + strconv.FormatInt(int64(m.TitleCommitInput.Width), 10))
		titleCommitInput := m.TitleCommitInput.View()
		return textInputStyle.Render(renderedTitle + "\n\n" + titleCommitInput)
	} else if state == 3 {
		s := fmt.Sprintf("Type of change  : %s\n", m.SelectedTypeOfChange.Title())
		s += fmt.Sprintf("Scope of change : %s\n", m.SelectedScopeOfChange.Title())
		s += fmt.Sprintf("Title of commit : %s\n", m.WrittenTitleCommit)
		s += grayLineStyle.Render(m.FullCommitMessage())
		s += "\n\nDo you want to commit? "
		s += m.FinalPromptInput.View()
		return textInputStyle.Render(s)
	}
	return "..."
}

func (m Model) FullCommitMessage() string {
	typeOfChange := m.SelectedTypeOfChange.Code()
	scope := m.SelectedScopeOfChange.Code()
	var sb strings.Builder

	if typeOfChange != "" {
		sb.WriteString(typeOfChange)
	}
	if scope != "" {
		if sb.Len() > 0 {
			sb.WriteString("(" + scope + ")")
		} else {
			sb.WriteString(scope)
		}
	}
	if m.WrittenTitleCommit != "" {
		if sb.Len() > 0 {
			sb.WriteString(": " + m.WrittenTitleCommit)
		} else {
			sb.WriteString(m.WrittenTitleCommit)
		}

	}
	return fmt.Sprintf("\n%s(%s): %s\n", typeOfChange, scope, m.WrittenTitleCommit)
}
