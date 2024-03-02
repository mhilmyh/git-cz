package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mhilmyh/git-cz/pkg/list2"
	"github.com/mhilmyh/git-cz/pkg/model"
)

func main() {
	m := model.Model{
		TypeOfChangeList:  list.New(list2.NewListItemTypeOfChange(), list.NewDefaultDelegate(), 0, 0),
		ScopeOfChangeList: list.New(list2.NewListItemScopeOfChange(), list.NewDefaultDelegate(), 0, 0),
	}
	m.TypeOfChangeList.Title = "Type of change"
	m.ScopeOfChangeList.Title = "Scope of change"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
