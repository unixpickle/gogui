package main

import (
	"github.com/unixpickle/gogui"
	"os"
)

func main() {
	// Setup the window once the loop starts.
	gogui.RunOnMain(func() {
		w, _ := gogui.NewWindow(gogui.Rect{0, 0, 400, 400})
		w.SetTitle("Demo")
		w.Center()
		w.Show()
		w.SetCloseHandler(func() {
			os.Exit(0)
		})
	})
	
	// Run the loop.
	gogui.Main(&gogui.AppInfo{Name: "Demo"})
}

func openWindow() {
}
