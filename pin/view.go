package pin

import "C"

import "unsafe"

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

//export ViewOnScroll
func ViewOnScroll(p unsafe.Pointer, x, y, xAbs, yAbs, dx, dy int, modifiers int,
	lines int, chars int, defaultLines int, defaultChars int,
	xMultiplier, yMultiplier float32) {
	return
}

//export ViewOnKey
func ViewOnKey(p unsafe.Pointer, typ, keyCode int, keyChars unsafe.Pointer,
	keyCharsCount int, modifiers int) {
	// addr := reflect.SliceHeader{
	// 	Data: uintptr(unsafe.Pointer(keyChars)),
	// 	Len:  keyCharsCount,
	// 	Cap:  keyCharsCount,
	// }
	// keyCharsSlice := *(*[]C.char)(unsafe.Pointer(&addr))
	return
}

//export ViewOnResize
func ViewOnResize(p unsafe.Pointer, width, height int) {

}
