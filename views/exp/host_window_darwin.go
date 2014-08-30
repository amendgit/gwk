package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

extern void gwk_window_run(void *slf);
extern void gwk_window_set_delegate(void *slf, void *delegate);
extern void *new_gwk_window();
*/
import "C"

import (
	"unsafe"
)

func window_will_use_fullscreen_content_rect(slf, window, content_rect uintptr) {
	//host_window := (*HostWindow)(unsafe.Pointer(slf))
}

type HostWindow struct {
	gwk_window unsafe.Pointer
}

func NewHostWindow() *HostWindow {
	h := new(HostWindow)
	h.gwk_window = unsafe.Pointer(C.new_gwk_window())
	C.gwk_window_set_delegate(h.gwk_window, unsafe.Pointer(h))
	return h
}

func (h *HostWindow) Run() {
	C.gwk_window_run(h.gwk_window)
}

func main() {
	h := NewHostWindow()
	h.Run()
}
