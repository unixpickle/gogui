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

- (BOOL)isFlipped {
	return YES;
}

@end

#define ASSERT_MAIN NSCAssert([NSThread isMainThread], \
	@"Call must be from main thread.")

void CanvasNeedsUpdate(void * v) {
	ASSERT_MAIN;
	[(NSView *)v setNeedsDisplay:YES];
}

void ContextBeginPath(void * c) {
	CGContextBeginPath((CGContextRef)c);
}

void ContextClosePath(void * c) {
	CGContextClosePath((CGContextRef)c);
}

void ContextFillEllipse(void * c, double x, double y, double w, double h) {
	CGContextFillEllipseInRect((CGContextRef)c, CGRectMake((CGFloat)x,
		(CGFloat)y, (CGFloat)w, (CGFloat)h));
}

void ContextFillPath(void * c) {
	CGContextFillPath((CGContextRef)c);
}

void ContextFillRect(void * c, double x, double y, double w, double h) {
	CGContextFillRect((CGContextRef)c, CGRectMake((CGFloat)x, (CGFloat)y,
		(CGFloat)w, (CGFloat)h));
}

void ContextLineTo(void * c, double x, double y) {
	CGContextAddLineToPoint((CGContextRef)c, (CGFloat)x, (CGFloat)y);
}

void ContextMoveTo(void * c, double x, double y) {
	CGContextMoveToPoint((CGContextRef)c, (CGFloat)x, (CGFloat)y);
}

void ContextSetFill(void * c, double r, double g, double b, double a) {
	CGContextSetRGBFillColor((CGContextRef)c, (CGFloat)r, (CGFloat)g,
		(CGFloat)b, (CGFloat)a);
}

void ContextSetStroke(void * c, double r, double g, double b, double a) {
	CGContextSetRGBStrokeColor((CGContextRef)c, (CGFloat)r, (CGFloat)g,
		(CGFloat)b, (CGFloat)a);
}

void ContextSetThickness(void * c, double thickness) {
	CGContextSetLineWidth((CGContextRef)c, (CGFloat)thickness);
}

void ContextStrokeEllipse(void * c, double x, double y, double w, double h) {
	CGContextStrokeEllipseInRect((CGContextRef)c, CGRectMake((CGFloat)x,
		(CGFloat)y, (CGFloat)w, (CGFloat)h));
}

void ContextStrokePath(void * c) {
	CGContextStrokePath((CGContextRef)c);
}

void ContextStrokeRect(void * c, double x, double y, double w, double h) {
	CGContextStrokeRect((CGContextRef)c, CGRectMake((CGFloat)x, (CGFloat)y,
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
	C.ContextBeginPath(d.pointer)
}

func (d *drawContext) ClosePath() {
	C.ContextClosePath(d.pointer)
}

func (d *drawContext) FillEllipse(r Rect) {
	C.ContextFillEllipse(d.pointer, C.double(r.X), C.double(r.Y),
		C.double(r.Width), C.double(r.Height))
}

func (d *drawContext) FillPath() {
	C.ContextFillPath(d.pointer)
}

func (d *drawContext) FillRect(r Rect) {
	C.ContextFillRect(d.pointer, C.double(r.X), C.double(r.Y),
		C.double(r.Width), C.double(r.Height))
}

func (d *drawContext) LineTo(x, y float64) {
	C.ContextLineTo(d.pointer, C.double(x), C.double(y))
}

func (d *drawContext) MoveTo(x, y float64) {
	C.ContextMoveTo(d.pointer, C.double(x), C.double(y))
}

func (d *drawContext) SetFill(r, g, b, a float64) {
	C.ContextSetFill(d.pointer, C.double(r), C.double(g), C.double(b),
		C.double(a))
}

func (d *drawContext) SetStroke(r, g, b, a float64) {
	C.ContextSetStroke(d.pointer, C.double(r), C.double(g), C.double(b),
		C.double(a))
}

func (d *drawContext) SetThickness(thickness float64) {
	C.ContextSetThickness(d.pointer, C.double(thickness))
}

func (d *drawContext) StrokeEllipse(r Rect) {
	C.ContextStrokeEllipse(d.pointer, C.double(r.X), C.double(r.Y),
		C.double(r.Width), C.double(r.Height))
}

func (d *drawContext) StrokePath() {
	C.ContextStrokePath(d.pointer)
}

func (d *drawContext) StrokeRect(r Rect) {
	C.ContextStrokeRect(d.pointer, C.double(r.X), C.double(r.Y),
		C.double(r.Width), C.double(r.Height))
}
