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

void RunMain(void (^ block)(void)) {
	if ([NSThread isMainThread]) {
		block();
	} else {
		dispatch_sync(dispatch_get_main_queue(), block);
	}
}
*/
import "C"

import (
	"runtime"
	"sync"
)

var globalLock sync.Mutex

func init() {
	// Make sure main.main runs on the main OS thread.
	runtime.LockOSThread()
	go mainEventLoop.main()
}

// Main runs the Cocoa runloop. You must call this from main.main.
func Main(info *AppInfo) {
	C.MainLoop(C.CString(info.Name))
}
