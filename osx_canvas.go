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
	CGContextRef c = (CGContextRef)[[NSGraphicsContext currentContext]
		graphicsPort];
	CGContextSetLineCap(c, kCGLineCapRound);
	CGContextSetLineJoin(c, kCGLineJoinRound);
	canvasDrawRect((void *)self.window, (void *)self, (void *)c);
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

void ContextText(char * text, double x, double y, double fontSize,
	char * fontName, double r, double g, double b, double a) {
	// Generate the font
	NSString * name = [NSString stringWithUTF8String:fontName];
	free((void *)fontName);
	NSFont * font = [NSFont fontWithName:name size:(CGFloat)fontSize];
	
	// Generate the color
	NSColor * color = [NSColor colorWithRed:(CGFloat)r green:(CGFloat)g
		blue:(CGFloat)b alpha:(CGFloat)a];
	
	// Generate the attributes and draw the string
	NSDictionary * dict = @{NSFontAttributeName: font,
		NSForegroundColorAttributeName: color};
	NSString * string = [NSString stringWithUTF8String:text];
	free((void *)text);
	[string drawAtPoint:NSMakePoint((CGFloat)x, (CGFloat)y)
		withAttributes:dict];
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
	fontSize  float64
	fontName  string
	fillColor Color
}

func newDrawContext(p unsafe.Pointer) *drawContext {
	return &drawContext{p, 18, "Helvetica", Color{0, 0, 0, 1}}
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

func (d *drawContext) FillText(text string, x, y float64) {
	c := d.fillColor
	C.ContextText(C.CString(text), C.double(x), C.double(y),
		C.double(d.fontSize), C.CString(d.fontName), C.double(c.R),
		C.double(c.G), C.double(c.B), C.double(c.A))
}

func (d *drawContext) LineTo(x, y float64) {
	C.ContextLineTo(d.pointer, C.double(x), C.double(y))
}

func (d *drawContext) MoveTo(x, y float64) {
	C.ContextMoveTo(d.pointer, C.double(x), C.double(y))
}

func (d *drawContext) SetFill(c Color) {
	C.ContextSetFill(d.pointer, C.double(c.R), C.double(c.G), C.double(c.B),
		C.double(c.A))
	d.fillColor = c
}

func (d *drawContext) SetFont(size float64, name string) {
	d.fontSize = size
	d.fontName = name
}

func (d *drawContext) SetStroke(c Color) {
	C.ContextSetStroke(d.pointer, C.double(c.R), C.double(c.G), C.double(c.B),
		C.double(c.A))
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
