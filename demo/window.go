package main

import (
	"github.com/unixpickle/gogui"
)

func main() {
	go openWindow()
	gogui.Main(&gogui.AppInfo{Name: "Demo"})
}

func openWindow() {
	w, _ := gogui.NewWindow(gogui.Rect{100, 100, 400, 400})
	w.SetTitle("Demo")
	w.Center()
	w.Show()
}
