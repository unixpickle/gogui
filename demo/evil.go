package main

import (
	"github.com/unixpickle/gogui"
)

func main() {
	openWindow()
	gogui.Main(&gogui.AppInfo{Name: "Demo"})
}

func openWindow() {
	w, _ := gogui.NewWindow(gogui.Rect{0, 0, 400, 400})
	w.Center()
	w.SetTitle("Demo")
	w.Show()
	w.SetCloseHandler(func() {
		w.Show()
	})
}
