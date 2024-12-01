tell application "System Events"
	tell process "Finder"
		repeat with w in windows
		if name of w contains "docs" then
			set frontmost to true
			return true
		end if
		end repeat
	end tell
end tell
