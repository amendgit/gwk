package vango

import (
	"errors"
	"gwk/vango/freetype"
	"image"
	"log"
)

type Context struct {
	font         *Font
	canvas       *Canvas
	dpi          float64
	stroke_color uint32
	fill_color   uint32
	font_color   uint32
}

func NewContext() *Context {
	ctxt := new(Context)
	ctxt.font = NewFont()
	ctxt.dpi = 72
	ctxt.stroke_color = 0x000000
	return ctxt
}

func (c *Context) SetStrokeColor(r, g, b byte) {
	c.stroke_color = uint32(b)<<8 | uint32(g)<<16 | uint32(r)<<24
}

func (c *Context) SetFillColor(r, g, b byte) {
	c.fill_color = uint32(b)<<8 | uint32(g)<<16 | uint32(r)<<24
}

func (c *Context) SetFontColor(r, g, b byte) {
	c.font_color = uint32(b)<<8 | uint32(g)<<16 | uint32(r)<<24
}

func (c *Context) SetFont(font_name string) {
	if font_name == "default" {
		c.font.font = g_default_font
	} else {
		log.Printf("NOT IMPLEMENTATION: font name %v", font_name)
	}
}

func (c *Context) SetCanvas(canvas *Canvas) *Canvas {
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
			pt.X += freetype.Fix32(c.font.Kerning(prev, idx)) << 2
		}

		mask, offset, err := c.font.GlyphAt(idx, pt)
		if err != nil {
			return freetype.RastPoint{}, err
		}

		pt.X += freetype.Fix32(c.font.HMetric(idx).AdvanceWidth) << 2
		glyph_rect := mask.Bounds().Add(offset)

		c.draw_text_mask(glyph_rect.Min.X, glyph_rect.Min.Y, mask)

		prev, has_prev = idx, true
	}
	return pt, nil
}

func (c *Context) draw_text_mask(x, y int, mask *image.Alpha) {
	src := mask
	dst := c.canvas
	rect := mask.Bounds()

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

	clr := c.font_color
	b, g, r := byte(clr>>8&0xff), byte(clr>>16&0xff), byte(clr>>24&0xff)

	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			o0, o1 := i0+int(x), i1+x*4 // pix offset in bytes

			a := int32(p0[o0])
			if a == 0 {
				continue
			}

			r1, g1, b1 := p1[o1+0], p1[o1+1], p1[o1+2]

			p1[o1+0] = byte((a*(int32(r)-int32(r1)))/256) + r1
			p1[o1+1] = byte((a*(int32(g)-int32(g1)))/256) + g1
			p1[o1+2] = byte((a*(int32(b)-int32(b1)))/256) + b1
		}
		i0 = i0 + s0
		i1 = i1 + s1
	}
}

func (c *Context) DrawColor(r, g, b byte) {
	dst := c.canvas
	i := dst.PixOffset(0, 0)
	dr := dst.LocalBounds()
	p := dst.Pix()

	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx()*4; x += 4 {
			p[i+x+0] = b
			p[i+x+1] = g
			p[i+x+2] = r
			p[i+x+3] = 255
		}
		i += dst.Stride()
	}
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

func (c *Context) DrawStretch(dst_rect image.Rectangle, src *Canvas, src_rect image.Rectangle) {

}

func (c *Context) FillRect(rect image.Rectangle) {
	dst := c.canvas
	l := dst.LocalBounds()
	dr := rect.Intersect(l) // draw rect

	if dr.Empty() {
		return
	}
	s := dst.Stride()
	p := dst.Pix()

	i0 := dst.PixOffset(dr.Min.X, dr.Min.Y) // offset of the pix that at the BEGIN of one line
	i1 := i0 + dr.Dx()*4                    // offset of the pix that at the END of one line

	clr := c.fill_color
	b, g, r := byte(clr>>8&0xff), byte(clr>>16&0xff), byte(clr>>24&0xff)
	for y := 0; y < dr.Dy(); y++ {
		for x := i0; x < i1; x = x + 4 {
			p[x+0] = b
			p[x+1] = g
			p[x+2] = r
			p[x+3] = 0
		}
		i0 += s
		i1 += s
	}
}

func (c *Context) StrokeRect(rect image.Rectangle) {
	dst := c.canvas
	i := dst.PixOffset(rect.Min.X, rect.Min.Y)
	s := dst.Stride()
	p := dst.Pix()

	clr := c.stroke_color
	b, g, r := byte(clr>>8&0xff), byte(clr>>16&0xff), byte(clr>>24&0xff)
	// r, g, b = 0x00, 0xff, 0xff
	// step 1
	for x := 0; x < rect.Dx(); x++ {
		o := i + x*4
		p[o+0] = r
		p[o+1] = g
		p[o+2] = b
		p[o+3] = 0xff
	}

	// step 2
	i = dst.PixOffset(rect.Min.X, rect.Max.Y)
	for x := 0; x < rect.Dx(); x++ {
		o := i + x*4
		p[o+0] = r
		p[o+1] = g
		p[o+2] = b
		p[o+3] = 0xff
	}

	// step 3
	i = dst.PixOffset(rect.Min.X, rect.Min.Y)
	for y := 0; y < rect.Dy(); y++ {
		o := i + y*s
		p[o+0] = r
		p[o+1] = g
		p[o+2] = b
		p[o+3] = 0xff
	}

	// step 4
	i = dst.PixOffset(rect.Max.X, rect.Min.Y)
	for y := 0; y < rect.Dy(); y++ {
		o := i + y*s
		p[o+0] = r
		p[o+1] = g
		p[o+2] = b
		p[o+3] = 0xff
	}
}
