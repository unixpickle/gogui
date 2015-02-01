// +build !darwin !cgo

package gogui

import (
	"errors"
)

var unsupportedError = errors.New("OS not supported.")

// NewCanvas creates a new canvas or fails with an error.
// The returned canvas will not be added to any window and will have a nil draw
// function by default.
func NewCanvas(r Rect) (Canvas, error) {
	return nil, unsupportedError
}

// NewWindow creates a new window or fails with an error.
// The returned window will not be shown until its Show() method is called.
func NewWindow(r Rect) (Window, error) {
	return nil, unsupportedError
}

// Main runs the main loop of the app. This should be called from the main
// function, since it may require execution on the main OS thread.
func Main(info *AppInfo) {
	select {
	}
}
