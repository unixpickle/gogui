// +build darwin,cgo

package gogui

import (
	"C"
	"strings"
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

func callKeyDown(w *window, evt KeyEvent) {
	if h := w.KeyDownHandler(); h != nil {
		h(evt)
	}
}

func callKeyPress(w *window, evt KeyEvent) {
	if h := w.KeyPressHandler(); h != nil {
		h(evt)
	}
}

func callKeyUp(w *window, evt KeyEvent) {
	if h := w.KeyUpHandler(); h != nil {
		h(evt)
	}
}

func findWindow(ptr unsafe.Pointer) *window {
	for _, w := range showingWindows {
		wptr := w.(*window)
		if wptr.pointer == ptr {
			return wptr
		}
	}
	return nil
}

func makeKeyEvent(keyCode, charCode int, flags int) KeyEvent {
	res := KeyEvent{KeyCode: keyCode, CharCode: charCode}
	if (flags & keyFlagAlt) != 0 {
		res.AltKey = true
	}
	if (flags & keyFlagCtrl) != 0 {
		res.CtrlKey = true
	}
	if (flags & keyFlagMeta) != 0 {
		res.MetaKey = true
	}
	if (flags & keyFlagShift) != 0 {
		res.ShiftKey = true
	}
	return res
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
	w := findWindow(ptr)
	if w == nil {
		return
	}
	difference := flags ^ w.modifiers
	w.modifiers = flags
	if (difference & keyFlagAlt) != 0 {
		evt := makeKeyEvent(-1, 18, flags)
		if (flags & keyFlagAlt) != 0 {
			callKeyDown(w, evt)
		} else {
			callKeyUp(w, evt)
		}
	}
	if (difference & keyFlagCtrl) != 0 {
		evt := makeKeyEvent(-1, 17, flags)
		if (flags & keyFlagCtrl) != 0 {
			callKeyDown(w, evt)
		} else {
			callKeyUp(w, evt)
		}
	}
	if (difference & keyFlagMeta) != 0 {
		evt := makeKeyEvent(-1, 91, flags)
		if (flags & keyFlagMeta) != 0 {
			callKeyDown(w, evt)
		} else {
			callKeyUp(w, evt)
		}
	}
	if (difference & keyFlagShift) != 0 {
		evt := makeKeyEvent(-1, 16, flags)
		if (flags & keyFlagShift) != 0 {
			callKeyDown(w, evt)
		} else {
			callKeyUp(w, evt)
		}
	}
}

//export windowKeyEvent
func windowKeyEvent(ptr unsafe.Pointer, eventType int, chars, modChars *C.char,
	keyCode, flags int) {
	w := findWindow(ptr)
	if w == nil {
		return
	}
	
	// Normally, the character code is just the ASCII character of the key.
	rawCodeStr := C.GoString(chars)
	rawCodeStr = strings.ToUpper(rawCodeStr)
	rawCode := int([]rune(rawCodeStr)[0])
	modCode := int([]rune(C.GoString(modChars))[0])
	
	// Certain keys have weird character codes that Cocoa doesn't report nicely.
	mapping := map[int]int{51: 8, 123: 37, 126: 38, 124: 39, 125: 40, 117: 46}
	if mapped, ok := mapping[keyCode]; ok {
		rawCode = mapped
		modCode = mapped
	}
	
	if eventType == keyEventDown {
		// Call down and press events.
		evt := makeKeyEvent(keyCode, rawCode, flags)
		callKeyDown(w, evt)
		evt.CharCode = modCode
		callKeyPress(w, evt)
	} else {
		// Call key up event.
		evt := makeKeyEvent(keyCode, rawCode, flags)
		callKeyUp(w, evt)
	}
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
