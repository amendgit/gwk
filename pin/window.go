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
	// TOIMPL
	var w = new(Window)
	C.NewWindow(unsafe.Pointer(w), unsafe.Pointer(nil), unsafe.Pointer(nil), 0, 0)
	return w
}

func (w *Window) StateOnChange(state int) {
	// TOIMPL
	return
}

func (w *Window) IsEnabled() bool {
	// TOIMPL
	return w.isEnabled
}

//export WindowIsEnabled
func WindowIsEnabled(p unsafe.Pointer) bool {
	// TOIMPL
	var w = (*Window)(unsafe.Pointer(p))
	return w.IsEnabled()
}

//export WindowOnStateChange
func WindowOnStateChange(p unsafe.Pointer, state C.int) {
	// TOIMPL
	var w = (*Window)(unsafe.Pointer(p))
	w.StateOnChange(int(state))
}

//export WindowOnNotifyFocus
func WindowOnNotifyFocus(p unsafe.Pointer, state C.int) {
	// TOIMPL
	return
}

//export WindowOnFocusDisabled
func WindowOnFocusDisabled(p unsafe.Pointer) {
	// TOIMPL
	return
}

//export WindowOnDestroy
func WindowOnDestroy(p unsafe.Pointer) {
	// TOIMPL
	return
}

//export WindowOnClose
func WindowOnClose(p unsafe.Pointer) {
	// TOIMPL
	return
}

//export WindowOnFocusUngrab
func WindowOnFocusUngrab(p unsafe.Pointer) {
	// TOIMPL
}
