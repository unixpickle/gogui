package main

import (
	"github.com/unixpickle/gogui"
)

func main() {
	
	// Open the evil window once the loop starts.
	gogui.RunOnMain(func() {
		w, _ := gogui.NewWindow(gogui.Rect{0, 0, 400, 400})
		w.Center()
		w.SetTitle("Evil")
		w.Show()
		w.SetCloseHandler(func() {
			w.Show()
		})
	})
	
	// Run the loop.
	gogui.Main(&gogui.AppInfo{Name: "Evil"})
}
