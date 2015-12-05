package pin

import "C"

import "unsafe"

type View struct {
}

func (v *View) OnRepaint(x, y, w, h int) {
	// TOIMPL
	return
}

//export ViewOnRepaint
func ViewOnRepaint(p unsafe.Pointer, x, y, w, h int) {
	// TOIMPL
	var v = (*View)(unsafe.Pointer(p))
	v.OnRepaint(x, y, w, h)
}

//export ViewOnMouse
func ViewOnMouse(p unsafe.Pointer, mouseEvent, mouseButton int, x, y int,
	xRoot, yRoot int, modifier int, isPopupTrigger, isSynthesized bool) {
	// TOIMPL
	return
}

//export ViewOnMenu
func ViewOnMenu(p unsafe.Pointer, x, y, xAbs, yAbs int, isKeyboardTrigger bool) {
	// TOIMPL
	return
}

//export ViewOnScroll
func ViewOnScroll(p unsafe.Pointer, x, y, xAbs, yAbs, dx, dy int, modifiers int,
	lines int, chars int, defaultLines int, defaultChars int,
	xMultiplier, yMultiplier float32) {
	// TOIMPL
	return
}

//export ViewOnKey
func ViewOnKey(p unsafe.Pointer, typ, keyCode int, keyChars unsafe.Pointer,
	keyCharsCount int, modifiers int) {
	// TOIMPL
	// keyCharsAddr := reflect.SliceHeader{
	// 	Data: uintptr(unsafe.Pointer(keyChars)),
	// 	Len:  keyCharsCount,
	// 	Cap:  keyCharsCount,
	// }
	// keyCharsSlice := *(*[]C.char)(unsafe.Pointer(&keyCharsAddr))
	return
}

//export ViewOnResize
func ViewOnResize(p unsafe.Pointer, width, height int) {
	// TOIMPL
}

func ViewOnDragEnter(p unsafe.Pointer, x, y, xAbs, yAbs, recommendedDropAction int) {

}
