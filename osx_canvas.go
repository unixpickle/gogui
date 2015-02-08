// +build darwin,cgo

package gogui

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

extern void canvasDrawRect(void * window, void * canvas, void * ctx);

@interface Canvas : NSView {
}
@end

@implementation Canvas

- (void)drawRect:(NSRect)ignored {
	canvasDrawRect((void *)self.window, (void *)self,
		[[NSGraphicsContext currentContext] graphicsPort]);
}

@end

#define ASSERT_MAIN NSCAssert([NSThread isMainThread], \
	@"Call must be from main thread.")

void CanvasNeedsUpdate(void * v) {
	ASSERT_MAIN;
	[(NSView *)v setNeedsDisplay:YES];
}

void ContextFillRect(void * c, double x, double y, double w, double h) {
	CGContextFillRect((CGContextRef)c, CGRectMake((CGFloat)x, (CGFloat)y,
		(CGFloat)w, (CGFloat)h));
}

void * CreateCanvas(double x, double y, double w, double h) {
	ASSERT_MAIN;
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	return (void *)[[Canvas alloc] initWithFrame:r];
}

void DestroyCanvas(void * c) {
	ASSERT_MAIN;
	NSView * v = (NSView *)c;
	[v removeFromSuperview];
	[v release];
}

void GetViewFrame(void * v, double * x, double * y, double * w,
	double * h) {
	ASSERT_MAIN;
	NSRect r = [(NSView *)v frame];
	*x = (double)r.origin.x;
	*y = (double)r.origin.y;
	*w = (double)r.size.width;
	*h = (double)r.size.height;
}

void SetViewFrame(void * v, double x, double y, double w, double h) {
	ASSERT_MAIN;
	NSRect r = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w,
		(CGFloat)h);
	[(NSView *)v setFrame:r];
}

*/
import "C"

import (
	"runtime"
	"unsafe"
)

type canvas struct {
	handler DrawHandler
	pointer unsafe.Pointer
	parent  parentRemover
}

func NewCanvas(r Rect) (Canvas, error) {
	ptr := C.CreateCanvas(C.double(r.X), C.double(r.Y), C.double(r.Width),
		C.double(r.Height))
	res := &canvas{pointer: ptr}
	runtime.SetFinalizer(res, finalizeCanvas)
	return res, nil
}

func (c *canvas) DrawHandler() DrawHandler {
	return c.handler
}

func (c *canvas) Frame() Rect {
	var x, y, w, h C.double
	C.GetViewFrame(c.pointer, &x, &y, &w, &h)
	return Rect{float64(x), float64(y), float64(w), float64(h)}
}

func (c *canvas) NeedsUpdate() {
	C.CanvasNeedsUpdate(c.pointer)
}

func (c *canvas) Parent() Widget {
	return c.parent
}

func (c *canvas) Remove() {
	if c.parent == nil {
		return
	}
	c.parent.removeView(c)
	c.parent = nil
}

func (c *canvas) SetDrawHandler(h DrawHandler) {
	c.handler = h
}

func (c *canvas) SetFrame(r Rect) {
	C.SetViewFrame(c.pointer, C.double(r.X), C.double(r.Y),
		C.double(r.Width), C.double(r.Height))
}

func (c *canvas) setParent(p parentRemover) {
	c.parent = p
}

func (c *canvas) viewPointer() unsafe.Pointer {
	return c.pointer
}

func finalizeCanvas(c *canvas) {
	RunOnMain(func() {
		c.Remove()
		C.DestroyCanvas(c.pointer)
	})
}

type drawContext struct {
	pointer unsafe.Pointer
}

func (d *drawContext) BeginPath() {
	panic("NYI")
}

func (d *drawContext) ClosePath() {
	panic("NYI")
}

func (d *drawContext) FillEllipse(r Rect) {
	panic("NYI")
}

func (d *drawContext) FillPath() {
	panic("NYI")
}

func (d *drawContext) FillRect(r Rect) {
	C.ContextFillRect(d.pointer, C.double(r.X), C.double(r.Y),
		C.double(r.Width), C.double(r.Height))
}

func (d *drawContext) LineTo(x, y float64) {
	panic("NYI")
}

func (d *drawContext) MoveTo(x, y float64) {
	panic("NYI")
}

func (d *drawContext) SetFill(r, g, b, a float64) {
	panic("NYI")
}

func (d *drawContext) SetStroke(r, g, b, a float64) {
	panic("NYI")
}

func (d *drawContext) SetThickness(thickness float64) {
	panic("NYI")
}

func (d *drawContext) StrokeEllipse(r Rect) {
	panic("NYI")
}

func (d *drawContext) StrokePath() {
	panic("NYI")
}

func (d *drawContext) StrokeRect(r Rect) {
	panic("NYI")
}
