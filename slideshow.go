package main

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

func slideshow(event app.FrameEvent, ops *op.Ops, theme *material.Theme) {
	gtx := app.NewContext(ops, event)
	paint.Fill(gtx.Ops, color.NRGBA{0, 0, 0, 255})
	layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				bar := material.ProgressBar(theme, progress)
				return bar.Layout(gtx)
			},
		),
	)
	event.Frame(gtx.Ops)
}
