// +build darwin,cgo

package gogui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

#define ASSERT_MAIN NSCAssert([NSThread isMainThread], \
	@"Call must be from main thread.")

extern void windowClosed(void * ptr);

@interface ContentView : NSView {
}

@end

@implementation ContentView

- (BOOL)isFlipped {
	return YES;
}

@end

@interface SimpleWindow : NSWindow {
	NSEventModifierFlags flags;
}

- (id)initWithFrame:(NSRect)rect;

@end

@implementation SimpleWindow

- (id)initWithFrame:(NSRect)r {
	self = [super initWithContentRect:r
		styleMask:(NSTitledWindowMask|NSClosableWindowMask)
		backing:NSBackingStoreBuffered
		defer:NO];
	if (self) {
		ContentView * cv = [[ContentView alloc]
			initWithFrame:NSMakeRect(0, 0, r.size.width, r.size.height)];
		[self setReleasedWhenClosed:NO];
		[self setContentView:cv];
		[cv release];
	}
	return self;
}

- (void)orderOut:(id)sender {
	if (sender) {
		windowClosed((void *)self);
	}
}

@end

void AddToWindow(void * w, void * v) {
	ASSERT_MAIN;
	[[(NSWindow *)w contentView] addSubview:(NSView *)v];
}

void CenterWindow(void * w) {
	ASSERT_MAIN;
	[(NSWindow *)w center];
}

void * CreateWindow(double x, double y, double w, double h) {
	ASSERT_MAIN;
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	return (void *)[[SimpleWindow alloc] initWithFrame:r];
}

void DestroyWindow(void * ptr) {
	ASSERT_MAIN;
	NSWindow * w = (NSWindow *)ptr;
	[w orderOut:nil];
	NSCAssert([w.contentView subviews].count == 0,
		@"Should not destroy window with subviews");
	[w release];
}

void GetWindowFrame(void * ptr, double * x, double * y, double * w,
	double * h) {
	ASSERT_MAIN;
	NSRect r = [(NSWindow *)ptr frame];
	*x = (double)r.origin.x;
	*y = (double)r.origin.y;
	*w = (double)r.size.width;
	*h = (double)r.size.height;
}

void HideWindow(void * ptr) {
	ASSERT_MAIN;
	NSWindow * w = (NSWindow *)ptr;
	[(NSWindow *)w orderOut:nil];
}

void RemoveFromSuperview(void * v) {
	ASSERT_MAIN;
	[(NSView *)v removeFromSuperview];
}

void SetWindowFrame(void * ptr, double x, double y, double w, double h) {
	ASSERT_MAIN;
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	[(NSWindow *)ptr setFrame:r display:YES];
}

void SetWindowTitle(void * w, const char * title) {
	ASSERT_MAIN;
	[(NSWindow *)w setTitle:[NSString stringWithUTF8String:title]];
}

void ShowWindow(void * w) {
	ASSERT_MAIN;
	[(NSWindow *)w makeKeyAndOrderFront:nil];
	[NSApp activateIgnoringOtherApps:YES];
}

*/
import "C"

import (
	"runtime"
	"unsafe"
)

var showingWindows = []Window{}

type window struct {
	windowEvents

	pointer unsafe.Pointer
	widgets []Widget
	showing bool
}

// NewWindow creates a new window with a given content rectangle.
// The created window will not be showing by default.
// You must call this from the main goroutine.
func NewWindow(r Rect) (Window, error) {
	ptr := C.CreateWindow(C.double(r.X), C.double(r.Y), C.double(r.Width),
		C.double(r.Height))
	res := &window{pointer: ptr, widgets: []Widget{}}
	runtime.SetFinalizer(res, finalizeWindow)
	return res, nil
}

// ShowingWindows returns a slice containing all the currently active windows.
// You must call this from the main goroutine.
func ShowingWindows() []Window {
	cpy := make([]Window, len(showingWindows))
	copy(cpy, showingWindows)
	return cpy
}

func (w *window) Add(widget Widget) {
	v, ok := widget.(ptrView)
	if !ok {
		panic("Widget is not ptrView")
	}
	w.widgets = append(w.widgets, widget)
	v.setParent(w)
	ptr := v.viewPointer()
	C.AddToWindow(w.pointer, ptr)
}

func (w *window) Center() {
	C.CenterWindow(w.pointer)
}

func (w *window) Children() []Widget {
	cpy := make([]Widget, len(w.widgets))
	copy(cpy, w.widgets)
	return cpy
}

func (w *window) Focus() {
	w.Show()
}

func (w *window) Frame() Rect {
	var x, y, width, height C.double
	C.GetWindowFrame(w.pointer, &x, &y, &width, &height)
	return Rect{float64(x), float64(y), float64(width), float64(height)}
}

func (w *window) Hide() {
	if !w.showing {
		return
	}
	w.showing = false
	C.HideWindow(w.pointer)
	for i, x := range showingWindows {
		if x.(*window) == w {
			// Remove the window from the list. Note also how we set the last
			// element of the list to nil in order to allow it to be garbage
			// collected sooner (since it's not held by the slice).
			showingWindows[i] = showingWindows[len(showingWindows)-1]
			showingWindows[len(showingWindows)-1] = nil
			showingWindows = showingWindows[0 : len(showingWindows)-1]
			break
		}
	}
}

func (w *window) Parent() Widget {
	return nil
}

func (w *window) Remove() {
}

func (w *window) SetFrame(r Rect) {
	C.SetWindowFrame(w.pointer, C.double(r.X), C.double(r.Y),
		C.double(r.Width), C.double(r.Height))
}

func (w *window) SetTitle(title string) {
	C.SetWindowTitle(w.pointer, C.CString(title))
}

func (w *window) Show() {
	if w.showing {
		return
	}
	w.showing = true
	C.ShowWindow(w.pointer)
	showingWindows = append(showingWindows, w)
}

func (w *window) Showing() bool {
	return w.showing
}

func (w *window) removeView(v ptrView) {
	ptr := v.viewPointer()
	C.RemoveFromSuperview(ptr)

	// Find the first widget with the same underlying pointer.
	for i, x := range w.widgets {
		aPtr := x.(ptrView).viewPointer()
		if aPtr == ptr {
			// We set the last item to nil to give garbage collection a little
			// nudge in the right direction.
			w.widgets[i] = w.widgets[len(w.widgets)-1]
			w.widgets[len(w.widgets)-1] = nil
			w.widgets = w.widgets[0 : len(w.widgets)-1]
			break
		}
	}
}

func finalizeWindow(w *window) {
	RunOnMain(func() {
		for len(w.widgets) > 0 {
			w.widgets[0].Remove()
		}
		C.DestroyWindow(w.pointer)
	})
}
