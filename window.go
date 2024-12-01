package main

import "strconv"

type ListEntry struct {
	title    string
	app      string
	Position int
}

// implement the list.Item interface
func (t ListEntry) FilterValue() string {
	return t.title + " " + t.app
}

func (t ListEntry) Title() string {
	if t.Position > 0 {
		return strconv.Itoa(t.Position) + " - " + t.title
	}
	return t.title
}

func (t ListEntry) Description() string {
	return t.app
}

func (t ListEntry) SetPosition(pos int) ListEntry {
	clone := t
	clone.Position = pos
	return clone
}
