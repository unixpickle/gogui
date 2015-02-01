// +build darwin,cgo

package gogui

import (
	"C"
	"unsafe"
)

//export windowOrderedOut
func windowOrderedOut(ptr unsafe.Pointer) {
	// Don't block the UI thread. This will not result in a race condition
	// because the Show() method won't do anything unless w.showing is false.
	go func() {
		globalLock.Lock()
		defer globalLock.Unlock()
		for i, x := range showingWindows {
			w := x.(*window)
			if w.pointer == ptr {
				if !w.showing {
					// Could happen in a race condition type situation where the
					// user closed the window just as the code ran w.Hide().
					break
				}
				w.showing = false
				showingWindows[i] = showingWindows[len(showingWindows) - 1]
				showingWindows[len(showingWindows) - 1] = nil
				showingWindows = showingWindows[0 : len(showingWindows)-1]
				if w.onClose != nil {
					go w.onClose()
				}
				break
			}
		}
	}()
}
