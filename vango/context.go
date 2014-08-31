package vango

import (
	"bufio"
	"errors"
	"gwk/vango/freetype"
	"image"
	"image/png"
	"log"
	"os"
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

	pt := freetype.Point(rect.Min.X, rect.Min.Y)

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
		fd, _ := os.Create("a.png")
		defer fd.Close()
		bio := bufio.NewWriter(fd)
		png.Encode(bio, mask)
		// dr := c.clip.Intersect(glyph_rect)
		//if !dr.Empty() {
		// mp := image.Point{0, dr.Min.Y - glyph_rect.Min.Y}
		// draw.DrawMask(c.dst, dr, c.src, image.ZP, mask, mp, draw.Over)
		c.DrawImage(rect.Min.X, rect.Min.Y, mask, glyph_rect)
		//}
		prev, has_prev = idx, true
	}
	return pt, nil
}

func (c *Context) DrawImage(x, y int, img image.Image, rect image.Rectangle) {
	switch typ := img.(type) {
	case *image.Alpha:
		c.DrawAlpha(x, y, typ, rect)
	}
}

func (c *Context) DrawAlpha(x, y int, src *image.Alpha, rect image.Rectangle) {
	rect = rect.Sub(rect.Min)
	dst := c.canvas

	i0, i1 := src.PixOffset(rect.Min.X, rect.Min.Y), dst.PixOffset(x, y) // pix offset
	s0, s1 := src.Stride, dst.Stride()                                   // stride
	p0, p1 := src.Pix, dst.Pix()                                         // pix
	// l0, l1 := src.Rect.Sub(src.Rect.Min), dst.LocalBounds()           // local bounds
	// b0, b1 := src.Bounds(), dst.Bounds()                              // bounds

	// draw rect
	r0 := rect.Sub(rect.Min)
	r1 := dst.Bounds()
	r1.Min = image.Pt(x, y)

	r := r0.Intersect(r1)
	if r.Empty() {
		return
	}

	log.Printf("r %v r0 %v r1 %v rect %v", r, r0, r1, rect)

	log.Printf("src %v pix %v", src.Bounds(), len(src.Pix))

	for y := 0; y < r.Dy(); y++ {
		for x := 0; x < r.Dx(); x += 4 {
			if p0[i0+x+3] == 0 {
				continue
			}
			p1[i1+x+0] = p0[i0+x+0]
			p1[i1+x+1] = p0[i0+x+1]
			p1[i1+x+2] = p0[i0+x+2]
			p1[i1+x+3] = p0[i0+x+3]
			// p1[i1+x+0] = 0x00
			// p1[i1+x+1] = 0x00
			// p1[i1+x+2] = 0x00
			// p1[i1+x+3] = 0x00

		}
		i0 = i0 + s0
		i1 = i1 + 4*s1
	}
}

func (c *Context) DrawImageNRGBA(x int, y int, src *image.NRGBA, rect image.Rectangle) {
	dst := c.canvas
	x0, y0 := rect.Min.X, rect.Min.Y
	x1, y1 := x, y

	i0, i1 := src.PixOffset(x0, y0), dst.PixOffset(x1, y1) // index
	p0, p1 := src.Pix, dst.Pix()                           // pix
	s0, s1 := src.Stride, dst.Stride()*4                   // stride
	// draw rect
	r0 := rect
	r1 := dst.LocalBounds()
	r1.Min = image.Pt(x, y)
	r := r0.Intersect(r1)
	if r.Empty() {
		return
	}

	for y := 0; y < r.Dy(); y++ {
		for x := 0; x < r.Dx(); x += 4 {
			p1[i1+x+0] = p0[i0+x+0]
			p1[i1+x+1] = p0[i0+x+1]
			p1[i1+x+2] = p0[i0+x+2]
			p1[i1+x+3] = p0[i0+x+3]
			x += 4
		}
		i0 = i0 + s0
		i1 = i1 + s1
	}
}

func (c *Context) DrawRGBA(x int, y int, src *image.RGBA, srcRc *image.Rectangle) {
	dst := c.canvas
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

func (c *Context) DrawCanvas(x int, y int, src *Canvas, rect image.Rectangle) {
	dst := c.canvas
	// 0 means src, 1 means dst.
	// b0, b1 := src.Bounds(), dst.Bounds()
	l0, l1 := src.LocalBounds(), dst.LocalBounds()
	x0, y0, x1, y1 := rect.Min.X, rect.Min.Y, x, y
	i0, i1 := src.PixOffset(x0, y0), dst.PixOffset(x1, y1)
	s0, s1 := src.Stride(), dst.Stride()
	p0, p1 := src.Pix(), dst.Pix()

	// TODO(BUG)
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

func (c *Context) AlphaBlend(x int, y int, src *Canvas) {
	dst := c.canvas
	// 0 means src, 1 means dst.
	l0, l1 := src.LocalBounds(), dst.LocalBounds()
	x0, y0, x1, y1 := 0, 0, x, y
	i0, i1 := src.PixOffset(x0, y0), dst.PixOffset(x1, y1)
	s0, s1 := src.Stride(), dst.Stride()
	p0, p1 := src.Pix(), dst.Pix()

	// TODO(BUG)
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

			// Alpha value
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
