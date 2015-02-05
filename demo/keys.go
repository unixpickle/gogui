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
	w.SetKeyDownHandler(func(k gogui.KeyEvent) {
		fmt.Println("KeyDown:", k)
	})
	w.SetKeyPressHandler(func(k gogui.KeyEvent) {
		fmt.Println("KeyPress:", k)
	})
	w.SetKeyUpHandler(func(k gogui.KeyEvent) {
		fmt.Println("KeyUp:", k)
	})
}