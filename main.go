package main

import (
	"gioui.org/app"
	"gioui.org/unit"
)

func main() {
	go func() {
		// Define Resolution
		width, height := 800, 1250

		// create new window
		w := app.NewWindow(
			app.Title("Reference Photos"),
			app.Size(unit.Dp(width), unit.Dp(height)),
		)

		// listen for events in the window
		for {
			w.NextEvent()
		}
	}()
	app.Main()
}
