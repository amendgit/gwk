// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vango

import (
	. "image"
	"log"
)

type Canvas struct {
	pix    []byte    // Pixels in RGBA order.
	bounds Rectangle // Bounds is the sub rectangle of the pixels's bounds.
	stride int       // The number of pixels in one line.
	opaque bool      // Does the canvas opaque.
}

func NewCanvas(width int, height int) *Canvas {
	log.Printf("NewCanvas")
	var c Canvas
	c.bounds = Rect(0, 0, width, height)
	c.pix = make([]byte, c.W()*c.H()*4)
	c.stride = c.W()
	return &c
}

func (c *Canvas) X() int {
	return c.bounds.Min.X
}

func (c *Canvas) Y() int {
	return c.bounds.Min.Y
}

func (c *Canvas) W() int {
	return c.bounds.Dx()
}

func (c *Canvas) H() int {
	return c.bounds.Dy()
}

func (c *Canvas) Stride() int {
	return c.stride
}

func (c *Canvas) SetStride(stride int) {
	c.stride = stride
}

func (c *Canvas) Pix() []byte {
	return c.pix
}

func (c *Canvas) SetPix(pix []byte) {
	c.pix = pix
}

func (c *Canvas) Bounds() Rectangle {
	return c.bounds
}

func (c *Canvas) LocalBounds() Rectangle {
	return c.bounds.Sub(c.bounds.Min)
}

func (c *Canvas) SetBounds(bounds Rectangle) {
	c.bounds = bounds
}

func (c *Canvas) Opaque() bool {
	return c.opaque
}

func (c *Canvas) SetOpaque(opaque bool) {
	c.opaque = opaque
}

func (c *Canvas) SubCanvas(r Rectangle) *Canvas {
	// The SubImage in the image pkg is need the r based on the absolute
	// coordinate. We need r based on the relative coordinate. So covnert
	// r to the parent's coordinate first.
	r = r.Add(c.Bounds().Min)
	r = r.Intersect(c.Bounds())
	if r.Empty() {
		return &Canvas{}
	}
	// i := c.PixOffset(r.Min.X, r.Min.Y)
	return &Canvas{
		pix:    c.pix,
		stride: c.stride,
		bounds: r,
	}
}

func (c *Canvas) PixOffset(x int, y int) int {
	return (x+c.bounds.Min.X)*4 + (y+c.bounds.Min.Y)*c.Stride()*4
}

func (dst *Canvas) DrawColor(r, g, b byte) {
	var i = dst.PixOffset(0, 0)
	var W = i + dst.W()*4
	var p = dst.Pix()
	var j = 0

	for j < dst.H() {
		for i < W {
			p[i+0] = b
			p[i+1] = g
			p[i+2] = r
			p[i+3] = 255
			i += 4
		}
		i = i + dst.Stride()*4 - dst.W()*4
		W = W + dst.Stride()*4
		j++
	}
}

func (c *Canvas) FillRect(rect Rectangle, r, g, b byte) {
	canvas_rect := c.Bounds()
	rect = rect.Add(canvas_rect.Min).Intersect(canvas_rect).Sub(canvas_rect.Min)
	line_begin_index := c.PixOffset(rect.Min.X, rect.Min.Y)
	line_end_index := line_begin_index + rect.Dx()*4
	stride := c.Stride() * 4
	pix := c.Pix()

	for j := 0; j < rect.Dy(); j++ {
		for index := line_begin_index; index < line_end_index; index = index + 4 {
			pix[index+0] = b
			pix[index+1] = g
			pix[index+2] = r
			pix[index+3] = 0
		}
		line_begin_index = line_begin_index + stride
		line_end_index = line_end_index + stride
	}
}

func (dst *Canvas) DrawLine(from Point, to Point) {
	return
}

func (dst *Canvas) DrawCanvas(x int, y int, src *Canvas) {
	// 0 means src, 1 means dst.
	// b0, b1 := src.Bounds(), dst.Bounds()
	l0, l1 := src.LocalBounds(), dst.LocalBounds()
	x0, y0, x1, y1 := 0, 0, x, y
	i0, i1 := src.PixOffset(x0, y0), dst.PixOffset(x1, y1)
	s0, s1 := src.Stride(), dst.Stride()
	p0, p1 := src.Pix(), dst.Pix()

	// the shared draw rect.
	r := l0.Intersect(l1)
	if r.Empty() {
		return
	}

	w, h := r.Dx(), r.Dy()

	// from src(x0, y0) draw |r| area to dst(x1, y1)
	for j := 0; j < h; j++ {
		for i := 0; i < w*4; i = i + 4 {
			p1[i1+i+0] = p0[i0+i+0]
			p1[i1+i+1] = p0[i0+i+1]
			p1[i1+i+2] = p0[i0+i+2]
			p1[i1+i+3] = p0[i0+i+3]
		}
		i0 = i0 + s0*4
		i1 = i1 + s1*4
	}
}

func (dst *Canvas) AlphaBlend(x int, y int, src *Canvas) {
	// 0 means src, 1 means dst.
	l0, l1 := src.LocalBounds(), dst.LocalBounds()
	x0, y0, x1, y1 := 0, 0, x, y
	i0, i1 := src.PixOffset(x0, y0), dst.PixOffset(x1, y1)
	s0, s1 := src.Stride(), dst.Stride()
	p0, p1 := src.Pix(), dst.Pix()

	// the shared draw rect.
	r := l0.Intersect(l1)
	if r.Empty() {
		return
	}

	w, h := r.Dx(), r.Dy()

	// from src(x0, y0) draw |r| area to dst(x1, y1)
	for j := 0; j < h; j++ {
		for i := 0; i < w*4; i = i + 4 {
			// http://archive.gamedev.net/archive/reference/articles/article817.html
			r0, g0, b0 := p0[i0+i+0], p0[i0+i+1], p0[i0+i+2]
			r1, g1, b1 := p1[i1+i+0], p1[i1+i+1], p1[i1+i+2]

			a := int32(p0[i0+i+3])

			p1[i1+i+0] = byte((a*(int32(r0)-int32(r1)))/256) + r1
			p1[i1+i+1] = byte((a*(int32(g0)-int32(g1)))/256) + g1
			p1[i1+i+2] = byte((a*(int32(b0)-int32(b1)))/256) + b1
			// p1[i1+i+3] = p1[i1+i+3]
		}
		i0 = i0 + s0*4
		i1 = i1 + s1*4
	}
}

func (dst *Canvas) DrawImageNRGBA(x int, y int, src *NRGBA, srcRc *Rectangle) {
	if srcRc == nil {
		srcRc = &(src.Rect)
	}

	var srcX, srcY = srcRc.Min.X, srcRc.Min.Y
	var dstX, dstY = x, y

	var bltW, bltH = srcRc.Dx(), srcRc.Dy()

	var srcI = src.PixOffset(srcX, srcY)
	var dstI = dst.PixOffset(dstX, dstY)

	var srcPix = src.Pix
	var dstPix = dst.Pix()

	var srcStride = src.Stride
	var dstStride = dst.Stride() * 4

	var i, j = 0, 0

	for j < bltH {
		i = 0
		for i < bltW*4 {
			dstPix[dstI+i+0] = srcPix[srcI+i+0]
			dstPix[dstI+i+1] = srcPix[srcI+i+1]
			dstPix[dstI+i+2] = srcPix[srcI+i+2]
			dstPix[dstI+i+3] = srcPix[srcI+i+3]
			i += 4
		}
		srcI = srcI + srcStride
		dstI = dstI + dstStride
		j++
	}
}

func (dst *Canvas) DrawImageRGBA(x int, y int, src *RGBA, srcRc *Rectangle) {
	if srcRc == nil {
		srcRc = &(src.Rect)
	}

	var srcX, srcY = srcRc.Min.X, srcRc.Min.Y
	var dstX, dstY = x, y

	var bltW, bltH = srcRc.Dx(), srcRc.Dy()

	var srcI = src.PixOffset(srcX, srcY)
	var dstI = dst.PixOffset(dstX, dstY)

	var srcPix = src.Pix
	var dstPix = dst.Pix()

	var srcStride = src.Stride
	var dstStride = dst.Stride() * 4

	var i, j = 0, 0

	for j < bltH {
		i = 0
		for i < bltW*4 {
			dstPix[dstI+i+0] = srcPix[srcI+i+2]
			dstPix[dstI+i+1] = srcPix[srcI+i+1]
			dstPix[dstI+i+2] = srcPix[srcI+i+0]
			dstPix[dstI+i+3] = srcPix[srcI+i+3]
			i += 4
		}
		srcI = srcI + srcStride
		dstI = dstI + dstStride
		j++
	}
}

func (dst *Canvas) DrawTexture(dstRc Rectangle, tex *Canvas, texRc Rectangle) {
	var texX, texY = texRc.Min.X, texRc.Min.Y
	var dstX, dstY = dstRc.Min.X, texRc.Min.Y

	var texW, texH = texRc.Dx(), texRc.Dy()
	var dstW, dstH = dstRc.Dx(), dstRc.Dy()

	var texI = tex.PixOffset(texX, texY)
	var dstI = dst.PixOffset(dstX, dstY)

	var texStride = tex.Stride() * 4
	var dstStride = dst.Stride() * 4

	var texPix = tex.Pix()
	var dstPix = dst.Pix()

	var i0, i1 = 0, 0
	var j0, j1 = 0, 0

	for {
		dstPix[dstI+i1+0] = texPix[texI+i0+0]
		dstPix[dstI+i1+1] = texPix[texI+i0+1]
		dstPix[dstI+i1+2] = texPix[texI+i0+2]
		dstPix[dstI+i1+3] = texPix[texI+i0+3]

		i0, i1 = i0+4, i1+4

		if i0 >= texW*4 {
			i0 = 0
		}

		if i1 >= dstW*4 {
			i0, i1 = 0, 0

			texI = texI + texStride
			dstI = dstI + dstStride

			j0, j1 = j0+1, j1+1

			if j0 >= texH {
				j0 = 0
				texI = tex.PixOffset(texX, texY)
			}

			if j1 >= dstH {
				break
			}
		}
	}
}

func (dst *Canvas) StretchDraw(dst_rect Rectangle, src *Canvas) {
	bounds0 := src.Bounds()
	bounds1 := dst.Bounds()

	rect0 := bounds0
	rect1 := dst_rect.Add(bounds1.Min).Intersect(bounds1)

	width0, height0 := rect0.Dx(), rect0.Dy()
	width1, height1 := rect1.Dx(), rect1.Dy()

	stride0 := src.Stride() * 4
	log.Printf("jishi debug: src stride %v", stride0)
	pix0 := src.Pix()
	pix1 := dst.Pix()

	max_pix_offset0 := len(pix0) - 1

	scale_x := float64(width0) / float64(width1)
	scale_y := float64(height0) / float64(height1)

	to_color_channel := func(f64 float64) byte {
		if f64 < 255 {
			return byte(f64)
		}
		return 255
	}

	for j := 0; j < height1; j++ {
		for i := 0; i < width1; i++ {
			xf := float64(i) * scale_x
			yf := float64(j) * scale_y

			x0, y0 := int(xf), int(yf)
			pix_offset0 := src.PixOffset(x0, y0)

			pix_offset1 := pix_offset0 + 4
			pix_offset2 := pix_offset0 + stride0
			pix_offset3 := pix_offset2 + 4

			if pix_offset3 > max_pix_offset0 {
				break
			}

			b0, g0, r0 := pix0[pix_offset0+0], pix0[pix_offset0+1], pix0[pix_offset0+2]
			b1, g1, r1 := pix0[pix_offset1+0], pix0[pix_offset1+1], pix0[pix_offset1+2]
			b2, g2, r2 := pix0[pix_offset2+0], pix0[pix_offset2+1], pix0[pix_offset2+2]
			b3, g3, r3 := pix0[pix_offset3+0], pix0[pix_offset3+1], pix0[pix_offset3+2]

			factor0 := xf - float64(int(xf))
			factor1 := 1 - factor0

			b4 := factor1*float64(b0) + factor0*float64(b1)
			g4 := factor1*float64(g0) + factor0*float64(g1)
			r4 := factor1*float64(r0) + factor0*float64(r1)
			b5 := factor1*float64(b2) + factor0*float64(b3)
			g5 := factor1*float64(g2) + factor0*float64(g3)
			r5 := factor1*float64(r2) + factor0*float64(r3)

			factor3 := yf - float64(int(yf))
			factor4 := 1 - factor3

			b := factor4*b4 + factor3*b5
			g := factor4*g4 + factor3*g5
			r := factor4*r4 + factor3*r5

			pix_offset := dst.PixOffset(i, j)
			pix1[pix_offset+0] = to_color_channel(b)
			pix1[pix_offset+1] = to_color_channel(g)
			pix1[pix_offset+2] = to_color_channel(r)
		}
	}
}
