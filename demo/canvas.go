package main

import (
	"github.com/unixpickle/gogui"
	"math"
	"os"
)

func main() {
	gogui.RunOnMain(openWindow)
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
	
	// Set the draw routine.
	c.SetDrawHandler(func(ctx gogui.DrawContext) {
		// Fill red rectangle
		ctx.SetFill(1, 0, 0, 1)
		ctx.FillRect(gogui.Rect{10, 10, 50, 50})
	
		// Draw a triangle
		ctx.SetStroke(0, 0, 0, 1)
		ctx.SetThickness(3)
		ctx.BeginPath()
		ctx.MoveTo(170, 200)
		ctx.LineTo(200, 260)
		ctx.LineTo(140, 260)
		ctx.ClosePath()
		ctx.StrokePath()
	
		// Draw a regular polygon
		vertices := 5
		ctx.SetFill(0, 0, 1, 1)
		ctx.SetThickness(10)
		ctx.BeginPath()
		for i := 0; i < vertices; i++ {
			angle := (math.Pi * 2) * (float64(i+1) / float64(vertices))
			y := math.Sin(angle)*30 + 300
			x := math.Cos(angle)*30 + 300
			if i == 0 {
				ctx.MoveTo(x, y)
			} else {
				ctx.LineTo(x, y)
			}
		}
		ctx.ClosePath()
		ctx.FillPath()
	
		// Draw a circle
		vertices = 100
		ctx.SetFill(0, 0.8, 0.8, 1)
		ctx.FillEllipse(gogui.Rect{240, 30, 80, 80})
	})
	c.NeedsUpdate()
}
