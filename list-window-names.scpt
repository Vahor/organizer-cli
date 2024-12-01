tell application "System Events"
    set windowNames to {}
	repeat with currentApp in (application processes where visible is true)
        try
			if name of currentApp contains "Dofus" then
            set appWindows to windows of currentApp
				repeat with w in appWindows
					set end of windowNames to name of w & " (" & name of currentApp & ")"
				end repeat
			end if
        end try
    end repeat

	return windowNames
end tell
