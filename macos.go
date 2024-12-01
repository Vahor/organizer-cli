package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/v2/list"
)

func listWindowNames() ([]string, error) {
	script := `
tell application "System Events"
    set windowNames to {}
    repeat with currentApp in (every process)
        try
            set appWindows to windows of currentApp
            repeat with w in appWindows
                set end of windowNames to name of w & " (" & name of currentApp & ")"
            end repeat
        end try
    end repeat

	return windowNames
end tell
`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split the output into individual window names
	windows := strings.Split(strings.TrimSpace(string(output)), ", ")
	return windows, nil
}

func formatToListEntry(windows []string) []list.Item {
	entries := []ListEntry{}
	r := regexp.MustCompile(`^(.*) \((.*)\)$`)
	for _, w := range windows {
		parts := r.FindStringSubmatch(w)
		if len(parts) != 3 {
			continue
		}
		entries = append(entries, ListEntry{
			title: parts[1], app: parts[2]})
	}

	// sort by app name
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].app) < strings.ToLower(entries[j].app)
	})

	// TODO: I don't know how to cast
	listEntries := []list.Item{}
	for _, e := range entries {
		listEntries = append(listEntries, e)
	}
	return listEntries
}

func focusWindow(entry ListEntry) error {
	script := fmt.Sprintf(`
tell application "System Events"
	tell process "%s"
		repeat with w in windows
			try
				if name of w contains "%s" then
					set frontmost to true
					return true
				end if
			end try
		end repeat
	end tell
end tell`, entry.app, entry.title)
	cmd := exec.Command("osascript", "-e", script)
	log.Println("Focusing window:", entry)
	_, err := cmd.Output()
	if err != nil {
		log.Println("Error focusing window:", err)
		return err
	}
	return nil
}
