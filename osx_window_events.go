// +build darwin,cgo

package gogui

import (
	"C"
	"unsafe"
)

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
