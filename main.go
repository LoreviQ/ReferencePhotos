package main

import (
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/LoreviQ/ReferencePhotos/internal/resolution"
)

func main() {
	go func() {
		// Define Resolution
		width, height := 400, 600
		res := resolution.GetPrimary()
		if res != nil {
			width, height = res.Width, res.Height
		}

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
