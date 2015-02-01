// +build darwin,cgo

package gogui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

extern void canvasDraw(unsigned long long i);
void RunMain(void (^ block)(void));

@interface AppDelegate : NSObject <NSApplicationDelegate> {
}

@property (nonatomic, strong) NSString * appName;

@end

@implementation AppDelegate

@synthesize appName;

- (void)applicationDidFinishLaunching:(NSNotification *)note {
	// TODO: figure out why the menu title is not bold like normal apps.
	NSMenu * menu = [[[NSApp mainMenu] itemAtIndex:0] submenu];
	[menu setTitle:self.appName];
}

@end

@interface Canvas : NSView {
}

@property (readwrite) unsigned long long canvasId;

@end

@implementation Canvas

@synthesize canvasId;

- (void)drawRect:(NSRect)dirtyRect {
	canvasDraw(self.canvasId);
}

@end

void AddToWindow(void * wind, void * view) {
	NSWindow * w = (NSWindow *)wind;
	NSView * v = (NSView *)view;
	RunMain(^{
		[w.contentView addSubview:v];
	});
}

void CanvasNeedsDisplay(void * ptr) {
	Canvas * c = (Canvas *)ptr;
	RunMain(^{
		[c setNeedsDisplay:YES];
	});
}

void CenterWindow(void * ptr) {
	NSWindow * w = (NSWindow *)ptr;
	RunMain(^{
		[w center];
	});
}

void * CreateCanvas(unsigned long long ident, double x, double y, double w,
	double h) {
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	__block Canvas * res = nil;
	RunMain(^{
		res = [[Canvas alloc] initWithFrame:r];
		res.canvasId = ident;
	});
	return (void *)res;
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

void DestroyCanvas(void * ptr) {
	Canvas * c = (Canvas *)ptr;
	RunMain(^{
		[c removeFromSuperview];
		[c release];
	});
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

void GetCanvasFrame(void * ptr, double * x, double * y, double * w,
	double * h) {
	Canvas * c = (Canvas *)ptr;
	RunMain(^{
		NSRect r = [c frame];
		*x = (double)r.origin.x;
		*y = (double)r.origin.y;
		*w = (double)r.size.width;
		*h = (double)r.size.height;
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

void RemoveCanvas(void * ptr) {
	Canvas * c = (Canvas *)ptr;
	RunMain(^{
		[c removeFromSuperview];
	});
}

void RunMain(void (^ block)(void)) {
	if ([NSThread isMainThread]) {
		block();
	} else {
		dispatch_sync(dispatch_get_main_queue(), block);
	}
}

void SetCanvasFrame(void * ptr, double x, double y, double w, double h) {
	Canvas * c = (Canvas *)ptr;
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	RunMain(^{
		[c setFrame:r];
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
