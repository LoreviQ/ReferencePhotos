package main

import (
	"time"

	"gioui.org/app"
)

var progress float32

func getProgressIncrementer(increment float32) chan float32 {
	progressIncrementer := make(chan float32)
	go func() {
		for {
			time.Sleep(time.Second / 100)
			progressIncrementer <- increment
		}
	}()
	return progressIncrementer
}

func incrementProgress(window *app.Window, progressIncrementer chan float32) {
	for p := range progressIncrementer {
		if localState.active && progress < 1 {
			progress += p
			window.Invalidate()
		}
	}
}
