package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
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
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/image/webp"
)

type ImageResult struct {
	Error    error
	Image    image.Image
	Filename string
	Size     image.Point // Pixels
	Filesize int64       // Bytes
}

type iconButton struct {
	icon   *widget.Icon
	button *widget.Clickable
	active bool
	label  string
}

type slideshowWidgets struct {
	currentImage    *ImageResult
	leftButton      *iconButton
	pauseButton     *iconButton
	rightButton     *iconButton
	exitButton      *iconButton
	infoButton      *iconButton
	folderButton    *iconButton
	volumeButton    *iconButton
	onTopButton     *iconButton
	greyscaleButton *iconButton
	timerButton     *iconButton
}

func slideshow(window *app.Window, ev app.FrameEvent, ops *op.Ops, theme *material.Theme, ss *slideshowWidgets) {
	gtx := app.NewContext(ops, ev)
	checkClick(ops, ev.Source, gtx)
	modifyStateSlideshow(window, gtx, ss)
	paint.Fill(gtx.Ops, color.NRGBA{0, 0, 0, 255})
	// Image
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
	)
	// Progress Bar
	if ss.timerButton.active {
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
			layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {
					bar := material.ProgressBar(theme, localState.progressBar.progress)
					bar.Height = unit.Dp(10)
					return bar.Layout(gtx)
				},
			),
		)
	}
	// Buttons
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
	// File Data
	if ss.infoButton.active {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Flexed(1,
				func(gtx layout.Context) layout.Dimensions {
					title := material.Label(theme, unit.Sp(20),
						fmt.Sprintf("%v\n%dx%d px - %d KB",
							ss.currentImage.Filename,
							ss.currentImage.Size.X,
							ss.currentImage.Size.Y,
							ss.currentImage.Filesize/1024,
						),
					)
					title.Alignment = text.Middle
					title.Color = myColours.text
					return title.Layout(gtx)
				},
			),
		)
	}
	ev.Frame(gtx.Ops)
}

func modifyStateSlideshow(window *app.Window, gtx layout.Context, ss *slideshowWidgets) {
	// Button Opacity
	var speed uint8 = 30
	if localState.showButtons && localState.opacity < 255 {
		if localState.opacity <= 255-speed {
			localState.opacity += speed
		} else {
			localState.opacity = 255
		}
	}
	if !localState.showButtons && localState.opacity > 0 {
		if localState.opacity >= speed {
			localState.opacity -= speed
		} else {
			localState.opacity = 0
		}
	}

	// Order
	if localState.order == nil || len(localState.order) == 0 {
		getRandomOrder()
	}

	// Current Image
	if ss.currentImage.Image == nil || localState.progressBar.progress >= 1 {
		err := ss.getNextImage()
		localState.progressBar.sounds = 0
		if err != nil {
			log.Printf("Cannot open %v. Err: %v", localState.order[0], err)
			localState.order = localState.order[1:]
			window.Invalidate()
		}
	}

	// Buttons
	if ss.exitButton.button.Clicked(gtx) {
		localState.active = !localState.active
		localState.progressBar.progress = 1
		localState.exit <- true
	}
	if ss.infoButton.button.Clicked(gtx) {
		ss.infoButton.active = !ss.infoButton.active
	}
	if ss.timerButton.button.Clicked(gtx) {
		ss.timerButton.active = !ss.timerButton.active
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
	// Opening File
	filename := localState.order[0]
	log.Printf("Opening: %v\n", filename)
	filepath := localState.cfg.Directory + string(os.PathSeparator) + filename
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	//Decoding Image
	var img image.Image
	var imgCfg image.Config
	if filename[len(filename)-4:] == "webp" {
		img, err = webp.Decode(file)
		if err != nil {
			return err
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		imgCfg, err = webp.DecodeConfig(file)
		if err != nil {
			return err
		}
	} else {
		img, _, err = image.Decode(file)
		if err != nil {
			return err
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		imgCfg, _, err = image.DecodeConfig(file)
		if err != nil {
			return err
		}
	}

	// Updating State
	ss.currentImage.Image = img
	ss.currentImage.Filename = filename
	ss.currentImage.Filesize = fileInfo.Size()
	ss.currentImage.Size = image.Point{imgCfg.Width, imgCfg.Height}
	localState.order = localState.order[1:]
	localState.progressBar.progress = 0
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

func slideshowImageButtons(iconButton *iconButton, theme *material.Theme) layout.FlexChild {
	iButton := material.IconButton(theme, iconButton.button, iconButton.icon, iconButton.label)
	iButton.Background = color.NRGBA{0, 0, 0, 0}
	if iconButton.active {
		iButton.Color = color.NRGBA{255, 255, 255, localState.opacity}
	} else {
		iButton.Color = color.NRGBA{255, 255, 255, uint8(localState.opacity / 2)}
	}
	return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		return iButton.Layout(gtx)
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
