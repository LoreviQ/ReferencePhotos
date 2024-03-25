package main

import (
	"log"
	"time"

	"gioui.org/app"
	"github.com/gen2brain/beeep"
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
		"10m": 0.002,
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
	var err error
	if localState.progressBar.progress > 1-p*300 && localState.progressBar.sounds == 0 {
		log.Println("Audio Alert 1")
		err = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		localState.progressBar.sounds++
	}
	if localState.progressBar.progress > 1-p*200 && localState.progressBar.sounds == 1 {
		log.Println("Audio Alert 2")
		err = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		localState.progressBar.sounds++
	}
	if localState.progressBar.progress > 1-p*100 && localState.progressBar.sounds == 2 {
		log.Println("Audio Alert 3")
		err = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		localState.progressBar.sounds++
	}
	if err != nil {
		log.Println(err)
	}
}
