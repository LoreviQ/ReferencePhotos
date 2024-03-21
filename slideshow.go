package main

import (
	"image"
	"image/color"

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

func slideshow(event app.FrameEvent, ops *op.Ops, theme *material.Theme, ss slideshowWidgets) {
	gtx := app.NewContext(ops, event)
	modifyStateSlideshow()
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

func modifyStateSlideshow() {

}
