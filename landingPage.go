package main

import (
	"fmt"
	"image/color"
	"log"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ncruces/zenity"
)

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

type dimensions struct {
	top    int
	bottom int
	right  int
	left   int
}

func landingPage(window *app.Window, event app.FrameEvent, ops *op.Ops, theme *material.Theme, lw landingPageWidgets) {
	gtx := app.NewContext(ops, event)
	modifyStateLandingPage(window, gtx, lw)
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
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(5,
					func(gtx layout.Context) layout.Dimensions {
						title := material.Label(theme, unit.Sp(60), "SLIDESHOW!")
						title.Alignment = text.Middle
						title.Color = myColours.text
						return title.Layout(gtx)
					},
				),
				layout.Flexed(2,
					func(gtx layout.Context) layout.Dimensions {
						title := material.Label(theme, unit.Sp(15), "v1.0 Oliver Jay, 2024")
						title.Alignment = text.End
						title.Color = myColours.text
						return title.Layout(gtx)
					},
				),
				layout.Flexed(1,
					layout.Spacer{Height: unit.Dp(d.bottom)}.Layout,
				),
				layout.Flexed(5,
					func(gtx layout.Context) layout.Dimensions {
						return myButton(gtx, theme, lw.sourceButton, "Select an image source", myColours.highlight)
					},
				),
				layout.Flexed(0.5,
					layout.Spacer{Height: unit.Dp(d.bottom)}.Layout,
				),
				layout.Flexed(2,
					func(gtx layout.Context) layout.Dimensions {
						dirSource := material.Label(theme, unit.Sp(15), fmt.Sprintf("Source: %v", localState.cfg.Directory))
						dirSource.Alignment = text.Start
						dirSource.Color = myColours.text
						return dirSource.Layout(gtx)
					},
				),
				layout.Flexed(2,
					layout.Spacer{Height: unit.Dp(d.bottom)}.Layout,
				),
				layout.Flexed(2,
					func(gtx layout.Context) layout.Dimensions {
						subtitle := material.Label(theme, unit.Sp(20), "Time Between Images")
						subtitle.Alignment = text.Middle
						subtitle.Color = myColours.text
						return subtitle.Layout(gtx)
					},
				),
				layout.Flexed(0.5,
					layout.Spacer{Height: unit.Dp(d.bottom)}.Layout,
				),
				layout.Flexed(5,
					func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{}.Layout(gtx,
							timerToggles("30s", lw.timeButton30s, theme),
							timerToggles("45s", lw.timeButton45s, theme),
							timerToggles("1m", lw.timeButton1m, theme),
							timerToggles("2m", lw.timeButton2m, theme),
							timerToggles("5m", lw.timeButton5m, theme),
							timerToggles("10m", lw.timeButton10m, theme),
						)
					},
				),
				layout.Flexed(1,
					layout.Spacer{Height: unit.Dp(d.bottom)}.Layout,
				),
				layout.Flexed(5,
					func(gtx layout.Context) layout.Dimensions {
						return myButton(gtx, theme, lw.startButton, "Start", myColours.highlight)
					},
				),
			)
		},
	)
	event.Frame(gtx.Ops)
}

func myButton(gtx layout.Context, theme *material.Theme, widget *widget.Clickable, text string, colour color.NRGBA) layout.Dimensions {
	button := material.Button(theme, widget, text)
	button.CornerRadius = unit.Dp(0)
	button.Inset = layout.UniformInset(unit.Dp(2))
	button.Background = colour
	return button.Layout(gtx)
}

func timerToggles(time string, widget *widget.Clickable, theme *material.Theme) layout.FlexChild {
	buttonColour := myColours.fg
	if localState.time == time {
		buttonColour = myColours.highlight
	}
	return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		return myButton(gtx, theme, widget, time, buttonColour)
	})
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

func modifyStateLandingPage(window *app.Window, gtx layout.Context, lw landingPageWidgets) {
	// Start Slideshow
	if lw.startButton.Clicked(gtx) && localState.cfg.Directory != "" {
		localState.active = !localState.active
		progressIncrementer := getProgressIncrementer(localState.time)
		go incrementProgress(window, progressIncrementer, localState.exit)
	}
	// Time
	if lw.timeButton30s.Clicked(gtx) {
		localState.time = "30s"
	}
	if lw.timeButton45s.Clicked(gtx) {
		localState.time = "45s"
	}
	if lw.timeButton1m.Clicked(gtx) {
		localState.time = "1m"
	}
	if lw.timeButton2m.Clicked(gtx) {
		localState.time = "2m"
	}
	if lw.timeButton5m.Clicked(gtx) {
		localState.time = "5m"
	}
	if lw.timeButton10m.Clicked(gtx) {
		localState.time = "10m"
	}

	// Get Directory Source
	if lw.sourceButton.Clicked(gtx) {
		dir, err := zenity.SelectFile(
			zenity.Filename(""),
			zenity.Directory())
		if err != nil {
			log.Print(err)
		}
		localState.cfg.Directory = dir
		localState.cfg.writeCFG()
	}
}
