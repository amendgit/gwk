// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vango

import (
	. "gwk/sysc"
	"image"
	"log"
	"unsafe"
)

type NativeContext Handle

type NativeCanvas struct {
	*Canvas
	hdc     Handle // lazily-created.
	hbitmap Handle // lazily-created.
}

func NewNativeCanvas(bounds image.Rectangle) *NativeCanvas {
	var c NativeCanvas

	// We can't use NewCanvas here, because it will allocate memory.
	c.Canvas = new(Canvas)

	c.SetBounds(bounds)

	var bmi BITMAPINFOHEADER

	var w, h = c.W(), c.H()
	if w == 0 || h == 0 {
		w, h = 1, 1
	}

	bmi.Size = uint32(unsafe.Sizeof(bmi))
	bmi.Width = int32(w)
	bmi.Height = int32(-h)
	bmi.Planes = 1
	bmi.BitCount = 32
	bmi.Compression = BI_RGB
	bmi.SizeImage = 0
	bmi.XPelsPerMeter = 1
	bmi.YPelsPerMeter = 1
	bmi.ClrUsed = 0
	bmi.ClrImportant = 0

	var p unsafe.Pointer
	var err error
	c.hbitmap, err = CreateDIBSection(NULL, &bmi, 0, &p, NULL, 0)

	if err != nil {
		log.Printf("CreateDIBSection %v", err)
	}

	c.SetPix(((*[1 << 30]byte)(unsafe.Pointer(p)))[:w*h*4])
	c.SetBounds(image.Rect(0, 0, w, h))
	c.SetStride(w * 4)

	return &c
}

func (c *NativeCanvas) Opaque() bool {
	return true
}

func (c *NativeCanvas) Release() bool {
	DeleteDC(c.hdc)
	DeleteObject(c.hbitmap)
	c.hdc = NULL
	return true
}

func (c *NativeCanvas) BeginPaint() NativeContext {
	if c.hdc == 0 {
		var err error
		if c.hdc, err = CreateCompatibleDC(NULL); err != nil {
			log.Printf("CreateCompatibleDC(NULL) %v", err)
		}
		var oldBitmap = SelectObject(c.hdc, c.hbitmap)
		// When the memory DC is created, its display surface is exactly one
		// monochrome pixel wide and one monochrome pixel high. Since we select
		// our own bitmap, we must delete the previous one.
		DeleteObject(oldBitmap)
	}
	return NativeContext(c.hdc)
}

func (c *NativeCanvas) EndPaint() {
	return
}

func (c *NativeCanvas) BlitToContext(nc NativeContext, x int, y int, srcRc *image.Rectangle) {
	// log.Printf("NativeCanvas.BlitToContext(...)")
	var srcDC = c.BeginPaint()

	if srcRc == nil {
		var tempRc = image.Rect(0, 0, c.W(), c.H())
		srcRc = &tempRc
	}

	var copyWidth = srcRc.Dx()
	var copyHeight = srcRc.Dy()

	if c.Opaque() {
		// log.Printf("BitBlt(%v %v %v)", x, y, *srcRc)

		var err = BitBlt(Handle(nc),
			int32(x), int32(y), int32(copyWidth), int32(copyHeight),
			Handle(srcDC),
			int32(srcRc.Min.X), int32(srcRc.Min.Y),
			SRCCOPY)

		if err != nil {
			log.Printf("BitBlt(...) %v", err)
		}
	} else {
		var bf = BLENDFUNCTION{AC_SRC_OVER, 0, 255, AC_SRC_ALPHA}
		GdiAlphaBlend(Handle(nc),
			int32(x), int32(y), int32(copyWidth), int32(copyHeight),
			Handle(srcDC),
			int32(srcRc.Min.X), int32(srcRc.Min.Y), int32(copyWidth), int32(copyHeight),
			bf)
	}

	c.EndPaint()
}
