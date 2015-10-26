package pin

// runtime.SetFinalizer(ptr, finalizerFunc)

// #cgo LDFLAGS: -lX11
// #include <stdlib.h>
// #include "bridge.h"
import "C"

import "unsafe"

type Window struct {
	isEnabled bool
}

func NewWindow() *Window {
	var w = new(Window)
	C.NewWindow(unsafe.Pointer(w), unsafe.Pointer(nil), unsafe.Pointer(nil), 0, 0)
	return w
}

func (w *Window) StateOnChange(state int) {
	return
}

func (w *Window) IsEnabled() bool {
	return w.isEnabled
}

//export WindowIsEnabled
func WindowIsEnabled(p unsafe.Pointer) bool {
	var w = (*Window)(unsafe.Pointer(p))
	return w.IsEnabled()
}

//export WindowOnStateChange
func WindowOnStateChange(p unsafe.Pointer, state C.int) {
	var w = (*Window)(unsafe.Pointer(p))
	w.StateOnChange(int(state))
}

//export WindowOnNotifyFocus
func WindowOnNotifyFocus(p unsafe.Pointer, state C.int) {
	return
}

//export WindowOnFocusDisabled
func WindowOnFocusDisabled(p unsafe.Pointer) {
	return
}

//export WindowOnDestroy
func WindowOnDestroy(p unsafe.Pointer) {
	return
}

//export WindowOnClose
func WindowOnClose(p unsafe.Pointer) {
	return
}
