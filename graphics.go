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
			landingPage(ops, startButton, theme, event)

		// Exit app
		case app.DestroyEvent:
			return event.Err
		}
	}
}

func landingPage(ops op.Ops, startButton widget.Clickable, theme *material.Theme, event app.FrameEvent) {
	gtx := app.NewContext(&ops, event)
	d := getDimensions(gtx)
	fmt.Println(d)
	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
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
						button := material.Button(theme, &startButton, "Start")
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
	const actionHeight int = 50
	const actionWidth int = 500

	var d dimensions
	fmt.Println(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	width, height := gtx.Constraints.Max.X, gtx.Constraints.Max.Y

	if height > actionHeight {
		d.top = int((height - actionHeight) / 2)
		d.bottom = (height - actionHeight) - d.top
	} else {
		d.top, d.bottom = 0, 0
	}

	if width > actionWidth {
		d.left = int((width - actionWidth) / 2)
		d.right = (width - actionWidth) - d.left
	} else {
		d.left, d.right = 0, 0
	}
	return d
}
