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

enum {
	canvasCommandBeginPath = 0,
	canvasCommandClosePath,
	canvasCommandFillPath,
	canvasCommandFillRect
};

@interface Canvas : NSView {
	int numCalls;
	int * callNames;
	double * callArgs;
}

- (void)applyCalls:(int)num names:(int *)names args:(double *)args;

@end

@implementation Canvas

- (void)applyCalls:(int)num names:(int *)names args:(double *)args {
	if (callNames) {
		free(callNames);
	}
	if (callArgs) {
		free(callArgs);
	}
	numCalls = num;
	callNames = (int *)malloc(sizeof(int) * num);
	callArgs = (double *)malloc(sizeof(double) * num * 4);
	memcpy(callNames, names, sizeof(int)*num);
	memcpy(callArgs, args, sizeof(double)*num*4);
	[self setNeedsDisplay:YES];
}

- (void)dealloc {
	[super dealloc];
	if (callNames) {
		free(callNames);
	}
	if (callArgs) {
		free(callArgs);
	}
}

- (void)drawRect:(NSRect)dirtyRect {
	if (numCalls == 0) {
		return;
	}
	CGContextRef c = (CGContextRef)[[NSGraphicsContext currentContext]
		graphicsPort];
	for (int i = 0; i < numCalls; ++i) {
		// Get the info for the call.
		int command = callNames[i];
		double * argsPtr = &callArgs[i * 4];
		CGFloat args[4] = {(CGFloat)argsPtr[0], (CGFloat)argsPtr[1],
			(CGFloat)argsPtr[2], (CGFloat)argsPtr[3]};

		// TODO: add the rest of the commands
		switch (command) {
		case canvasCommandFillRect:
			CGContextFillRect(c, CGRectMake(args[0], args[1], args[2],
				args[3]));
			break;
		default:
			break;
		}
	}
}

- (BOOL)isFlipped {
	return YES;
}

@end

void ApplyCalls(void * ptr, int count, int * commands, double * args) {
	Canvas * c = (Canvas *)ptr;
	RunMain(^{
		[c applyCalls:count names:commands args:args];
	});
}

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

import (
	"runtime"
	"unsafe"
)

type canvas struct {
	parent  Widget
	pointer unsafe.Pointer
	
	commands []C.int
	args     []C.double
}

const (
	canvasCommandBeginPath = iota
	canvasCommandClosePath = iota
	canvasCommandFillPath = iota
	canvasCommandFillRect = iota
)

// NewCanvas creates a new canvas with the given frame.
func NewCanvas(r Rect) (Canvas, error) {
	globalLock.Lock()
	defer globalLock.Unlock()
	ptr := C.CreateCanvas(C.double(r.X), C.double(r.Y), C.double(r.Width),
		C.double(r.Height))
	res := &canvas{nil, ptr, []C.int{}, []C.double{}}
	runtime.SetFinalizer(res, finalizeCanvas)
	return res, nil
}

func (c *canvas) BeginPath() {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	c.addEmptyCommand(canvasCommandBeginPath)
}

func (c *canvas) ClosePath() {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	c.addEmptyCommand(canvasCommandClosePath)
}

func (c *canvas) FillPath() {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	c.addEmptyCommand(canvasCommandFillPath)
}

func (c *canvas) FillRect(r Rect) {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	c.addFullCommand(canvasCommandFillRect, r.X, r.Y, r.Width, r.Height)
}

func (c *canvas) Flush() {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	C.ApplyCalls(c.pointer, C.int(len(c.commands)), &c.commands[0], &c.args[0])
	c.commands = []C.int{}
	c.args = []C.double{}
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

func (c *canvas) LineTo(x, y float64) {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	// TODO: this
}

func (c *canvas) MoveTo(x, y float64) {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	// TODO: this
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
		w.removeWidget(c)
	} else {
		panic("Unknown parent type.")
	}
	c.parent = nil
	
	// Remove the actual view
	C.RemoveView(c.pointer)
}

func (c *canvas) SetFill(r, g, b, a float64) {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	// TODO: this
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

func (c *canvas) SetStroke(r, g, b, a float64) {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	// TODO: this
}

func (c *canvas) SetThickness(thickness float64) {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	// TODO: this
}

func (c *canvas) StrokePath() {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	// TODO: this
}

func (c *canvas) StrokeRect(r Rect) {
	globalLock.Lock();
	defer globalLock.Unlock();
	if c.pointer == nil {
		panic("Canvas is invaild.")
	}
	// TODO: this
}

func (c *canvas) addEmptyCommand(cmd int) {
	c.addFullCommand(cmd, 0, 0, 0, 0)
}

func (c *canvas) addFullCommand(cmd int, w, x, y, z float64) {
	c.commands = append(c.commands, C.int(cmd))
	c.args = append(c.args, C.double(w), C.double(x), C.double(y), C.double(z))
}

func finalizeCanvas(c *canvas) {
	globalLock.Lock()
	defer globalLock.Unlock()
	if c.pointer == nil {
		panic("Canvas is invalid.")
	}
	
	// I do not call c.Remove() here because the finalizer will only be called
	// if nothing, including a superview, contains the Canvas.
	
	C.DestroyCanvas(c.pointer)
	c.pointer = nil
}
