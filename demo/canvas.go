package main

import (
	"github.com/unixpickle/gogui"
	"math"
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
	
	// Fill red rectangle
	c.SetFill(1, 0, 0, 1)
	c.FillRect(gogui.Rect{10, 10, 50, 50})
	
	// Draw a triangle
	c.SetStroke(0, 0, 1, 1)
	c.BeginPath()
	c.MoveTo(170, 200)
	c.LineTo(200, 260)
	c.LineTo(140, 260)
	c.ClosePath()
	c.StrokePath()
	
	// Draw a regular polygon
	vertices := 5
	c.SetFill(0, 1, 0, 1)
	c.SetThickness(10)
	c.BeginPath()
	for i := 0; i < vertices; i++ {
		angle := (math.Pi * 2) * (float64(i+1) / float64(vertices))
		y := math.Sin(angle)*30 + 300
		x := math.Cos(angle)*30 + 300
		if i == 0 {
			c.MoveTo(x, y)
		} else {
			c.LineTo(x, y)
		}
	}
	c.ClosePath()
	c.FillPath()
	
	// Draw a "circle"
	vertices = 100
	c.SetFill(0, 0.8, 0.8, 1)
	c.SetThickness(10)
	c.BeginPath()
	for i := 0; i < vertices; i++ {
		angle := (math.Pi * 2) * (float64(i+1) / float64(vertices))
		y := math.Sin(angle)*40 + 70
		x := math.Cos(angle)*40 + 280
		if i == 0 {
			c.MoveTo(x, y)
		} else {
			c.LineTo(x, y)
		}
	}
	c.ClosePath()
	c.FillPath()
	
	c.Flush()
}
