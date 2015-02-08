package main

import (
	"github.com/unixpickle/gogui"
	"math"
	"os"
	"strconv"
	"time"
)

const ClockSize = 400

func drawClock(c gogui.DrawContext) {
	c.SetFill(gogui.Color{0, 0, 0, 1})
	c.FillEllipse(gogui.Rect{0, 0, ClockSize, ClockSize})
	
	// Draw the 12 numbers around the clock
	c.SetFill(gogui.Color{1, 1, 1, 1})
	for i := 0; i < 12; i++ {
		numberStr := strconv.Itoa(((i + 2) % 12) + 1)
		angle := math.Pi * 2.0 * float64(i) / float64(12)
		x := math.Cos(angle) * (ClockSize - 20) / 2
		y := math.Sin(angle) * (ClockSize - 20) / 2
		x += ClockSize / 2
		y += ClockSize / 2
		c.FillText(numberStr, x-5, y-11)
	}
	
	// Get the times for the hands
	t := time.Now()
	hour, min, sec := t.Clock()
	
	// Draw the hour hand
	hour = hour % 12
	c.SetThickness(5)
	c.SetStroke(gogui.Color{1, 0, 0, 1})
	drawHand(c, float64(hour)/12.0, ClockSize/4.5)
	
	// Draw the minute hand
	c.SetThickness(3)
	c.SetStroke(gogui.Color{0, 0, 1, 1})
	drawHand(c, float64(min)/60.0, ClockSize/3.5)
	
	// Draw second hand
	c.SetThickness(2)
	c.SetStroke(gogui.Color{1, 1, 1, 1})
	drawHand(c, float64(sec)/60.0, ClockSize/3.0)
}

func drawHand(c gogui.DrawContext, fraction float64, length float64) {
	angle := math.Pi*2.0*fraction - math.Pi/2.0
	x := ClockSize/2 + math.Cos(angle)*length
	y := ClockSize/2 + math.Sin(angle)*length
	c.BeginPath()
	c.MoveTo(ClockSize/2, ClockSize/2)
	c.LineTo(x, y)
	c.StrokePath()
}

func main() {
	gogui.RunOnMain(setupClock)
	gogui.Main(&gogui.AppInfo{Name: "Clock"})
}

func setupClock() {
	w, _ := gogui.NewWindow(gogui.Rect{0, 0, ClockSize, ClockSize})
	w.SetTitle("Clock")
	w.Center()
	w.Show()
	
	w.SetCloseHandler(func() {
		os.Exit(1)
	})
	
	c, _ := gogui.NewCanvas(gogui.Rect{0, 0, ClockSize, ClockSize})
	w.Add(c)
	
	c.SetDrawHandler(drawClock)
	
	// Redraw periodically
	go func() {
		for {
			gogui.RunOnMain(func() {
				c.NeedsUpdate()
			})
			time.Sleep(time.Second)
		}
	}()
}
