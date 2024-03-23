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
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type colours struct {
	bg        color.NRGBA
	fg        color.NRGBA
	text      color.NRGBA
	highlight color.NRGBA
}

var tag = new(bool)

var myColours colours

type state struct {
	cfg         config
	time        string
	active      bool
	paused      bool
	showButtons bool
	progress    float32
	opacity     uint8
	order       []string
	exit        chan bool
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
	// Initialising State
	localState = state{
		cfg:         InitialiseConfig("./config.json"),
		time:        "30s",
		active:      false,
		paused:      false,
		showButtons: false,
		progress:    0,
		opacity:     0,
		order:       nil,
		exit:        make(chan bool),
	}

	//	My colours
	myColours = colours{
		bg:        color.NRGBA{41, 40, 45, 255},
		fg:        color.NRGBA{53, 54, 62, 255},
		text:      color.NRGBA{255, 255, 255, 255},
		highlight: color.NRGBA{63, 81, 182, 255},
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
	ss := slideshowWidgets{
		currentImage: &ImageResult{
			Error: nil,
			Image: nil,
		},
		leftButton: &iconButton{
			button: &widget.Clickable{},
			label:  "Left Arrow",
		},
		pauseButton: &iconButton{
			button: &widget.Clickable{},
			label:  "Pause",
		},
		rightButton: &iconButton{
			button: &widget.Clickable{},
			label:  "Right Arrow",
		},
		exitButton: &iconButton{
			button: &widget.Clickable{},
			label:  "Exit",
		},
		infoButton: &iconButton{
			button: &widget.Clickable{},
			label:  "Info",
		},
		folderButton: &iconButton{
			button: &widget.Clickable{},
			label:  "Folder",
		},
		volumeButton: &iconButton{
			button: &widget.Clickable{},
			label:  "Volume",
		},
		onTopButton: &iconButton{
			button: &widget.Clickable{},
			label:  "Always on Top",
		},
		greyscaleButton: &iconButton{
			button: &widget.Clickable{},
			label:  "Greyscale",
		},
		timerButton: &iconButton{
			button: &widget.Clickable{},
			label:  "Time Button",
		},
	}
	ss.exitButton.icon, _ = widget.NewIcon(icons.NavigationCancel)
	ss.infoButton.icon, _ = widget.NewIcon(icons.ActionInfo)
	ss.folderButton.icon, _ = widget.NewIcon(icons.FileFolderOpen)
	ss.volumeButton.icon, _ = widget.NewIcon(icons.AVVolumeUp)
	ss.onTopButton.icon, _ = widget.NewIcon(icons.AVLibraryAdd)
	ss.greyscaleButton.icon, _ = widget.NewIcon(icons.ImageColorize)
	ss.timerButton.icon, _ = widget.NewIcon(icons.ActionHourglassEmpty)

	// Main event loop
	for {
		switch ev := window.NextEvent().(type) {

		// Re-render app
		case app.FrameEvent:
			if localState.active {
				slideshow(window, ev, &ops, theme, &ss)
			} else {
				landingPage(window, ev, &ops, theme, lw)
			}

		// Exit app
		case app.DestroyEvent:
			return ev.Err
		}

	}
}
