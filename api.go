// package gogui provides a very simple library for creating user interfaces in
// Go.
package gogui

// The AppInfo object represents information about the application which the
// implementation may choose to display to the user in some form.
type AppInfo struct {
	Name string
}

// A Canvas is a widget that can be drawn into.
type Canvas interface {
	Widget
	
	DrawHandler() DrawHandler
	SetDrawHandler(d DrawHandler)
}

// A DrawContext receives draw commands.
type DrawContext interface {
	// BeginPath starts a path which can be filled or stroked.
	BeginPath()

	// ClosePath closes the current path by connecting the first and last points
	// in it.
	ClosePath()

	// FillEllipse fills an ellipse inside a rectangle.
	FillEllipse(r Rect)

	// FillPath fills the current path.
	// If the path was not closed, the behaviour of FillPath may vary on
	// different platforms.
	FillPath()

	// FillRect fills a rectangle.
	FillRect(r Rect)

	// LineTo adds a line from the current point in the path to another point.
	LineTo(x, y float64)

	// MoveTo moves the current path to a point without connecting said point to
	// the rest of the path.
	MoveTo(x, y float64)

	// SetFill sets the color for every Fill method.
	SetFill(r, g, b, a float64)

	// SetStroke sets the color for every Stroke method.
	SetStroke(r, g, b, a float64)
	
	// SetThickness sets the thickness for every Stroke method.
	SetThickness(thickness float64)

	// StrokeEllipse strokes an ellipse inside a rectangle.
	StrokeEllipse(r Rect)

	// StrokePath outlines the current path.
	StrokePath()

	// StrokeRect outlines a rectangle.
	StrokeRect(r Rect)
}

// A DrawHandler is called to draw into a canvas's drawing context.
type DrawHandler func(DrawContext)

// A KeyEvent holds information for a key event.
type KeyEvent struct {
	// CharCode stores a char code similar to the char codes in JavaScript.
	CharCode int
	
	// KeyCode stores the OS-specific key code (if available).
	KeyCode  int
	
	AltKey   bool
	CtrlKey  bool
	MetaKey  bool
	ShiftKey bool
}

// A KeyEventer listens to key events and sends them to handlers.
type KeyEventer interface {
	KeyDownHandler() KeyHandler
	KeyPressHandler() KeyHandler
	KeyUpHandler() KeyHandler
	SetKeyDownHandler(k KeyHandler)
	SetKeyPressHandler(k KeyHandler)
	SetKeyUpHandler(k KeyHandler)
}

// A KeyHandler handles keyboard events.
type KeyHandler func(KeyEvent)

// A MouseEvent holds information for a mouse event.
type MouseEvent struct {
	X float64
	Y float64
}

// A MouseEventer listens to mouse events and sends them to handlers.
type MouseEventer interface {
	MouseDownHandler() MouseHandler
	MouseDragHandler() MouseHandler
	MouseMoveHandler() MouseHandler
	MouseUpHandler() MouseHandler
	SetMouseDownHandler(m MouseHandler)
	SetMouseDragHandler(m MouseHandler)
	SetMouseMoveHandler(m MouseHandler)
	SetMouseUpHandler(m MouseHandler)
}

// A MouseHandler handles mouse events.
type MouseHandler func(MouseEvent)

// A Rect holds the position and dimensions for a Widget.
//
// The X value starts from the left of the parent. The Y value starts from the
// top of the parent.
//
// Width extends to the right and Height extends downward.
type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// A Widget is any item that can be shown visually to the user.
type Widget interface {
	// Frame returns the bounding box for this widget.
	Frame() Rect

	// Parent returns the widget which contains this widget, or nil.
	Parent() Widget

	// Remove removes this widget from its parent if it has one.
	Remove()

	// SetFrame sets the bounding box for this widget.
	SetFrame(r Rect)
}

// A Window is container Widget which shows the user its sub-widgets.
type Window interface {
	KeyEventer
	MouseEventer
	
	// Add adds a widget to the window. The widget cannot already be added to
	// something else.
	Add(w Widget)

	// Center centers the window on the screen.
	Center()

	// Children returns every direct child of this window.
	Children() []Widget
	
	// CloseHandler returns the window's close handler.
	CloseHandler() func()

	// Focus brings the window to the front if it is showing.
	Focus()

	// Frame returns the content rectangle for the window.
	Frame() Rect

	// Hide closes the window if it was open.
	Hide()

	// Parent returns nil; it exists to implement the Widget interface.
	Parent() Widget

	// Remove does nothing; it exists to implement the Widget interface.
	Remove()

	// SetCloseHandler sets a function to be called when the user closes the
	// window.
	SetCloseHandler(h func())

	// SetFrame sets the content rectangle for the window.
	SetFrame(r Rect)
	
	// SetTitle sets the title of the window.
	SetTitle(t string)

	// Show opens the window if it was not open before.
	Show()
	
	// Showing returns whether the window is showing or not.
	Showing() bool
}
