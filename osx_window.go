// +build darwin,cgo

package gogui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void RunMain(void (^ block)(void));
extern void windowKeyDown(void * ptr, int charCode, int keyCode, int flags);
extern void windowKeyPress(void * ptr, int charCode, int keyCode, int flags);
extern void windowKeyUp(void * ptr, int charCode, int keyCode, int flags);
extern void windowMouseDown(void * ptr, double x, double y);
extern void windowMouseDrag(void * ptr, double x, double y);
extern void windowMouseMove(void * ptr, double x, double y);
extern void windowMouseUp(void * ptr, double x, double y);
extern void windowOrderedOut(void * ptr);

static int eventCharCode(NSEvent * e, BOOL press) {
	NSString * s;
	if (press) {
		s = e.characters;
		// Set the case (because it's not set for us)
		if ([NSEvent modifierFlags] & NSShiftKeyMask) {
			s = [s uppercaseString];
		}
	} else {
		s = e.charactersIgnoringModifiers;
		// JavaScript char codes are always uppercase by default, but NSEvents
		// are lowercase.
		s = [s uppercaseString];
	}
	if (s.length == 0) {
		return 0;
	}
	return (int)[s characterAtIndex:0];
}

static int generateFlags() {
	int res = 0;
	NSEventModifierFlags f = [NSEvent modifierFlags];
	if (f & NSAlternateKeyMask) {
		res |= 1;
	}
	if (f & NSControlKeyMask) {
		res |= 2;
	}
	if (f & NSCommandKeyMask) {
		res |= 4;
	}
	if (f & NSShiftKeyMask) {
		res |= 8;
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

- (void)flagsChanged:(NSEvent *)evt {
	NSEventModifierFlags f = evt.modifierFlags;
	NSEventModifierFlags changed = f ^ flags;
	if (changed & NSShiftKeyMask) {
		if (f & NSShiftKeyMask) {
			windowKeyDown((void *)self, 0x10, -1, generateFlags());
		} else {
			windowKeyUp((void *)self, 0x10, -1, generateFlags());
		}
	}
	if (changed & NSAlternateKeyMask) {
		if (f & NSAlternateKeyMask) {
			windowKeyDown((void *)self, 18, -1, generateFlags());
		} else {
			windowKeyUp((void *)self, 18, -1, generateFlags());
		}
	}
	if (changed & NSCommandKeyMask) {
		if (f & NSCommandKeyMask) {
			windowKeyDown((void *)self, 91, -1, generateFlags());
		} else {
			windowKeyUp((void *)self, 91, -1, generateFlags());
		}
	}
	if (changed & NSControlKeyMask) {
		if (f & NSControlKeyMask) {
			windowKeyDown((void *)self, 17, -1, generateFlags());
		} else {
			windowKeyUp((void *)self, 17, -1, generateFlags());
		}
	}
	flags = f;
}

- (void)keyDown:(NSEvent *)evt {
	windowKeyDown((void *)self, eventCharCode(evt, NO), (int)[evt keyCode],
		generateFlags());
	windowKeyPress((void *)self, eventCharCode(evt, YES), (int)[evt keyCode],
		generateFlags());
}

- (void)keyUp:(NSEvent *)evt {
	windowKeyUp((void *)self, eventCharCode(evt, NO), (int)[evt keyCode],
		generateFlags());
}

- (void)mouseDown:(NSEvent *)evt {
	NSPoint p = [evt locationInWindow];
	p.y = [self.contentView frame].size.height - p.y;
	windowMouseDown((void *)self, (double)p.x, (double)p.y);
}

- (void)mouseDragged:(NSEvent *)evt {
	NSPoint p = [evt locationInWindow];
	p.y = [self.contentView frame].size.height - p.y;
	windowMouseDrag((void *)self, (double)p.x, (double)p.y);
}

- (void)mouseMoved:(NSEvent *)evt {
	NSPoint p = [evt locationInWindow];
	p.y = [self.contentView frame].size.height - p.y;
	windowMouseMove((void *)self, (double)p.x, (double)p.y);
}

- (void)mouseUp:(NSEvent *)evt {
	NSPoint p = [evt locationInWindow];
	p.y = [self.contentView frame].size.height - p.y;
	windowMouseUp((void *)self, (double)p.x, (double)p.y);
}

- (void)orderOut:(id)sender {
	[super orderOut:sender];
	if (sender) {
		windowOrderedOut((void *)self);
	}
}

@end

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
		res = [[SimpleWindow alloc] initWithFrame:r];
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

void SetWindowMouseMove(void * ptr, int flag) {
	NSWindow * w = (NSWindow *)ptr;
	BOOL b = (flag != 0);
	RunMain(^{
		[w setAcceptsMouseMovedEvents:b];
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

import (
	"runtime"
	"unsafe"
)

var showingWindows = []Window{}

type window struct {
	pointer unsafe.Pointer
	widgets []Widget
	showing bool
	
	onClose     func()
	onKeyDown   KeyHandler
	onKeyPress  KeyHandler
	onKeyUp     KeyHandler
	onMouseDown MouseHandler
	onMouseDrag MouseHandler
	onMouseMove MouseHandler
	onMouseUp   MouseHandler
}

func NewWindow(r Rect) (Window, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	ptr := C.CreateWindow(C.double(r.Y), C.double(r.Y), C.double(r.Width),
		C.double(r.Height))
	res := &window{pointer: ptr, widgets: []Widget{}}
	runtime.SetFinalizer(res, finalizeWindow)
	return res, nil
}

func ShowingWindows() []Window {
	globalLock.Lock()
	defer globalLock.Unlock()
	res := make([]Window, len(showingWindows))
	copy(res, showingWindows)
	return res
}

func (w *window) Add(widget Widget) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if canvas, ok := widget.(*canvas); ok {
		if canvas.parent != nil {
			panic("Widget already has a parent.")
		}
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

func (w *window) CloseHandler() func() {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	return w.onClose
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
	if w.showing {
		w.showing = false
		C.HideWindow(w.pointer)
		for i, x := range showingWindows {
			if x.(*window) == w {
				showingWindows[i] = showingWindows[len(showingWindows) - 1]
				showingWindows[len(showingWindows) - 1] = nil
				showingWindows = showingWindows[0 : len(showingWindows)-1]
				break
			}
		}
	}
}

func (w *window) KeyDownHandler() KeyHandler {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	return w.onKeyDown
}

func (w *window) KeyPressHandler() KeyHandler {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	return w.onKeyPress
}

func (w *window) KeyUpHandler() KeyHandler {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	return w.onKeyUp
}

func (w *window) MouseDownHandler() MouseHandler {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	return w.onMouseDown
}

func (w *window) MouseDragHandler() MouseHandler {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	return w.onMouseDrag
}

func (w *window) MouseMoveHandler() MouseHandler {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	return w.onMouseMove
}

func (w *window) MouseUpHandler() MouseHandler {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	return w.onMouseUp
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

func (w *window) SetCloseHandler(f func()) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	w.onClose = f
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

func (w *window) SetKeyDownHandler(f KeyHandler) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	w.onKeyDown = f
}

func (w *window) SetKeyPressHandler(f KeyHandler) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	w.onKeyPress = f
}

func (w *window) SetKeyUpHandler(f KeyHandler) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	w.onKeyUp = f
}

func (w *window) SetMouseDownHandler(f MouseHandler) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	w.onMouseDown = f
}

func (w *window) SetMouseDragHandler(f MouseHandler) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	w.onMouseDrag = f
}

func (w *window) SetMouseMoveHandler(f MouseHandler) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	w.onMouseMove = f
	
	// If nothing is listening to mouse moved events, we shouldn't waste CPU
	// handling them since they're so plentiful.
	if f != nil {
		C.SetWindowMouseMove(w.pointer, C.int(1))
	} else {
		C.SetWindowMouseMove(w.pointer, C.int(0))
	}
}

func (w *window) SetMouseUpHandler(f MouseHandler) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	w.onMouseUp = f
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
	if !w.showing {
		w.showing = true
		C.ShowWindow(w.pointer)
		showingWindows = append(showingWindows, w)
	}
}

// Showing returns whether the window is showing or not.
func (w *window) Showing() bool {
	globalLock.Lock()
	defer globalLock.Unlock()
	if w.pointer == nil {
		panic("Window is invalid.")
	}
	return w.showing
}

func (w *window) removeWidget(widget Widget) {
	for i, x := range w.widgets {
		if x == widget {
			w.widgets[i] = w.widgets[len(w.widgets) - 1]
			w.widgets[len(w.widgets) - 1] = nil
			w.widgets = w.widgets[0 : len(w.widgets)-1]
			break
		}
	}
}

func finalizeWindow(w *window) {
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
