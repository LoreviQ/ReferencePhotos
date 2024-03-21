package main

import (
	"image"
	"image/color"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"slices"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/image/webp"
)

type ImageResult struct {
	Error error
	Image image.Image
}

type slideshowWidgets struct {
	currentImage *ImageResult
}

func slideshow(window *app.Window, event app.FrameEvent, ops *op.Ops, theme *material.Theme, ss *slideshowWidgets) {
	gtx := app.NewContext(ops, event)
	modifyStateSlideshow(window, ss)
	paint.Fill(gtx.Ops, color.NRGBA{0, 0, 0, 255})
	layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if ss.currentImage.Error == nil && ss.currentImage.Image == nil {
				// Blank
				return layout.Dimensions{}
			} else if ss.currentImage.Error != nil {
				// Print Error
				return material.H6(theme, ss.currentImage.Error.Error()).Layout(gtx)
			}
			// Draw image
			return widget.Image{
				Src: paint.NewImageOp(ss.currentImage.Image),
				Fit: widget.Contain,
			}.Layout(gtx)
		}),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				bar := material.ProgressBar(theme, localState.progress)
				return bar.Layout(gtx)
			},
		),
	)
	event.Frame(gtx.Ops)
}

func modifyStateSlideshow(window *app.Window, ss *slideshowWidgets) {
	if localState.order == nil || len(localState.order) == 0 {
		getRandomOrder()
	}
	if ss.currentImage.Image == nil || localState.progress >= 1 {
		err := ss.getNextImage()
		if err != nil {
			log.Printf("Cannot open %v. Err: %v", localState.order[0], err)
			localState.order = localState.order[1:]
			window.Invalidate()
		}
	}
}

func getRandomOrder() {
	// Gets files in directory
	files, err := os.ReadDir(localState.cfg.Directory)
	if err != nil {
		log.Fatal(err)
	}
	// Filters out unsupported files
	var filename string
	supportedFiletypes := []string{"jpeg", ".jpg", ".png", "webp"}
	supportedFiles := make([]fs.DirEntry, 0, len(files))
	for _, file := range files {
		filename = file.Name()
		if slices.Contains(supportedFiletypes, filename[len(filename)-4:]) {
			supportedFiles = append(supportedFiles, file)
		}
	}
	// Generates random order
	slideshowLength := len(supportedFiles)
	if slideshowLength > 50 {
		slideshowLength = 50
	}
	order := rand.Perm(len(supportedFiles))
	localState.order = make([]string, 0, slideshowLength)
	for i := 0; i < 50; i++ {
		localState.order = append(localState.order, supportedFiles[order[i]].Name())
	}
}

func (ss *slideshowWidgets) getNextImage() error {
	log.Printf("Opening: %v\n", localState.order[0])
	filepath := localState.cfg.Directory + string(os.PathSeparator) + localState.order[0]
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	var img image.Image
	if localState.order[0][len(localState.order[0])-4:] == "webp" {
		img, err = webp.Decode(file)
	} else {
		img, _, err = image.Decode(file)
	}
	if err != nil {
		return err
	}
	ss.currentImage.Image = img
	localState.order = localState.order[1:]
	localState.progress = 0
	return nil
}
