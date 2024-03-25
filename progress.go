package main

import (
	"fmt"
	"time"

	"gioui.org/app"
)

type progressBar struct {
	progress float32
	time     string
	paused   bool
	sounds   uint8
}

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
			if localState.active && localState.progressBar.progress < 1 && !localState.progressBar.paused {
				localState.progressBar.progress += p
				playSound(p)
				window.Invalidate()
			}
		case <-exit:
			return
		}
	}
}

func playSound(p float32) {
	if localState.progressBar.progress > 1-p*300 && localState.progressBar.sounds == 0 {
		fmt.Println("one")
		localState.progressBar.sounds++
	}
	if localState.progressBar.progress > 1-p*200 && localState.progressBar.sounds == 1 {
		fmt.Println("two")
		localState.progressBar.sounds++
	}
	if localState.progressBar.progress > 1-p*100 && localState.progressBar.sounds == 2 {
		fmt.Println("three")
		localState.progressBar.sounds++
	}
}
