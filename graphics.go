package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
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

type state struct {
	time      string
	directory string
	active    bool
}

var localState state

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
	localState = state{time: "30s", directory: "", active: false}

	//	My colours
	myColours = colours{
		bg:        color.NRGBA{41, 40, 45, 255},
		fg:        color.NRGBA{53, 54, 62, 255},
		text:      color.NRGBA{255, 255, 255, 255},
		highlight: color.NRGBA{60, 92, 115, 255},
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

	// Main event loop
	for {

		switch event := window.NextEvent().(type) {

		// Re-render app
		case app.FrameEvent:
			if localState.active {
				slideshow(event, &ops, theme)
			} else {
				landingPage(window, event, &ops, theme, lw)
			}

		// Exit app
		case app.DestroyEvent:
			return event.Err
		}
	}
}
