package main

// TODO: I gotta rewrite this entire thing...

import (
	"github.com/unixpickle/gogui"
	"os"
)

func main() {
	gogui.RunOnMain(openWindow)
	gogui.Main(&gogui.AppInfo{Name: "Demo"})
}

func openWindow() {
	w, _ := gogui.NewWindow(gogui.Rect{0, 0, 400, 400})
	c, _ := gogui.NewCanvas(gogui.Rect{0, 0, 400, 400})
	w.Add(c)
	w.SetTitle("Demo")
	w.Center()
	w.Show()
	w.SetCloseHandler(func() {
		os.Exit(0)
	})
	
	path := []gogui.MouseEvent{}
	w.SetMouseDownHandler(func(evt gogui.MouseEvent) {
		path = []gogui.MouseEvent{evt}
		c.NeedsUpdate()
	})
	w.SetMouseDragHandler(func(evt gogui.MouseEvent) {
		path = append(path, evt)
		c.NeedsUpdate()
	})
	w.SetMouseMoveHandler(func(evt gogui.MouseEvent) {
		path = []gogui.MouseEvent{evt}
		c.NeedsUpdate()
	})
	w.SetMouseUpHandler(func(evt gogui.MouseEvent) {
		path = []gogui.MouseEvent{evt}
		c.NeedsUpdate()
	})
	
	c.SetDrawHandler(func(ctx gogui.DrawContext) {
		drawLines(ctx, path)
	})
}

func drawCircle(c gogui.DrawContext, evt gogui.MouseEvent) {
	c.SetFill(0, 0, 0, 1)
	c.FillEllipse(gogui.Rect{evt.X - 5, evt.Y - 5, 10, 10})
}

func drawLines(c gogui.DrawContext, evts []gogui.MouseEvent) {
	if len(evts) == 0 {
		return
	} else if len(evts) == 1 {
		drawCircle(c, evts[0])
		return
	}
	
	c.SetStroke(1, 0, 0, 1)
	c.SetThickness(8)
	c.BeginPath()
	c.MoveTo(evts[0].X, evts[0].Y)
	for i := 1; i < len(evts); i++ {
		c.LineTo(evts[i].X, evts[i].Y)
	}
	c.StrokePath()
	
	// Fill the last point
	evt := evts[len(evts) - 1]
	c.SetFill(0.82, 0.29, 0.29, 1)
	c.FillEllipse(gogui.Rect{evt.X - 5, evt.Y - 5, 10, 10})
}
