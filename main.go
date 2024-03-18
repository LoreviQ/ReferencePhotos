package main

import (
	"gioui.org/app"
)

func main() {
	go createGUI(800, 1250)
	app.Main()
}
