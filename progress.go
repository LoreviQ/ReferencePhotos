package main

import (
	"time"

	"gioui.org/app"
)

func getProgressIncrementer(increment string) chan float32 {
	incrementMap := map[string]float32{
		"30s": 0.00033333333,
		"45s": 0.00022222222,
		"1m":  0.00016666666,
		"2m":  0.00008333333,
		"5m":  0.00003333333,
		"10m": 0.01,
		//"10m": 0.00001666666,
	}
	progressIncrementer := make(chan float32)
	go func() {
		for {
			time.Sleep(time.Second / 100)
			progressIncrementer <- incrementMap[increment]
		}
	}()
	return progressIncrementer
}

func incrementProgress(window *app.Window, progressIncrementer chan float32, exit chan bool) {
	for {
		select {
		case p := <-progressIncrementer:
			if localState.active && localState.progress < 1 && !localState.paused {
				localState.progress += p
				window.Invalidate()
			}
		case <-exit:
			return
		}
	}
}
