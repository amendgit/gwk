// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sysc

import (
	"unsafe"
)

type Handle uintptr

const NULL = 0

func MAKEINTRESOURCE(n int32) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(n)))
}

//sys	CreateWindowEx(exStyle uint32, className *uint16, windowName *uint16, style uint32, x int, y int, width int, height int, wndParent Handle, menu Handle, instance Handle, param uintptr) (hwnd Handle, err error) = user32.CreateWindowExW
//sys	RegisterClassEx(wcx *WNDCLASSEX) (atom uint16, err error) = user32.RegisterClassExW
//sys	ShowWindow(hwnd Handle, cmdShow int32) (visiable int) = user32.ShowWindow
//sys	UpdateWindow(hwnd Handle) (isUpdated int) = user32.UpdateWindow
//sys	LoadIcon(instance Handle, iconName *uint16) (hicon Handle, err error) = user32.LoadIconW
//sys	LoadCursor(instance Handle, cursorName *uint16) (hcursor Handle, err error) = user32.LoadCursorW
//sys	GetMessage(msg *MSG, hwnd Handle, msgFilterMin uint32, msgFilterMax uint32) (EOL int) = user32.GetMessageW
//sys	PeekMessage(msg *MSG, hwnd Handle, msgFilterMin uint32, msgFilterMax uint32, removeMsg uint32) (has_msg bool) = user32.PeekMessageW
//sys	TranslateMessage(msg *MSG) (isTranslated int) = user32.TranslateMessage
//sys	DispatchMessage(msg *MSG) (ingore int) = user32.DispatchMessageW
//sys	GetClientRect(hwnd Handle, rect *RECT) (err error) = user32.GetClientRect
//sys	GetWindowRect(hwnd Handle, rect *RECT) (err error) = user32.GetWindowRect
//sys	GetDC(hwnd Handle) (hDC Handle, err error) = user32.GetDC
//sys	BeginPaint(hwnd Handle, ps *PAINTSTRUCT) (hDC Handle) = user32.BeginPaint
//sys	EndPaint(hwnd Handle, ps *PAINTSTRUCT) (err error) = user32.EndPaint
//sys	DrawTextEx(hDC Handle, text *uint16, length int32, rc *RECT, format uint32, params *DRAWTEXTPARAMS) (retval int32) = user32.DrawTextExW
//sys	PostQuitMessage(exitcode int32) = user32.PostQuitMessage
//sys	DefWindowProc(hwnd Handle, msg uint32, warg uintptr, larg uintptr) (retval uintptr) = user32.DefWindowProcW
//sys	GetUpdateRect(hwnd Handle, rect *RECT, erase int32) (isempty int32) = user32.GetUpdateRect
//sys	RedrawWindow(hwnd Handle, rect *RECT, hrgn Handle, flags uint32) = user32.RedrawWindow
//sys	InvalidateRect(hwnd Handle, rect *RECT, isErased int) (ok int) = user32.InvalidateRect

//sys	BitBlt(hDC Handle, xDext int32, yDext int32, width int32, height int32, hDCSrc Handle, xSrc int32, ySrc int32, rop uint32) (err error) = gdi32.BitBlt
//sys	SetDIBitsToDevice(hDC Handle, xDext int32, yDest int32, width int32, height int32, xSrc int32, ySrc int32, startScan uint32, scanLines uint32, bits uintptr, bmi *BITMAPINFO, colorUse uint32) (lines int32) = gdi32.SetDIBitsToDevice
//sys	DeleteDC(hDC Handle) (err error) = gdi32.DeleteDC
//sys	CreateCompatibleDC(hdc Handle) (cmptDC Handle, err error) = gdi32.CreateCompatibleDC
//sys	SelectObject(hDC Handle, hgdiobj Handle) (oldObj Handle) = gdi32.SelectObject
//sys	DeleteObject(hObject Handle) (isDeleted bool) = gdi32.DeleteObject
//sys	gdiAlphaBlend(hdcDest Handle, xoriginDest int32, yoriginDest int32, wDest int32, hDest int32, hdcSrc Handle, xoriginSrc int32, yoriginSrc int32, wSrc int32, hSrc int32, ftn uint32) (err error) = gdi32.GdiAplhaBlend
//sys	CreateDIBSection(hdc Handle, bmi *BITMAPINFOHEADER, iUsage uint, ppvBits *unsafe.Pointer, hSection Handle, offset uint32) (hbitmp Handle, err error) = gdi32.CreateDIBSection

func GdiAlphaBlend(hdcDest Handle, xoriginDest int32, yoriginDest int32, wDest int32, hDest int32, hdcSrc Handle, xoriginSrc int32, yoriginSrc int32, wSrc int32, hSrc int32, bf BLENDFUNCTION) error {
	var ftn uint32
	ftn = ftn & uint32(bf.BlendOp)
	ftn = (ftn << 8) & uint32(bf.BlendFlags)
	ftn = (ftn << 8) & uint32(bf.SourceConstantAlpha)
	ftn = (ftn << 8) & uint32(bf.AlphaFormat)
	return gdiAlphaBlend(hdcDest, xoriginDest, yoriginDest, wDest, hDest, hdcSrc, xoriginSrc, yoriginSrc, wSrc, hSrc, ftn)
}
