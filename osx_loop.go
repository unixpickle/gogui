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
	// TODO: figure out why the menu title is not bold like normal apps.
	NSMenu * menu = [[[NSApp mainMenu] itemAtIndex:0] submenu];
	[menu setTitle:self.appName];
}

@end

extern void runNextEvent();

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

void DispatchMainEvent() {
	dispatch_async(dispatch_get_main_queue(), ^{
		runNextEvent();
	});
}

*/
import "C"

import (
	"runtime"
)

func init() {
	// Make sure main.main runs on the main OS thread.
	runtime.LockOSThread()
}

// Main runs the Cocoa runloop. You must call this from main.main.
func Main(info *AppInfo) {
	C.MainLoop(C.CString(info.Name))
}

// RunOnMain runs a function on the main goroutine asynchronously using the
// dispatch_async() API.
func RunOnMain(f func()) {
	pushEvent(f)
	C.DispatchMainEvent()
}
