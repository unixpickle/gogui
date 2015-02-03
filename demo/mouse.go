package main

import (
	"fmt"
	"github.com/unixpickle/gogui"
	"os"
)

func main() {
	go openWindow()
	gogui.Main(&gogui.AppInfo{Name: "Demo"})
}

func openWindow() {
	w, _ := gogui.NewWindow(gogui.Rect{0, 0, 400, 400})
	w.SetTitle("Demo")
	w.Center()
	w.Show()
	w.SetCloseHandler(func() {
		os.Exit(0)
	})
	w.SetMouseDownHandler(func(evt gogui.MouseEvent) {
		fmt.Println("Got mouse down", evt)
	})
}
