package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func createGUI(width, height int) {
	// Define Resolution
	if width == 0 || height == 0 {
		width, height = 800, 1250
	}

	// Create new window
	window := app.NewWindow(
		app.Title("Reference Photos"),
		app.Size(unit.Dp(width), unit.Dp(height)),
	)

	// Update window until exit
	if err := draw(window); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func draw(window *app.Window) error {
	var startButton widget.Clickable

	var ops op.Ops
	theme := material.NewTheme()

	// Main event loop
	for {
		switch event := window.NextEvent().(type) {

		// Re-render app
		case app.FrameEvent:
			landingPage(ops, startButton, theme, event)

		// Exit app
		case app.DestroyEvent:
			return event.Err
		}
	}
}

func landingPage(ops op.Ops, startButton widget.Clickable, theme *material.Theme, event app.FrameEvent) {
	gtx := app.NewContext(&ops, event)
	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				margins := layout.Inset{
					Top:    unit.Dp(5),
					Bottom: unit.Dp(5),
					Right:  unit.Dp(100),
					Left:   unit.Dp(100),
				}

				return margins.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						fmt.Printf("gtx.Constraints: %v\n", gtx.Constraints)
						button := material.Button(theme, &startButton, "Start")
						return button.Layout(gtx)
					},
				)

			},
		),
		layout.Rigid(
			layout.Spacer{Height: unit.Dp(25)}.Layout,
		),
	)
	event.Frame(gtx.Ops)
}
