package main

import "github.com/charmbracelet/bubbles/v2/key"

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},     // first column
		{k.Quit, k.Reload}, // second column
		{k.Delete, k.Edit}, // 3 column
		{k.Enable},         // 4 column
	}
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Delete key.Binding
	Reload key.Binding
	Enable key.Binding
	Quit   key.Binding
	Edit   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Reload: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reload"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d", "backspace"),
		key.WithHelp("⌫/d", "delete"),
	),
	Edit: key.NewBinding(
		key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
		key.WithHelp("[0-9]", "edit"),
	),
	Enable: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "enable"),
	),
}

func (k keyMap) SetEmpty(empty bool) {
	k.Down.SetEnabled(!empty)
	k.Up.SetEnabled(!empty)
	k.Edit.SetEnabled(!empty)
	k.Delete.SetEnabled(!empty)
}
