package main

import (
	"github.com/unixpickle/gogui"
)

func main() {
	// Open the evil window once the loop starts.
	gogui.RunOnMain(doEvil)
	
	// Run the loop.
	gogui.Main(&gogui.AppInfo{Name: "Evil"})
}

func doEvil() {
	// Create teh window
	w, _ := gogui.NewWindow(gogui.Rect{0, 0, 400, 400})
	w.Center()
	w.SetTitle("Evil")
	w.Show()
	
	// If they somehow get to the close button, the window automatically
	// re-opens itself.
	w.SetCloseHandler(func() {
		w.Show()
	})
	
	// When the move their mouse, move the window so that their mouse is in
	// the center of the window.
	w.SetMouseMoveHandler(func(e gogui.MouseEvent) {
		frame := w.Frame()
		frame.X += e.X - 200
		frame.Y += e.Y - 200
		w.SetFrame(frame)
	})
	
	// Create a canvas with the taunt.
	c, _ := gogui.NewCanvas(gogui.Rect{0, 0, 400, 400})
	c.SetDrawHandler(func(c gogui.DrawContext) {
		c.FillText("Try moving your mouse...", 10, 10)
	})
	w.Add(c)
}
