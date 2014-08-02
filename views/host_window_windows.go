// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	"fmt"
	. "gwk/sysc"
	. "gwk/vango"
	"image"
	"log"
	"syscall"
	"unsafe"
)

type HostWindow struct {
	root_view *RootView
	hwnd      Handle
	handled   bool
	bounds    image.Rectangle
}

func NewHostWindow(bounds image.Rectangle) *HostWindow {
	var hw HostWindow
	hw.Init(bounds)
	return &hw
}

func (hw *HostWindow) Init(bounds image.Rectangle) {
	host_window_class := fmt.Sprintf("www.ustc.edu.cn/gwk/host_window/%p", hw)

	var wc WNDCLASSEX

	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.Style = CS_HREDRAW | CS_VREDRAW
	wc.FnWndProc = syscall.NewCallback(host_window_wnd_proc)
	wc.ClassExtra = 0
	wc.WindowExtra = 0
	wc.HInstance = NULL
	wc.HIcon, _ = LoadIcon(NULL, MAKEINTRESOURCE(IDI_APPLICATION))
	wc.HCursor, _ = LoadCursor(NULL, MAKEINTRESOURCE(IDC_ARROW))
	wc.HbrBackground = COLOR_WINDOWFRAME
	wc.MenuName = nil
	wc.ClassName = syscall.StringToUTF16Ptr(host_window_class)

	if _, err := RegisterClassEx(&wc); err != nil {
		log.Panicf("RegisterClassEx Failed %v", err)
	}

	var x, y, width, height int

	if bounds.Empty() {
		x, y = CW_USEDEFAULT, CW_USEDEFAULT
		width, height = CW_USEDEFAULT, CW_USEDEFAULT
	} else {
		x, y = bounds.Min.X, bounds.Min.Y
		width, height = bounds.Dx(), bounds.Dy()
	}

	hwnd, err := CreateWindowEx(WS_OVERLAPPED,
		syscall.StringToUTF16Ptr(host_window_class),
		nil,
		WS_OVERLAPPEDWINDOW,
		x, y, width, height,
		NULL, NULL, 0,
		uintptr(unsafe.Pointer(hw)))

	if err != nil {
		log.Panicf("CreateWindowEx Failed: %v", err)
	}

	hw.hwnd = hwnd
	hw.bounds = bounds
}

func (h *HostWindow) set_msg_handled(handled bool) {
	h.handled = handled
}

func (h *HostWindow) is_msg_handled() bool {
	return h.handled
}

func (h *HostWindow) on_wnd_proc(msg uint32, warg uintptr, larg uintptr) uintptr {
	switch msg {
	case WM_MOUSEMOVE:
		h.set_msg_handled(true)
		x := int(larg & 0x0000ffff)
		y := int((larg & 0xffff0000) >> 16)
		h.root_view.DispatchMouseMove(image.Pt(x, y))
		if h.is_msg_handled() {
			return TRUE
		}
	case WM_PAINT:
		h.set_msg_handled(true)
		var r RECT
		GetClientRect(h.hwnd, &r)
		if r.Right-r.Left != uint32(h.root_view.W()) &&
			r.Bottom-r.Top != uint32(h.root_view.H()) {
			// If the size of the window differs from the size of the root view
			// it means we're being asked to paint before we've gotten a
			// WM_SIZE. This can happen when the user is interactively resizing
			// the window. To avoid mass flickering we don't do anything here.
			// Once we get the WM_SIZE we'll reset the region of the window
			// which triggers another WM_NCPAINT and all is well.
			return TRUE
		}

		var dirtyRect image.Rectangle

		if GetUpdateRect(h.hwnd, &r, FALSE) != FALSE {
			dirtyRect = image.Rect(int(r.Left), int(r.Top), int(r.Right),
				int(r.Bottom))
		}

		// Why I only can get the client rect during the WM_PAINT. Otherwise
		// I'll get the whole window bounds.
		GetClientRect(h.hwnd, &r)
		h.bounds = image.Rect(int(r.Left), int(r.Top), int(r.Right),
			int(r.Bottom))

		var ps PAINTSTRUCT
		var dc = BeginPaint(h.hwnd, &ps)
		h.OnHostPaint(NativeContext(dc), dirtyRect)
		EndPaint(h.hwnd, &ps)

		if h.is_msg_handled() {
			// Handle this msg so that the window will not draw background
			// before draw the content. So that the window won't flicker.
			return TRUE
		}

	case WM_SIZE:
		h.set_msg_handled(true)

		var r RECT
		GetClientRect(h.hwnd, &r)

		rc := image.Rect(int(r.Left), int(r.Top), int(r.Right), int(r.Bottom))

		// GetWindowRect return the bounds based on screen coordinate. Which
		// means the top_left point may not (0,0). But gwk need the bounds based
		// on app coordinate.
		rc = rc.Sub(rc.Min)
		h.root_view.SetBounds(rc)
		h.root_view.DispatchLayout()

		RedrawWindow(h.hwnd, &r, NULL, RDW_INVALIDATE|RDW_ALLCHILDREN)

		if h.is_msg_handled() {
			return TRUE
		}

	case WM_ERASEBKGND:
		// Do not allow erase background. Which may caused flickering.
		return TRUE
	}

	return DefWindowProc(h.hwnd, msg, warg, larg)
}

func host_window_wnd_proc(hwnd Handle, msg uint32, warg uintptr, larg uintptr) uintptr {
	if msg == WM_NCCREATE {
		var cs = (*CREATESTRUCT)(unsafe.Pointer(larg))
		var w = (*HostWindow)(unsafe.Pointer(cs.CreateParams))
		SetWindowLongPtr(hwnd, GWLP_USERDATA, uintptr(unsafe.Pointer(w)))
		return TRUE
	}

	var lp uintptr
	if lp = GetWindowLongPtr(hwnd, GWLP_USERDATA); lp == NULL {
		return DefWindowProc(hwnd, msg, warg, larg)
	}

	var host = (*HostWindow)(unsafe.Pointer(lp))

	return host.on_wnd_proc(msg, warg, larg)
}

func (h *HostWindow) Run() {
	var msg MSG
	for GetMessage(&msg, h.hwnd, 0, 0) != 0 {
		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}
}

func (hw *HostWindow) Show() {
	ShowWindow(hw.hwnd, SW_SHOWNORMAL)
}

func (hw *HostWindow) ClientBounds() image.Rectangle {
	// var r RECT
	// var err = GetClientRect(hw.hwnd, &r)
	// return image.Rect(int(r.Left), int(r.Top), int(r.Right), int(r.Bottom))
	return hw.bounds
}

func (hw *HostWindow) Bounds() image.Rectangle {
	var r RECT
	GetWindowRect(hw.hwnd, &r)
	return image.Rect(int(r.Left), int(r.Top), int(r.Right), int(r.Bottom))
}

func (h *HostWindow) UpdateWindow() {
	UpdateWindow(h.hwnd)
}

func (h *HostWindow) InvalidateRect(r image.Rectangle) {
	var rect = RECT{uint32(r.Min.X), uint32(r.Min.Y), uint32(r.Max.X), uint32(r.Max.Y)}
	InvalidateRect(h.hwnd, &rect, FALSE)
}
