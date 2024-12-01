package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		item := m.SelectedItem().(ListEntry)
		if item.Position == 0 {
			keys.Delete.SetEnabled(false)
			keys.Edit.SetEnabled(true)
		} else {
			keys.Delete.SetEnabled(true)
			keys.Edit.SetEnabled(false)
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Delete):
				m.SetItem(m.Index(), item.SetPosition(0))
			}
		}

		return nil
	}

	return d
}
