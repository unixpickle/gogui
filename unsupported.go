// +build !darwin !cgo

package gogui

import (
	"errors"
)

var unsupportedError = errors.New("OS not supported.")

func NewCanvas(r Rect) (Canvas, error) {
	return nil, unsupportedError
}

func NewWindow(r Rect) (Window, error) {
	return nil, unsupportedError
}

// Main does nothing.
func Main(info *AppInfo) {
	select {
	}
}
