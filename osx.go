// +build darwin,cgo

package gogui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSObject <NSApplicationDelegate> {
}

@property (nonatomic, strong) NSString * appName;

@end

@implementation AppDelegate

@synthesize appName;

- (void)applicationDidFinishLaunching:(NSNotification *)note {
	// TODO: figure out why the menu title is not bold like other app.
	NSMenu * menu = [[[NSApp mainMenu] itemAtIndex:0] submenu];
	[menu setTitle:self.appName];
}

@end

void CenterWindow(void * ptr) {
	NSWindow * w = (NSWindow *)ptr;
	dispatch_sync(dispatch_get_main_queue(), ^{
		[w center];
	});
}

void * CreateWindow(double x, double y, double w, double h) {
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	__block NSWindow * res = nil;
	dispatch_sync(dispatch_get_main_queue(), ^{
		res = [[NSWindow alloc] initWithContentRect:r
			styleMask:(NSTitledWindowMask|NSClosableWindowMask)
			backing:NSBackingStoreBuffered
			defer:NO];
		[res setReleasedWhenClosed:NO];
	});
	return (void *)res;
}

void GetWindowFrame(void * ptr, double * x, double * y, double * w,
	double * h) {
	// TODO: use the content frame, not the window frame.
	NSWindow * window = (NSWindow *)ptr;
	dispatch_sync(dispatch_get_main_queue(), ^{
		NSRect r = [window frame];
		*x = (double)r.origin.x;
		*y = (double)r.origin.y;
		*w = (double)r.size.width;
		*h = (double)r.size.height;
	});
}

void MainLoop(const char * name) {
	NSString * appName = [NSString stringWithUTF8String:name];
	AppDelegate * delegate = [[AppDelegate alloc] init];
	delegate.appName = appName;

	NSAutoreleasePool * pool = [[NSAutoreleasePool alloc] init];
    [NSApplication sharedApplication];

	// Make sure the application behaves normally.
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
	[[NSApplication sharedApplication] setDelegate:delegate];

	// Create the main menu.
    NSMenu * menubar = [[NSMenu alloc] init];
    NSMenuItem * appMenuItem = [[NSMenuItem alloc] init];
    [menubar addItem:appMenuItem];
    [NSApp setMainMenu:menubar];

	// Add a quit button to the menu.
    NSMenu * appMenu = [[NSMenu alloc] init];
    NSMenuItem * quitMenuItem = [[NSMenuItem alloc]
		initWithTitle:[@"Quit " stringByAppendingString:appName]
        action:@selector(terminate:) keyEquivalent:@"q"];
    [appMenu addItem:quitMenuItem];
    [appMenuItem setSubmenu:appMenu];

	// Run the loop.
    [NSApp run];

	// Release menus.
	[NSApp setMainMenu:nil];
	[menubar release];
	[appMenuItem release];
	[appMenu release];
	[quitMenuItem release];
	[pool release];
}

void SetWindowFrame(void * ptr, double x, double y, double w, double h) {
	// TODO: use the content frame, not the window frame.
	NSWindow * window = (NSWindow *)ptr;
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	dispatch_sync(dispatch_get_main_queue(), ^{
		[window setFrame:r display:YES];
	});
}

void SetWindowTitle(void * ptr, const char * title) {
	NSWindow * w = (NSWindow *)ptr;
	dispatch_sync(dispatch_get_main_queue(), ^{
		[w setTitle:[NSString stringWithUTF8String:title]];
	});
}

void ShowWindow(void * ptr) {
	NSWindow * w = (NSWindow *)ptr;
	dispatch_sync(dispatch_get_main_queue(), ^{
		[w makeKeyAndOrderFront:nil];
		[NSApp activateIgnoringOtherApps:YES];
	});
}
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

func init() {
	// Make sure main.main runs on the main OS thread.
	runtime.LockOSThread()
}

func NewCanvas(r Rect) (Canvas, error) {
	return nil, errors.New("Canvas not yet supported.")
}

type window struct {
	pointer unsafe.Pointer
}

func NewWindow(r Rect) (Window, error) {
	ptr := C.CreateWindow(C.double(r.Y), C.double(r.Y), C.double(r.Width),
		C.double(r.Height))
	return &window{ptr}, nil
}

func (w *window) Add(widget Widget) {
	// TODO: this
}

func (w *window) Center() {
	C.CenterWindow(w.pointer)
}

func (w *window) Children() []Widget {
	// TODO: this
	return []Widget{}
}

func (w *window) Destroy() {
	// TODO: this
}

func (w *window) Focus() {
	// TODO: this
}

func (window *window) Frame() Rect {
	var x, y, w, h C.double
	C.GetWindowFrame(window.pointer, &x, &y, &w, &h)
	return Rect{float64(x), float64(y), float64(w), float64(h)}
}

func (w *window) Hide() {
	// TODO: this
}

func (w *window) Parent() {
	// TODO: this
}

func (w *window) Remove() {
	// TODO: this
}

func (w *window) SetFrame(r Rect) {
	// TODO: this
}

func (w *window) SetTitle(title string) {
	C.SetWindowTitle(w.pointer, C.CString(title))
}

func (w *window) Show() {
	C.ShowWindow(w.pointer)
}

// Main runs the Cocoa run-loop. You must call this from main.main.
func Main(info *AppInfo) {
	C.MainLoop(C.CString(info.Name))
}
