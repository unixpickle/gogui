// +build darwin,cgo

package gogui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void RunMain(void (^ block)(void));
void GetViewFrame(void * ptr, double * x, double * y, double * w, double * h);
void RemoveView(void * ptr);
void SetViewFrame(void * ptr, double x, double y, double w, double h);

@interface Canvas : NSView {
}
@end

@implementation Canvas

- (void)drawRect:(NSRect)dirtyRect {
	// TODO: this
}

@end

void * CreateCanvas(double x, double y, double w, double h) {
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	__block Canvas * res = nil;
	RunMain(^{
		res = [[Canvas alloc] initWithFrame:r];
	});
	return (void *)res;
}

void DestroyCanvas(void * ptr) {
	Canvas * c = (Canvas *)ptr;
	RunMain(^{
		[c removeFromSuperview];
		[c release];
	});
}

void RemoveCanvas(void * ptr) {
	Canvas * c = (Canvas *)ptr;
	RunMain(^{
		[c removeFromSuperview];
	});
}
*/
import "C"

import "unsafe"

type canvas struct {
	parent  Widget
	pointer unsafe.Pointer 
}

// NewCanvas creates a new canvas with the given frame.
func NewCanvas(r Rect) (Canvas, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	ptr := C.CreateCanvas(C.double(r.X), C.double(r.Y), C.double(r.Width),
		C.double(r.Height))
	res := &canvas{nil, ptr}
	return res, nil
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
	C.GetViewFrame(c.pointer, &x, &y, &w, &h)
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
	C.RemoveView(c.pointer)
}

func (c *canvas) SetFrame(r Rect) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if c.pointer == nil {
		panic("Canvas is invalid.")
	}
	C.SetViewFrame(c.pointer, C.double(r.Y), C.double(r.Y), C.double(r.Width),
		C.double(r.Height))
}