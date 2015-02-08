// +build darwin,cgo

package gogui

import (
	"unsafe"
)

type keyEvents struct {
	down  KeyHandler
	press KeyHandler
	up    KeyHandler
}

func (k *keyEvents) KeyDownHandler() KeyHandler {
	return k.down
}

func (k *keyEvents) KeyPressHandler() KeyHandler {
	return k.press
}

func (k *keyEvents) KeyUpHandler() KeyHandler {
	return k.up
}

func (k *keyEvents) SetKeyDownHandler(h KeyHandler) {
	k.down = h
}

func (k *keyEvents) SetKeyPressHandler(h KeyHandler) {
	k.press = h
}

func (k *keyEvents) SetKeyUpHandler(h KeyHandler) {
	k.up = h
}

type mouseEvents struct {
	down MouseHandler
	drag MouseHandler
	move MouseHandler
	up   MouseHandler
}

func (m *mouseEvents) MouseDownHandler() MouseHandler {
	return m.down
}

func (m *mouseEvents) MouseDragHandler() MouseHandler {
	return m.drag
}

func (m *mouseEvents) MouseMoveHandler() MouseHandler {
	return m.move
}

func (m *mouseEvents) MouseUpHandler() MouseHandler {
	return m.up
}

func (m *mouseEvents) SetMouseDownHandler(h MouseHandler) {
	m.down = h
}

func (m *mouseEvents) SetMouseDragHandler(h MouseHandler) {
	m.drag = h
}

func (m *mouseEvents) SetMouseMoveHandler(h MouseHandler) {
	m.move = h
}

func (m *mouseEvents) SetMouseUpHandler(h MouseHandler) {
	m.up = h
}

type parentRemover interface {
	Widget
	removeView(v ptrView)
}

type ptrView interface {
	viewPointer() unsafe.Pointer
	setParent(p parentRemover)
}

type windowEvents struct {
	keyEvents
	mouseEvents
	onClose func()
}

func (w *windowEvents) CloseHandler() func() {
	return w.onClose
}

func (w *windowEvents) SetCloseHandler(h func()) {
	w.onClose = h
}
