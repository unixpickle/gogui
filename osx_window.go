// +build darwin,cgo

package gogui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void RunMain(void (^ block)(void));

void AddToWindow(void * wind, void * view) {
	NSWindow * w = (NSWindow *)wind;
	NSView * v = (NSView *)view;
	RunMain(^{
		[w.contentView addSubview:v];
	});
}

void CenterWindow(void * ptr) {
	NSWindow * w = (NSWindow *)ptr;
	RunMain(^{
		[w center];
	});
}

void * CreateWindow(double x, double y, double w, double h) {
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	__block NSWindow * res = nil;
	RunMain(^{
		res = [[NSWindow alloc] initWithContentRect:r
			styleMask:(NSTitledWindowMask|NSClosableWindowMask)
			backing:NSBackingStoreBuffered
			defer:NO];
		[res setReleasedWhenClosed:NO];
	});
	return (void *)res;
}

void DestroyWindow(void * ptr) {
	NSWindow * w = (NSWindow *)ptr;
	RunMain(^{
		[w orderOut:nil];
		NSCAssert([w.contentView subviews].count == 0,
			@"Should not destroy window with subviews");
		[w release];
	});
}

void GetWindowFrame(void * ptr, double * x, double * y, double * w,
	double * h) {
	// TODO: use the content frame, not the window frame.
	NSWindow * window = (NSWindow *)ptr;
	RunMain(^{
		NSRect r = [window frame];
		*x = (double)r.origin.x;
		*y = (double)r.origin.y;
		*w = (double)r.size.width;
		*h = (double)r.size.height;
	});
}

void HideWindow(void * ptr) {
	NSWindow * w = (NSWindow *)ptr;
	RunMain(^{
		[w orderOut:nil];
	});
}

void SetWindowFrame(void * ptr, double x, double y, double w, double h) {
	// TODO: use the content frame, not the window frame.
	NSWindow * window = (NSWindow *)ptr;
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	RunMain(^{
		[window setFrame:r display:YES];
	});
}

void SetWindowTitle(void * ptr, const char * title) {
	NSWindow * w = (NSWindow *)ptr;
	RunMain(^{
		[w setTitle:[NSString stringWithUTF8String:title]];
	});
}

void ShowWindow(void * ptr) {
	NSWindow * w = (NSWindow *)ptr;
	RunMain(^{
		[w makeKeyAndOrderFront:nil];
		[NSApp activateIgnoringOtherApps:YES];
	});
}
*/
import "C"

import "unsafe"

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
