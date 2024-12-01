package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	help    help.Model
	list    list.Model
	spinner spinner.Model

	height int

	Mode ViewMode

	// custom position => position in list +1
	positions map[int]int

	Loading  bool
	Quitting bool
}

func initialModel() *model {
	// list
	items := []list.Item{
		ListEntry{title: "Finder", app: "oui"},
		ListEntry{title: "Chrome", app: "some"},
		ListEntry{title: "Terminal", app: "Terminal"},
	}

	delegate := newItemDelegate()

	l := list.New(items, delegate, 0, 0)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#72EDB2"))

	// help
	h := help.New()
	h.ShowAll = true

	return &model{list: l, help: h, spinner: s, positions: make(map[int]int)}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, keys.Reload):
			if m.Loading {
				return m, nil
			}
			m.list.SetItems(nil)
			return m, m.loadWindows()

		case key.Matches(msg, keys.Quit):
			m.Quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.Edit):
			item := m.list.SelectedItem().(ListEntry)
			pos, err := strconv.Atoi(msg.String())
			if err != nil {
				return m, nil
			}
			oldPos := m.positions[pos]
			if oldPos != 0 {
				idx := oldPos - 1
				oldItem := m.list.Items()[idx]
				if oldItem != nil {
					m.list.SetItem(idx, oldItem.(ListEntry).SetPosition(0))
				}
			}
			idx := m.list.Index()
			m.positions[pos] = idx + 1
			m.list.SetItem(idx, item.SetPosition(pos))
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-6)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case loadedWindowsMsg:
		m.list.SetItems(msg)
		m.Loading = false
		empty := len(msg) == 0
		keys.SetEmpty(empty)

		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func (m model) View() string {
	if m.Quitting {
		return "Bye!"
	}

	titleView := titleStyle.Render("Organizer")
	if m.Loading {
		titleView += " " + m.spinner.View() + descStyle.Render("Loading windows...")
	}

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		titleView,
		m.list.View(),
		m.help.View(keys)))
}
