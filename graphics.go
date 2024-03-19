package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
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

type landingPageWidgets struct {
	startButton   *widget.Clickable
	sourceButton  *widget.Clickable
	timeButton30s *widget.Clickable
	timeButton45s *widget.Clickable
	timeButton1m  *widget.Clickable
	timeButton2m  *widget.Clickable
	timeButton5m  *widget.Clickable
	timeButton10m *widget.Clickable
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
	var ops op.Ops
	theme := material.NewTheme()

	// Landing Page Widgets
	lw := landingPageWidgets{
		startButton:   &widget.Clickable{},
		sourceButton:  &widget.Clickable{},
		timeButton30s: &widget.Clickable{},
		timeButton45s: &widget.Clickable{},
		timeButton1m:  &widget.Clickable{},
		timeButton2m:  &widget.Clickable{},
		timeButton5m:  &widget.Clickable{},
		timeButton10m: &widget.Clickable{},
	}

	// Slideshow Widgets
	var progressIncrementer chan float32
	go incrementProgress(window, progressIncrementer)

	// Main event loop
	for {
		switch event := window.NextEvent().(type) {

		// Re-render app
		case app.FrameEvent:
			landingPage(event, &ops, theme, lw)

		// Exit app
		case app.DestroyEvent:
			return event.Err
		}
	}
}

func landingPage(event app.FrameEvent, ops *op.Ops, theme *material.Theme, lw landingPageWidgets) {
	gtx := app.NewContext(ops, event)
	modifyState(gtx, lw)
	d := getDimensions(gtx)
	// Whitespace
	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		// Title
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				title := material.Label(theme, unit.Sp(40), "SLIDESHOW")
				title.Alignment = text.Middle
				return title.Layout(gtx)
			},
		),
		// Image Source
		middleAlign(d,
			func(gtx layout.Context) layout.Dimensions {
				button := material.Button(theme, lw.sourceButton, "Select an image source")
				return button.Layout(gtx)
			},
		),
		// Timer
		middleAlign(d,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						button := material.Button(theme, lw.timeButton30s, "30s")
						return button.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						button := material.Button(theme, lw.timeButton45s, "45s")
						return button.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						button := material.Button(theme, lw.timeButton1m, "1m")
						return button.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						button := material.Button(theme, lw.timeButton2m, "2m")
						return button.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						button := material.Button(theme, lw.timeButton5m, "5m")
						return button.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						button := material.Button(theme, lw.timeButton10m, "10m")
						return button.Layout(gtx)
					}),
				)
			},
		),

		// Start
		middleAlign(d,
			func(gtx layout.Context) layout.Dimensions {
				var text string
				if !active {
					text = "Start"
				} else {
					text = "Stop"
				}
				button := material.Button(theme, lw.startButton, text)
				return button.Layout(gtx)
			},
		),
		// Whitespace
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

func middleAlign(d dimensions, element func(layout.Context) layout.Dimensions) layout.FlexChild {
	return layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			margins := layout.Inset{
				Top:    unit.Dp(5),
				Bottom: unit.Dp(5),
				Right:  unit.Dp(d.right),
				Left:   unit.Dp(d.left),
			}
			return margins.Layout(gtx, element)
		},
	)
}

func modifyState(gtx layout.Context, lw landingPageWidgets) {
	if lw.startButton.Clicked(gtx) {
		active = !active
	}
}
