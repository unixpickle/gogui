// +build darwin,cgo

package gogui

import (
	"C"
	"unsafe"
)

const (
	keyEventDown = iota
	keyEventUp   = iota
)

const (
	keyFlagAlt   = 1
	keyFlagCtrl  = 2
	keyFlagMeta  = 4
	keyFlagShift = 8
)

const (
	mouseEventDown = iota
	mouseEventDrag = iota
	mouseEventMove = iota
	mouseEventUp   = iota
)

func findWindow(ptr unsafe.Pointer) *window {
	for _, w := range showingWindows {
		wptr := w.(*window)
		if wptr.pointer == ptr {
			return wptr
		}
	}
	return nil
}

//export windowClosed
func windowClosed(ptr unsafe.Pointer) {
	for i, w := range showingWindows {
		wptr := w.(*window)
		if wptr.pointer == ptr {
			// Remove the window from the list and set it not showing.
			wptr.showing = false
			showingWindows[i] = showingWindows[len(showingWindows)-1]
			showingWindows[len(showingWindows)-1] = nil
			showingWindows = showingWindows[0 : len(showingWindows)-1]

			// Call the close handler
			if h := w.CloseHandler(); h != nil {
				h()
			}
		}
	}
}

//export windowKeyFlagsChanged
func windowKeyFlagsChanged(ptr unsafe.Pointer, flags int) {
	// TODO: this
}

//export windowKeyEvent
func windowKeyEvent(ptr unsafe.Pointer, eventType int, chars, modChars *C.char,
	keyCode, mods int) {
	// TODO: this
}

//export windowMouseEvent
func windowMouseEvent(ptr unsafe.Pointer, eventType int, x, y C.double) {
	w := findWindow(ptr)
	if w == nil {
		return
	}
	
	// Get the handler.
	var handler MouseHandler
	switch eventType {
	case mouseEventDown:
		handler = w.MouseDownHandler()
	case mouseEventDrag:
		handler = w.MouseDragHandler()
	case mouseEventMove:
		handler = w.MouseMoveHandler()
	case mouseEventUp:
		handler = w.MouseUpHandler()
	default:
		panic("Unknown mouse event.")
	}
	
	// Call the handler if there is one.
	if handler == nil {
		return
	}
	handler(MouseEvent{float64(x), float64(y)})
}
