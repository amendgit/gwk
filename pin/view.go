package pin

import "C"

import (
	"unsafe"
)

type View struct {
}

func (v *View) OnRepaint(x, y, w, h int) {
	return
}

//export ViewOnRepaint
func ViewOnRepaint(p unsafe.Pointer, x, y, w, h int) {
	var v = (*View)(unsafe.Pointer(p))
	v.OnRepaint(x, y, w, h)
}

//export ViewOnMouse
func ViewOnMouse(p unsafe.Pointer, mouseEvent, mouseButton int, x, y int,
	xRoot, yRoot int, modifier int, isPopupTrigger, isSynthesized bool) {
	return
}

//export ViewOnMenu
func ViewOnMenu(p unsafe.Pointer, x, y, xAbs, yAbs int, isKeyboardTrigger bool) {
	return
}
