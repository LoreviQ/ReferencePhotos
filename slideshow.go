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
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/webp"
)

type ImageResult struct {
	Error error
	Image image.Image
}

type slideshowWidgets struct {
	currentImage    *ImageResult
	leftButton      *widget.Clickable
	pauseButton     *widget.Clickable
	rightButton     *widget.Clickable
	exitButton      *widget.Clickable
	infoButton      *widget.Clickable
	folderButton    *widget.Clickable
	volumeButton    *widget.Clickable
	onTopButton     *widget.Clickable
	greyscaleButton *widget.Clickable
	timerButton     *widget.Clickable
}

func slideshow(window *app.Window, ev app.FrameEvent, ops *op.Ops, theme *material.Theme, ss *slideshowWidgets) {
	gtx := app.NewContext(ops, ev)
	checkClick(ops, ev.Source, gtx)
	modifyStateSlideshow(window, ss)
	paint.Fill(gtx.Ops, color.NRGBA{0, 0, 0, 255})
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1,
			layout.Spacer{}.Layout,
		),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1,
						layout.Spacer{}.Layout,
					),
					layout.Rigid(
						drawImage(ss.currentImage, theme),
					),
					layout.Flexed(1,
						layout.Spacer{}.Layout,
					),
				)
			},
		),
		layout.Flexed(1,
			layout.Spacer{}.Layout,
		),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				bar := material.ProgressBar(theme, localState.progress)
				return bar.Layout(gtx)
			},
		),
	)
	if localState.opacity > 0 {
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
			layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(5, layout.Spacer{}.Layout),
						slideshowImageButtons(ss.exitButton, theme),
						slideshowImageButtons(ss.infoButton, theme),
						slideshowImageButtons(ss.folderButton, theme),
						slideshowImageButtons(ss.volumeButton, theme),
						slideshowImageButtons(ss.onTopButton, theme),
						slideshowImageButtons(ss.greyscaleButton, theme),
						slideshowImageButtons(ss.timerButton, theme),
						layout.Flexed(5, layout.Spacer{}.Layout),
					)
				},
			),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		)
	}
	ev.Frame(gtx.Ops)
}

func modifyStateSlideshow(window *app.Window, ss *slideshowWidgets) {
	if localState.showButtons && localState.opacity < 100 {
		localState.opacity += 10
	}
	if !localState.showButtons && localState.opacity > 0 {
		localState.opacity -= 10
	}
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

func checkClick(ops *op.Ops, q input.Source, gtx layout.Context) {
	width, height := gtx.Constraints.Max.X, gtx.Constraints.Max.Y
	defer clip.Rect{
		Max: image.Pt(width, height),
	}.Push(ops).Pop()
	event.Op(ops, tag)
	for {
		ev, ok := q.Event(pointer.Filter{
			Target: tag,
			Kinds:  pointer.Release,
		})
		if !ok {
			break
		}

		if x, ok := ev.(pointer.Event); ok {
			switch x.Kind {
			case pointer.Release:
				localState.showButtons = !localState.showButtons
			}
		}
	}
}

func slideshowImageButtons(button *widget.Clickable, theme *material.Theme) layout.FlexChild {
	icon, err := widget.NewIcon(icons.ContentAdd)
	if err != nil {
		log.Print(err)
	}
	iconButton := material.IconButton(theme, button, icon, "Exit Button")
	iconButton.Background = color.NRGBA{0, 0, 0, localState.opacity}
	iconButton.Color = color.NRGBA{255, 255, 255, localState.opacity * 2}
	return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		return iconButton.Layout(gtx)
	})
}

func drawImage(image *ImageResult, theme *material.Theme) func(layout.Context) layout.Dimensions {
	return func(gtx layout.Context) layout.Dimensions {
		if image.Error == nil && image.Image == nil {
			// Blank
			return layout.Dimensions{}
		} else if image.Error != nil {
			// Print Error
			return material.H6(theme, image.Error.Error()).Layout(gtx)
		}
		// Draw image
		return widget.Image{
			Src: paint.NewImageOp(image.Image),
			Fit: widget.Contain,
		}.Layout(gtx)
	}
}
