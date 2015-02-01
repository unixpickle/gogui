// +build darwin,cgo

package gogui

/*
void AddToWindow(void * wind, void * view);
void CanvasNeedsDisplay(void * ptr);
void CenterWindow(void * ptr);
void * CreateCanvas(unsigned long long ident, double x, double y, double w,
	double h);
void * CreateWindow(double x, double y, double w, double h);
void DestroyCanvas(void * ptr);
void DestroyWindow(void * ptr);
void GetCanvasFrame(void * ptr, double * x, double * y, double * w,
	double * h);
void GetWindowFrame(void * ptr, double * x, double * y, double * w,
	double * h);
void HideWindow(void * ptr);
void MainLoop(const char * name);
void RemoveCanvas(void * ptr);
void SetCanvasFrame(void * ptr, double x, double y, double w, double h);
void SetWindowFrame(void * ptr, double x, double y, double w, double h);
void SetWindowTitle(void * ptr, const char * title);
void ShowWindow(void * ptr);
*/
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

var allCanvases = map[uint64]*canvas{}
var canvasIdCounter uint64
var globalLock sync.Mutex

func init() {
	// Make sure main.main runs on the main OS thread.
	runtime.LockOSThread()
}

type canvas struct {
	id      uint64
	parent  Widget
	pointer unsafe.Pointer 
}

// NewCanvas creates a new canvas with the given frame.
func NewCanvas(r Rect) (Canvas, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	id := canvasIdCounter
	canvasIdCounter++
	ptr := C.CreateCanvas(C.ulonglong(id), C.double(r.X), C.double(r.Y),
		C.double(r.Width), C.double(r.Height))
	res := &canvas{id, nil, ptr}
	allCanvases[id] = res
	return res, nil
}

//export canvasDraw
func canvasDraw(id C.ulonglong) {
	
}

func (c *canvas) Begin() DrawContext {
	// TODO: this
	return nil
}

func (c *canvas) Destroy() {
	globalLock.Lock()
	defer globalLock.Unlock()
	if c.pointer == nil {
		panic("Canvas is invalid.")
	}
	C.DestroyCanvas(c.pointer)
	c.pointer = nil
	delete(allCanvases, c.id)
}

func (c *canvas) Flush() {
	// TODO: this
}

func (c *canvas) Frame() Rect {
	globalLock.Lock()
	defer globalLock.Unlock()
	if c.pointer == nil {
		panic("Canvas is invalid.")
	}
	var x, y, w, h C.double
	C.GetCanvasFrame(c.pointer, &x, &y, &w, &h)
	return Rect{float64(x), float64(y), float64(w), float64(h)}
}

func (c *canvas) Parent() Widget {
	globalLock.Lock()
	defer globalLock.Unlock()
	if c.pointer == nil {
		panic("Canvas is invalid.")
	}
	return c.parent
}

func (c *canvas) Remove() {
	globalLock.Lock()
	defer globalLock.Unlock()
	if c.pointer == nil {
		panic("Canvas is invalid.")
	}
	if c.parent == nil {
		return
	}
	
	// Remove references to this widget from its parent.
	if w, ok := c.parent.(*window); ok {
		for i, x := range w.widgets {
			if widget, ok := x.(*canvas); ok && widget == c {
				w.widgets[i] = w.widgets[len(w.widgets) - 1]
				w.widgets = w.widgets[0 : len(w.widgets)-1]
			}
		}
	} else {
		panic("Unknown parent type.")
	}
	c.parent = nil
	
	// Remove the actual view
	C.RemoveCanvas(c.pointer)
}

func (c *canvas) SetFrame(r Rect) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if c.pointer == nil {
		panic("Canvas is invalid.")
	}
	C.SetCanvasFrame(c.pointer, C.double(r.Y), C.double(r.Y), C.double(r.Width),
		C.double(r.Height))
}

type window struct {
	pointer unsafe.Pointer
	widgets []Widget
}

func NewWindow(r Rect) (Window, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	ptr := C.CreateWindow(C.double(r.Y), C.double(r.Y), C.double(r.Width),
		C.double(r.Height))
	return &window{ptr, []Widget{}}, nil
}

func (w *window) Add(widget Widget) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if canvas, ok := widget.(*canvas); ok {
		canvas.parent = w
		C.AddToWindow(w.pointer, canvas.pointer)
	} else {
		panic("Unknown widget type.")
	}
	w.widgets = append(w.widgets, widget)
}

func (w *window) Center() {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	C.CenterWindow(w.pointer)
}

func (w *window) Children() []Widget {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	cpy := make([]Widget, len(w.widgets))
	copy(cpy, w.widgets)
	return cpy
}

func (w *window) Destroy() {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	
	// Remove all children.
	for len(w.widgets) > 0 {
		w.widgets[0].Remove()
	}
	
	// Destroy the window.
	C.DestroyWindow(w.pointer)
	
	// Invalidate the object.
	w.pointer = nil
}

func (w *window) Focus() {
	w.Show()
}

func (window *window) Frame() Rect {
	globalLock.Lock()
	defer globalLock.Unlock()
	if window.pointer == nil {
		panic("Window is invalid.")
	}
	var x, y, w, h C.double
	C.GetWindowFrame(window.pointer, &x, &y, &w, &h)
	return Rect{float64(x), float64(y), float64(w), float64(h)}
}

func (w *window) Hide() {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	C.HideWindow(w.pointer)
}

func (w *window) Parent() Widget {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	return nil
}

func (w *window) Remove() {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
}

func (w *window) SetFrame(r Rect) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	C.SetWindowFrame(w.pointer, C.double(r.X), C.double(r.Y),
		C.double(r.Width), C.double(r.Height))
}

func (w *window) SetTitle(title string) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	C.SetWindowTitle(w.pointer, C.CString(title))
}

func (w *window) Show() {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	C.ShowWindow(w.pointer)
}

// Main runs the Cocoa run-loop. You must call this from main.main.
func Main(info *AppInfo) {
	C.MainLoop(C.CString(info.Name))
}
