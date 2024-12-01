package main

import (
	"log"

	hook "github.com/robotn/gohook"
)

type Callback func()

func InitHotkeys(m *model) {
	// register 0->9 hotkey and pageUp/Down
	log.Println("register hotkey")

	hook.Register(hook.KeyDown, []string{"num1", "shift"}, func(e hook.Event) {
		app.positionSub <- -1
	})
	hook.Register(hook.KeyDown, []string{"num2", "shift"}, func(e hook.Event) {
		app.positionSub <- +1
	})

	s := hook.Start()
	<-hook.Process(s)
}
