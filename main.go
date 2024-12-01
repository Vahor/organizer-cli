package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04B575")).
				Render
)

type ViewMode int

const (
	ViewHome ViewMode = iota
)

type loadedWindowsMsg []list.Item

func (m *model) loadWindows() tea.Cmd {
	m.Loading = true
	return func() tea.Msg {
		windows, err := listWindowNames()
		if err != nil {
			fmt.Println("Error listing windows:", err)
			m.Loading = false
			return nil
		}
		entries := formatToListEntry(windows)
		return loadedWindowsMsg(entries)
	}
}

var app *model

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	app = initialModel()
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
