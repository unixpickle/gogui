// +build darwin,cgo

package gogui

import "fmt"
import (
	"C"
	"unsafe"
)

//export canvasDrawRect
func canvasDrawRect(windowPtr, canvas, ctx unsafe.Pointer) {
	for _, w := range showingWindows {
		wptr := w.(*window)
		if wptr.pointer == windowPtr {
			// Found the window
			for _, child := range w.Children() {
				if child.(ptrView).viewPointer() == canvas {
					// Found the canvas; call the draw handler.
					canvas := child.(Canvas)
					if h := canvas.DrawHandler(); h != nil {
						c := &drawContext{ctx}
						h(c)
						c.pointer = nil
					}
				}
			}
		}
	}
}
