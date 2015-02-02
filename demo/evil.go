package main

import (
	"github.com/unixpickle/gogui"
)

func main() {
	for i := 0; i < 10; i++ {
		go openWindow(float64(i) * 100)
	}
	gogui.Main(&gogui.AppInfo{Name: "Demo"})
}

func openWindow(coord float64) {
	w, _ := gogui.NewWindow(gogui.Rect{coord, coord, 400, 400})
	w.SetTitle("Demo")
	w.Show()
	w.SetCloseHandler(func() {
		w.Show()
	})
}
