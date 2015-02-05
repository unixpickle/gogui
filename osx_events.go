// +build darwin,cgo

package gogui

import (
	"C"
	"sync"
	"unsafe"
)

type eventLoop struct {
	lock    sync.Mutex
	waiting []func()
	trigger chan struct{}
}

var mainEventLoop = eventLoop{sync.Mutex{}, []func(){}, make(chan struct{})}

func (e *eventLoop) main() {
	for {
		<-e.trigger
		e.lock.Lock()
		waiting := e.waiting
		e.waiting = []func(){}
		e.lock.Unlock()
		for _, evt := range waiting {
			evt()
		}
	}
}

func (e *eventLoop) push(evt func()) {
	e.lock.Lock()
	e.waiting = append(e.waiting, evt)
	e.lock.Unlock()
	select {
	case e.trigger <- struct{}{}:
	default:
	}
}

func findWindow(ptr unsafe.Pointer) *window {
	globalLock.Lock()
	defer globalLock.Unlock()
	for _, x := range showingWindows {
		w := x.(*window)
		if w.pointer == ptr {
			return w
		}
	}
	return nil
}

func makeKeyEvent(charCode, keyCode, flags C.int) KeyEvent {
	res := KeyEvent{CharCode: int(charCode), KeyCode: int(keyCode)}
	f := int(flags)
	if (f & 1) != 0 {
		res.AltKey = true
	}
	if (f & 2) != 0 {
		res.CtrlKey = true
	}
	if (f & 4) != 0 {
		res.MetaKey = true
	}
	if (f & 8) != 0 {
		res.ShiftKey = true
	}
	return res
}

//export windowKeyDown
func windowKeyDown(ptr unsafe.Pointer, charCode, keyCode, flags C.int) {
	mainEventLoop.push(func() {
		window := findWindow(ptr)
		if window == nil {
			return
		}
		if handler := window.KeyDownHandler(); handler != nil {
			handler(makeKeyEvent(charCode, keyCode, flags))
		}
	})
}

//export windowKeyPress
func windowKeyPress(ptr unsafe.Pointer, charCode, keyCode, flags C.int) {
	mainEventLoop.push(func() {
		window := findWindow(ptr)
		if window == nil {
			return
		}
		if handler := window.KeyPressHandler(); handler != nil {
			handler(makeKeyEvent(charCode, keyCode, flags))
		}
	})
}

//export windowKeyUp
func windowKeyUp(ptr unsafe.Pointer, charCode, keyCode, flags C.int) {
	mainEventLoop.push(func() {
		window := findWindow(ptr)
		if window == nil {
			return
		}
		if handler := window.KeyUpHandler(); handler != nil {
			handler(makeKeyEvent(charCode, keyCode, flags))
		}
	})
}

//export windowMouseDown
func windowMouseDown(ptr unsafe.Pointer, x, y C.double) {
	mainEventLoop.push(func() {
		window := findWindow(ptr)
		if window == nil {
			return
		}
		if handler := window.MouseDownHandler(); handler != nil {
			evt := MouseEvent{float64(x), float64(y)}
			handler(evt)
		}
	})
}

//export windowMouseDrag
func windowMouseDrag(ptr unsafe.Pointer, x, y C.double) {
	mainEventLoop.push(func() {
		window := findWindow(ptr)
		if window == nil {
			return
		}
		if handler := window.MouseDragHandler(); handler != nil {
			evt := MouseEvent{float64(x), float64(y)}
			handler(evt)
		}
	})
}

//export windowMouseMove
func windowMouseMove(ptr unsafe.Pointer, x, y C.double) {
	mainEventLoop.push(func() {
		window := findWindow(ptr)
		if window == nil {
			return
		}
		if handler := window.MouseMoveHandler(); handler != nil {
			evt := MouseEvent{float64(x), float64(y)}
			handler(evt)
		}
	})
}

//export windowMouseUp
func windowMouseUp(ptr unsafe.Pointer, x, y C.double) {
	mainEventLoop.push(func() {
		window := findWindow(ptr)
		if window == nil {
			return
		}
		if handler := window.MouseUpHandler(); handler != nil {
			evt := MouseEvent{float64(x), float64(y)}
			handler(evt)
		}
	})
}

//export windowOrderedOut
func windowOrderedOut(ptr unsafe.Pointer) {
	mainEventLoop.push(func() {
		globalLock.Lock()
		for i, x := range showingWindows {
			w := x.(*window)
			if w.pointer == ptr {
				if !w.showing {
					// Could happen in a race condition type situation where the
					// user closed the window just as the code ran w.Hide().
					break
				}
				
				// Mark the window as hidden.
				w.showing = false
				showingWindows[i] = showingWindows[len(showingWindows) - 1]
				showingWindows[len(showingWindows) - 1] = nil
				showingWindows = showingWindows[0 : len(showingWindows)-1]
				
				// Unlock, call the handler, and return.
				globalLock.Unlock()
				if w.onClose != nil {
					w.onClose()
				}
				return
			}
		}
		globalLock.Unlock()
	})
}
