package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type model struct {
	help    help.Model
	list    list.Model
	spinner spinner.Model

	height int

	Mode ViewMode

	// custom position => position in list +1
	positions map[int]int

	currentPosition int
	positionSub     chan int

	Loading  bool
	Enabled  bool
	Quitting bool
}

func initialModel() *model {
	// list
	items := []list.Item{}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false) // Bugged, we can't update position of items

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#72EDB2"))

	// help
	h := help.New()
	h.ShowAll = true
	h.FullSeparator = "     "

	return &model{list: l, help: h, spinner: s, positions: make(map[int]int), positionSub: make(chan int), currentPosition: 1}
}

func (m model) Init() (tea.Model, tea.Cmd) {
	return m, tea.Batch(m.spinner.Tick, m.loadWindows(), waitForPositionChange(m.positionSub))
}

type updatePositionMsg int

func waitForPositionChange(sub chan int) tea.Cmd {
	return func() tea.Msg {
		return updatePositionMsg(<-sub)
	}
}

// focus window

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updatePositionMsg:
		m.currentPosition += int(msg)
		if m.currentPosition < 1 {
			m.currentPosition = max(len(m.positions), 1)
		} else if m.currentPosition > len(m.positions) {
			m.currentPosition = 1
		}
		if m.Enabled {
			itemPos := m.positions[m.currentPosition]
			item := m.list.Items()[itemPos]
			focusWindow(item.(ListEntry))
		}
		return m, waitForPositionChange(m.positionSub)

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
			m.positions = make(map[int]int)
			return m, m.loadWindows()

		case key.Matches(msg, keys.Quit):
			m.Quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.Delete):
			if len(m.list.Items()) == 0 {
				return m, nil
			}
			item := m.list.SelectedItem().(ListEntry)
			if item.Position != 0 {
				m.list.SetItem(m.list.GlobalIndex(), item.SetPosition(0))
				delete(m.positions, item.Position)
			}
			return m, nil

		case key.Matches(msg, keys.Edit):
			if len(m.list.Items()) == 0 {
				return m, nil
			}

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
					delete(m.positions, pos)
				}
			}

			idx := m.list.GlobalIndex()
			m.positions[pos] = idx + 1
			item := m.list.SelectedItem().(ListEntry)
			if item.Position != 0 {
				delete(m.positions, item.Position)
			}
			m.list.SetItem(idx, item.SetPosition(pos))
		case key.Matches(msg, keys.Enable):
			m.Enabled = !m.Enabled
			if m.Enabled {
				keys.Enable.SetHelp("e", "disable")
				titleStyle = titleStyle.Background(titleEnabledColor)
			} else {
				keys.Enable.SetHelp("e", "enable")
				titleStyle = titleStyle.Background(titleDisabledColor)
			}

		}

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-7)
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

	title := "Organizer"
	if !m.Enabled {
		title += " (disabled)"
	} else {
		title += " (enabled)"
	}
	titleView := titleStyle.Render(title)
	if m.Loading {
		titleView += " " + m.spinner.View() + descStyle.Render("Loading windows...")
	}

	curr := strconv.Itoa(m.currentPosition)
	maxPos := strconv.Itoa(len(m.positions))
	positionView := ""
	if len(m.positions) > 0 {
		positionView = titleStyle.Render("Position: " + curr + "/" + maxPos)
	}

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		titleView,
		positionView,
		"",
		m.list.View(),
		m.help.View(keys)))
}
