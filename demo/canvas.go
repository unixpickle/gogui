package main

import (
	"github.com/unixpickle/gogui"
	"os"
)

func main() {
	go openWindow()
	gogui.Main(&gogui.AppInfo{Name: "Demo"})
}

func openWindow() {
	// Create the window.
	w, _ := gogui.NewWindow(gogui.Rect{0, 0, 400, 400})
	w.SetTitle("Demo")
	w.Center()
	w.Show()
	w.SetCloseHandler(func() {
		os.Exit(0)
	})
	
	// Create the canvas.
	c, _ := gogui.NewCanvas(gogui.Rect{0, 0, 400, 400})
	w.Add(c)
	c.FillRect(gogui.Rect{0, 0, 50, 50})
	c.Flush()
}
