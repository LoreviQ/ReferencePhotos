package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ImageResult struct {
	Error error
	Image image.Image
}

type slideshowWidgets struct {
	currentImage *ImageResult
}

func slideshow(event app.FrameEvent, ops *op.Ops, theme *material.Theme, ss *slideshowWidgets) {
	gtx := app.NewContext(ops, event)
	modifyStateSlideshow(ss)
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

func modifyStateSlideshow(ss *slideshowWidgets) {
	fmt.Println(localState.progress)
	if localState.order == nil || len(localState.order) == 0 {
		getRandomOrder()
	}
	if ss.currentImage.Image == nil || localState.progress >= 1 {
		ss.getNextImage()
	}
}

func getRandomOrder() {
	files, err := os.ReadDir(localState.cfg.Directory)
	if err != nil {
		log.Fatal(err)
	}
	slideshowLength := len(files)
	if slideshowLength > 50 {
		slideshowLength = 50
	}
	order := rand.Perm(len(files))
	localState.order = make([]string, 0, slideshowLength)
	for i := 0; i < 50; i++ {
		localState.order = append(localState.order, files[order[i]].Name())
	}
}

func (ss *slideshowWidgets) getNextImage() error {
	filepath := localState.cfg.Directory + string(os.PathSeparator) + localState.order[0]
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	ss.currentImage.Image = img
	localState.order = localState.order[1:]
	localState.progress = 0
	return nil
}
