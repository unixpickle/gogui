// +build darwin,cgo

package gogui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void RunMain(void (^ block)(void));

void GetViewFrame(void * ptr, double * x, double * y, double * w, double * h) {
	NSView * v = (NSView *)ptr;
	RunMain(^{
		NSRect r = [v frame];
		*x = (double)r.origin.x;
		*y = (double)r.origin.y;
		*w = (double)r.size.width;
		*h = (double)r.size.height;
	});
}

void RemoveView(void * ptr) {
	NSView * v = (NSView *)ptr;
	RunMain(^{
		[v removeFromSuperview];
	});
}

void SetViewFrame(void * ptr, double x, double y, double w, double h) {
	NSView * v = (NSView *)ptr;
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	RunMain(^{
		[v setFrame:r];
	});
}
*/
import "C"