package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type colours struct {
	bg        color.NRGBA
	fg        color.NRGBA
	text      color.NRGBA
	highlight color.NRGBA
}

var myColours colours

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

	//	My colours
	myColours = colours{
		bg:        color.NRGBA{41, 40, 45, 1},
		fg:        color.NRGBA{53, 54, 62, 1},
		text:      color.NRGBA{255, 255, 255, 1},
		highlight: color.NRGBA{60, 92, 115, 1},
	}

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
			landingPage(event, &ops, theme, myColours, lw)

		// Exit app
		case app.DestroyEvent:
			return event.Err
		}
	}
}

func landingPage(event app.FrameEvent, ops *op.Ops, theme *material.Theme, myColours colours, lw landingPageWidgets) {
	gtx := app.NewContext(ops, event)
	modifyState(gtx, lw)
	paint.Fill(gtx.Ops, myColours.bg)
	d := getDimensions(gtx, 400, 500)
	margins := layout.Inset{
		Top:    unit.Dp(d.top),
		Bottom: unit.Dp(d.bottom),
		Right:  unit.Dp(d.right),
		Left:   unit.Dp(d.left),
	}
	margins.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Flexed(5,
					func(gtx layout.Context) layout.Dimensions {
						title := material.Label(theme, unit.Sp(40), "SLIDESHOW")
						title.Alignment = text.Middle
						return title.Layout(gtx)
					},
				),
				layout.Flexed(1,
					layout.Spacer{Height: unit.Dp(d.bottom)}.Layout,
				),
				layout.Flexed(5,
					func(gtx layout.Context) layout.Dimensions {
						return myButton(gtx, theme, lw.sourceButton, "Select an image source")
					},
				),
				layout.Flexed(1,
					layout.Spacer{Height: unit.Dp(d.bottom)}.Layout,
				),
				layout.Flexed(5,
					func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Horizontal,
						}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return myButton(gtx, theme, lw.timeButton30s, "30s")
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return myButton(gtx, theme, lw.timeButton45s, "45s")
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return myButton(gtx, theme, lw.timeButton1m, "1m")
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return myButton(gtx, theme, lw.timeButton2m, "2m")
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return myButton(gtx, theme, lw.timeButton5m, "5m")
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return myButton(gtx, theme, lw.timeButton10m, "10m")
							}),
						)
					},
				),
				layout.Flexed(1,
					layout.Spacer{Height: unit.Dp(d.bottom)}.Layout,
				),
				layout.Flexed(5,
					func(gtx layout.Context) layout.Dimensions {
						var text string
						if !active {
							text = "Start"
						} else {
							text = "Stop"
						}
						return myButton(gtx, theme, lw.startButton, text)
					},
				),
			)
		},
	)
	event.Frame(gtx.Ops)
}

func myButton(gtx layout.Context, theme *material.Theme, widget *widget.Clickable, text string) layout.Dimensions {
	button := material.Button(theme, widget, text)
	button.CornerRadius = unit.Dp(0)
	button.Inset = layout.UniformInset(unit.Dp(2))
	return button.Layout(gtx)
}

func getDimensions(gtx layout.Context, interactableHeight, interactableWidth int) dimensions {
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

func modifyState(gtx layout.Context, lw landingPageWidgets) {
	if lw.startButton.Clicked(gtx) {
		active = !active
	}
}
