package gogui

import (
	"C"
	"sync"
)

var eventLock sync.Mutex
var eventQueue []func()

func pushEvent(f func()) {
	eventLock.Lock()
	defer eventLock.Unlock()
	eventQueue = append(eventQueue, f)
}

//export runNextEvent
func runNextEvent() {
	eventLock.Lock()
	if len(eventQueue) == 0 {
		eventLock.Unlock()
		return
	}
	f := eventQueue[0]
	eventQueue = eventQueue[1:]
	eventLock.Unlock()
	f()
}
