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

//export windowOrderedOut
func windowOrderedOut(ptr unsafe.Pointer) {
	mainEventLoop.push(func() {
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
	})
}
