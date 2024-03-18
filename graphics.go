package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var progress float32
var progressIncrementer chan float32
var clicked bool

type dimensions struct {
	top    int
	bottom int
	right  int
	left   int
}

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
			landingPage(event, &ops, theme, &startButton)

		// Exit app
		case app.DestroyEvent:
			return event.Err
		}
	}
}

func landingPage(event app.FrameEvent, ops *op.Ops, theme *material.Theme, startButton *widget.Clickable) {
	gtx := app.NewContext(ops, event)
	if startButton.Clicked(gtx) {
		clicked = !clicked
	}
	d := getDimensions(gtx)
	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				bar := material.ProgressBar(theme, progress)
				return bar.Layout(gtx)
			},
		),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				margins := layout.Inset{
					Top:    unit.Dp(5),
					Bottom: unit.Dp(5),
					Right:  unit.Dp(d.right),
					Left:   unit.Dp(d.left),
				}

				return margins.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						var text string
						if !clicked {
							text = "Start"
						} else {
							text = "Stop"
						}
						button := material.Button(theme, startButton, text)
						return button.Layout(gtx)
					},
				)

			},
		),
		layout.Rigid(
			layout.Spacer{Height: unit.Dp(d.bottom)}.Layout,
		),
	)
	event.Frame(gtx.Ops)
}

func getDimensions(gtx layout.Context) dimensions {
	const interactableHeight int = 50
	const interactableWidth int = 500

	var d dimensions
	width, height := gtx.Constraints.Max.X, gtx.Constraints.Max.Y

	if height > interactableHeight {
		d.top = int((height - interactableHeight) / 2)
		d.bottom = (height - interactableHeight) - d.top
	} else {
		d.top, d.bottom = 0, 0
	}

	if width > interactableWidth {
		d.left = int((width - interactableWidth) / 2)
		d.right = (width - interactableWidth) - d.left
	} else {
		d.left, d.right = 0, 0
	}
	return d
}
