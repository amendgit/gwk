package vango

import (
	"errors"
	"gwk/vango/freetype"
	"image"
	"log"
)

type Context struct {
	font   *Font
	canvas *Canvas
	dpi    float64
}

func NewContext() *Context {
	ctxt := new(Context)
	ctxt.font = NewFont()
	ctxt.dpi = 72
	return ctxt
}

func (c *Context) SelectFont(font_name string) {
	if font_name == "default" {
		c.font.font = g_default_font
	} else {
		log.Printf("NOT IMPLEMENTATION: font name %v", font_name)
	}
}

func (c *Context) SelectCanvas(canvas *Canvas) *Canvas {
	old := c.canvas
	c.canvas = canvas
	return old
}

func (c *Context) DrawText(text string, rect image.Rectangle) (freetype.RastPoint, error) {
	if c.font == nil {
		return freetype.RastPoint{}, errors.New("vango DrawString called with nil font.")
	}

	pt := freetype.Point(rect.Min.X+1, rect.Min.Y+10)

	prev, has_prev := uint16(0), false
	for _, rune := range text {
		idx := c.font.Index(rune)
		if has_prev {
			pt.X += freetype.Fix32(c.font.Kerning(768, prev, idx)) << 2
		}

		mask, offset, err := c.font.GlyphAt(idx, pt)
		if err != nil {
			return freetype.RastPoint{}, err
		}

		pt.X += freetype.Fix32(c.font.font.HMetric(768, idx).AdvanceWidth) << 2
		glyph_rect := mask.Bounds().Add(offset)

		c.DrawImage(glyph_rect.Min.X, glyph_rect.Min.Y, mask, mask.Bounds())

		prev, has_prev = idx, true
	}
	return pt, nil
}

func (c *Context) DrawImage(x, y int, src image.Image, rect image.Rectangle) {
	switch typ := src.(type) {
	case *image.Alpha:
		c.DrawAlpha(x, y, typ, rect)
	}
}

func (c *Context) DrawAlpha(x, y int, src *image.Alpha, rect image.Rectangle) {
	dst := c.canvas

	i0, i1 := src.PixOffset(rect.Min.X, rect.Min.Y), dst.PixOffset(x, y) // pix offset
	s0, s1 := src.Stride, dst.Stride()                                   // stride
	p0, p1 := src.Pix, dst.Pix()                                         // pix

	// calculate the draw rect
	r0 := rect.Sub(rect.Min) // align in (0, 0)
	r1 := dst.LocalBounds()  // bounds in self coordinate.
	r1.Min = image.Pt(x, y)  // dst draw position start at (x, y)
	r1 = r1.Sub(r1.Min)      // align in (0, 0)

	dr := r0.Intersect(r1)
	if dr.Empty() {
		return
	}

	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			o0, o1 := i0+int(x), i1+x*4 // pix offset in bytes

			a := int32(p0[o0])
			if a == 0 {
				continue
			}

			r0, g0, b0 := 0x00, 0x00, 0x00
			r1, g1, b1 := p1[o1+0], p1[o1+1], p1[o1+2]

			p1[o1+0] = byte((a*(int32(r0)-int32(r1)))/256) + r1
			p1[o1+1] = byte((a*(int32(g0)-int32(g1)))/256) + g1
			p1[o1+2] = byte((a*(int32(b0)-int32(b1)))/256) + b1
		}
		i0 = i0 + s0
		i1 = i1 + s1
	}
}

func (c *Context) DrawNRGBA(x int, y int, src *image.NRGBA, rect image.Rectangle) {
	dst := c.canvas
	x0, y0 := rect.Min.X, rect.Min.Y
	x1, y1 := x, y

	i0, i1 := src.PixOffset(x0, y0), dst.PixOffset(x1, y1) // index
	p0, p1 := src.Pix, dst.Pix()                           // pix
	s0, s1 := src.Stride, dst.Stride()                     // stride

	// draw rect
	r0 := rect
	r1 := dst.LocalBounds()
	r1.Min = image.Pt(x, y)
	r1 = r1.Sub(r1.Min)
	dr := r0.Intersect(r1)
	if dr.Empty() {
		return
	}

	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			o0, o1 := i0+4*x, i1+4*x
			p1[o1+0] = p0[o0+0]
			p1[o1+1] = p0[o0+1]
			p1[o1+2] = p0[o0+2]
			p1[o1+3] = p0[o0+3]
		}
		i0 = i0 + s0
		i1 = i1 + s1
	}
}

// TODO(refactory)
func (c *Context) DrawRGBA(x, y int, src *image.RGBA, rect image.Rectangle) {
	dst := c.canvas

	x0, y0 := rect.Min.X, rect.Min.Y
	x1, y1 := x, y

	i0, i1 := src.PixOffset(x0, y0), dst.PixOffset(x1, y1)
	p0, p1 := src.Pix, dst.Pix()
	s0, s1 := src.Stride, dst.Stride()

	// calculate draw rect.
	r0 := rect
	r1 := dst.LocalBounds()
	r1.Min = image.Pt(x, y)
	r1 = r1.Sub(r1.Min)
	dr := r0.Intersect(r1)

	if dr.Empty() {
		return
	}

	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			o0, o1 := i0+4*x, i1+4*x
			p1[o1+0] = p0[o0+0]
			p1[o1+1] = p0[o0+1]
			p1[o1+2] = p0[o0+2]
			p1[o1+3] = p0[o0+3]
		}
		i0 = i0 + s0
		i1 = i1 + s1
	}
}

func (c *Context) DrawCanvas(x, y int, src *Canvas, rect image.Rectangle) {
	dst := c.canvas
	// 0 means src, 1 means dst.
	// b0, b1 := src.Bounds(), dst.Bounds()
	// l0, l1 := src.LocalBounds(), dst.LocalBounds()
	x0, y0, x1, y1 := rect.Min.X, rect.Min.Y, x, y
	i0, i1 := src.PixOffset(x0, y0), dst.PixOffset(x1, y1)
	s0, s1 := src.Stride(), dst.Stride()
	p0, p1 := src.Pix(), dst.Pix()

	// the shared draw rect.
	r0 := rect.Sub(rect.Min)
	r1 := dst.LocalBounds()
	r1.Min = image.Pt(x, y)

	dr := r0.Intersect(r1)
	if dr.Empty() {
		return
	}

	// from src(x0, y0) draw |r| area to dst(x1, y1)
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			o0, o1 := i0+4*x, i1+4*x
			p1[o1+0] = p0[o0+0]
			p1[o1+1] = p0[o0+1]
			p1[o1+2] = p0[o0+2]
			p1[o1+3] = p0[o0+3]
		}
		i0 = i0 + s0
		i1 = i1 + s1
	}
}

func (c *Context) AlphaBlend(x int, y int, src *Canvas, rect image.Rectangle) {
	dst := c.canvas
	// 0 means src, 1 means dst.
	// l0, l1 := src.LocalBounds(), dst.LocalBounds()
	x0, y0, x1, y1 := 0, 0, x, y
	i0, i1 := src.PixOffset(x0, y0), dst.PixOffset(x1, y1)
	s0, s1 := src.Stride(), dst.Stride()
	p0, p1 := src.Pix(), dst.Pix()

	// the shared draw rect.
	r0 := rect.Sub(rect.Min)
	r1 := dst.LocalBounds()
	r1.Min = image.Pt(x, y)
	r1 = r1.Sub(r1.Min)
	dr := r0.Intersect(r1)
	if dr.Empty() {
		return
	}

	// from src(x0, y0) draw |r| area to dst(x1, y1)
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			// http://archive.gamedev.net/archive/reference/articles/article817.html
			o0, o1 := i0+4*x, i1+4*x
			r0, g0, b0 := p0[o0+0], p0[o0+1], p0[o0+2]
			r1, g1, b1 := p1[o1+0], p1[o1+1], p1[o1+2]

			// alpha value
			a := int32(p0[o0+3])

			p1[o1+0] = byte((a*(int32(r0)-int32(r1)))/256) + r1
			p1[o1+1] = byte((a*(int32(g0)-int32(g1)))/256) + g1
			p1[o1+2] = byte((a*(int32(b0)-int32(b1)))/256) + b1
		}

		i0 = i0 + s0
		i1 = i1 + s1
	}
}
