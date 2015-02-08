// +build darwin,cgo

package gogui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

#define ASSERT_MAIN NSCAssert([NSThread isMainThread], \
	@"Call must be from main thread.")

enum {
	keyEventDown = 0,
	keyEventUp
};

enum {
	keyFlagAlt   = 1,
	keyFlagCtrl  = 2,
	keyFlagMeta  = 4,
	keyFlagShift = 8
};

enum {
	mouseEventDown = 0,
	mouseEventDrag,
	mouseEventMove,
	mouseEventUp
};

extern void windowClosed(void * ptr);
extern void windowKeyFlagsChanged(void * ptr, int flags);
extern void windowKeyEvent(void * ptr, int type, const char * chars,
	const char * modChars, int keyCode, int modifiers);
extern void windowMouseEvent(void * ptr, int type, double x, double y);

static int generateFlags() {
	int res = 0;
	NSEventModifierFlags f = [NSEvent modifierFlags];
	if (f & NSAlternateKeyMask) {
		res |= keyFlagAlt;
	}
	if (f & NSControlKeyMask) {
		res |= keyFlagCtrl;
	}
	if (f & NSCommandKeyMask) {
		res |= keyFlagMeta;
	}
	if (f & NSShiftKeyMask) {
		res |= keyFlagShift;
	}
	return res;
}

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

- (void)flagsChanged:(NSEvent *)evt {
	windowKeyFlagsChanged((void *)self, generateFlags());
}

- (NSRect)flippedContentRect {
	NSRect frame = self.frame;
	NSRect contentRect = [self contentRectForFrameRect:frame];
	NSRect screenFrame = [self screen].frame;
	contentRect.origin.y = screenFrame.size.height -
		(contentRect.origin.y+contentRect.size.height);
	return contentRect;
}

- (id)initWithFrame:(NSRect)r {
	self = [super initWithContentRect:r
		styleMask:(NSTitledWindowMask|NSClosableWindowMask)
		backing:NSBackingStoreBuffered
		defer:NO];
	if (self) {
		[self setAcceptsMouseMovedEvents:YES];
		ContentView * cv = [[ContentView alloc]
			initWithFrame:NSMakeRect(0, 0, r.size.width, r.size.height)];
		[self setReleasedWhenClosed:NO];
		[self setContentView:cv];
		[cv release];
	}
	return self;
}

- (void)keyDown:(NSEvent *)e {
	const char * chars = e.charactersIgnoringModifiers.UTF8String;
	const char * modChars = e.characters.UTF8String;
	int modifiers = generateFlags();
	int keyCode = (int)e.keyCode;
	windowKeyEvent((void *)self, keyEventDown, chars, modChars, keyCode,
		modifiers);
}

- (void)keyUp:(NSEvent *)e {
	const char * chars = e.charactersIgnoringModifiers.UTF8String;
	const char * modChars = e.characters.UTF8String;
	int modifiers = generateFlags();
	int keyCode = (int)e.keyCode;
	windowKeyEvent((void *)self, keyEventUp, chars, modChars, keyCode,
		modifiers);
}

- (void)mouseDown:(NSEvent *)evt {
	NSPoint p = [evt locationInWindow];
	p.y = [self.contentView frame].size.height - p.y;
	windowMouseEvent((void *)self, mouseEventDown, (double)p.x, (double)p.y);
}

- (void)mouseDragged:(NSEvent *)evt {
	NSPoint p = [evt locationInWindow];
	p.y = [self.contentView frame].size.height - p.y;
	windowMouseEvent((void *)self, mouseEventDrag, (double)p.x, (double)p.y);
}

- (void)mouseMoved:(NSEvent *)evt {
	NSPoint p = [evt locationInWindow];
	p.y = [self.contentView frame].size.height - p.y;
	windowMouseEvent((void *)self, mouseEventMove, (double)p.x, (double)p.y);
}

- (void)mouseUp:(NSEvent *)evt {
	NSPoint p = [evt locationInWindow];
	p.y = [self.contentView frame].size.height - p.y;
	windowMouseEvent((void *)self, mouseEventUp, (double)p.x, (double)p.y);
}

- (void)orderOut:(id)sender {
	[super orderOut:sender];
	if (sender) {
		windowClosed((void *)self);
	}
}

- (void)setFlippedContentRect:(NSRect)r {
	NSRect screenFrame = [self screen].frame;
	r.origin.y = screenFrame.size.height - (r.origin.y+r.size.height);
	[self setFrame:[self frameRectForContentRect:r] display:YES];
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
	NSRect r = [(SimpleWindow *)ptr flippedContentRect];
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
	[(SimpleWindow *)ptr setFlippedContentRect:r];
}

void SetWindowTitle(void * w, char * title) {
	ASSERT_MAIN;
	[(NSWindow *)w setTitle:[NSString stringWithUTF8String:title]];
	free((void *)title);
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
	
	modifiers int
	pointer   unsafe.Pointer
	showing   bool
	widgets   []Widget
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
