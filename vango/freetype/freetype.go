// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

import (
	"errors"
	"image"
	"image/draw"
)

const (
	kGlyphNum      = 256
	kXFractionsNum = 4
	kYFractionsNum = 1
)

type cache_entry_t struct {
	valid  bool
	glyph  uint16
	mask   *image.Alpha
	offset image.Point
}

func ParseFont(b []byte) (*Font, error) {
	return Parse(b)
}

func Point(x, y int) RastPoint {
	return RastPoint{
		X: Fix32(x << 8),
		Y: Fix32(y << 8),
	}
}

type Context struct {
	rast           *Rast
	font           *Font
	glyph          *Glyph
	clip           image.Rectangle
	dst            draw.Image
	src            image.Image
	font_size, dpi float64
	scale          int32
	cache          [kGlyphNum * kXFractionsNum * kYFractionsNum]cache_entry_t
}

func (c *Context) PointToFix32(x float64) Fix32 {
	return Fix32(x * float64(c.dpi) * (256.0 / 72.0))
}

func (c *Context) drawContour(pt_array []FontPoint, dx, dy Fix32) {
	if len(pt_array) == 0 {
		return
	}

	start := RastPoint{
		X: dx + Fix32(pt_array[0].X<<2),
		Y: dy - Fix32(pt_array[0].Y<<2),
	}

	c.rast.Start(start)
	q0, on0 := start, true
	for _, pt := range pt_array[1:] {
		q := RastPoint{
			X: dx + Fix32(pt.X<<2),
			Y: dy - Fix32(pt.Y<<2),
		}
		on := pt.Flag&0x01 != 0
		if on {
			if on0 {
				c.rast.Add1(q)
			} else {
				c.rast.Add2(q0, q)
			}
		} else {
			if on0 {
				// empty
			} else {
				mid := RastPoint{
					X: (q0.X + q.X) / 2,
					Y: (q0.Y + q.Y) / 2,
				}
				c.rast.Add2(q0, mid)
			}
		}
		q0, on0 = q, on
	}
	// Close the curve.
	if on0 {
		c.rast.Add1(start)
	} else {
		c.rast.Add2(q0, start)
	}
}

func (c *Context) rasterize(glyph uint16, fx, fy Fix32) (*image.Alpha, image.Point, error) {
	if err := c.glyph.Load(c.font, c.scale, glyph, nil); err != nil {
		return nil, image.ZP, err
	}

	xmin := int(fx+Fix32(c.glyph.Rect.XMin<<2)) >> 8
	ymin := int(fy-Fix32(c.glyph.Rect.YMax<<2)) >> 8
	xmax := int(fx+Fix32(c.glyph.Rect.XMax<<2)+0xff) >> 8
	ymax := int(fy-Fix32(c.glyph.Rect.YMin<<2)+0xff) >> 8
	if xmin > xmax || ymin > ymax {
		return nil, image.ZP, errors.New("freetype negative sized glyph")
	}

	fx += Fix32(-xmin << 8)
	fy += Fix32(-ymin << 8)

	c.rast.Clear()

	e0 := 0

	for _, e1 := range c.glyph.EndIndexArray {
		c.drawContour(c.glyph.AllPoints[e0:e1], fx, fy)
		e0 = e1
	}

	a := image.NewAlpha(image.Rect(0, 0, xmax-xmin, ymax-ymin))
	c.rast.Rast(NewAlphaSrcDrawer(a))

	return a, image.Point{xmin, ymin}, nil
}

func (c *Context) glyph_at(glyph uint16, pt RastPoint) (*image.Alpha, image.Point, error) {
	ix, fx := int(pt.X>>8), pt.X&0xff
	iy, fy := int(pt.Y>>8), pt.Y&0xff

	tg := int(glyph) % kGlyphNum
	tx := int(fx) / (256 / kXFractionsNum)
	ty := int(fy) / (256 / kYFractionsNum)
	t := ((tg*kXFractionsNum)+tx)*kYFractionsNum + ty

	if c.cache[t].valid && c.cache[t].glyph == glyph {
		return c.cache[t].mask, c.cache[t].offset.Add(image.Point{ix, iy}), nil
	}

	mask, offset, err := c.rasterize(glyph, fx, fy)
	if err != nil {
		return nil, image.ZP, err
	}

	c.cache[t] = cache_entry_t{true, glyph, mask, offset}
	return mask, offset.Add(image.Point{ix, iy}), nil
}

func (c *Context) DrawString(str string, pt RastPoint) (RastPoint, error) {
	if c.font == nil {
		return RastPoint{}, errors.New("freetype DrawString called with nil font.")
	}

	prev, has_prev := uint16(0), false
	for _, rune := range str {
		idx := c.font.Index(rune)
		if has_prev {
			pt.X += Fix32(c.font.Kerning(c.scale, prev, idx)) << 2
		}

		mask, offset, err := c.glyph_at(idx, pt)
		if err != nil {
			return RastPoint{}, err
		}

		pt.X += Fix32(c.font.HMetric(c.scale, idx).AdvanceWidth) << 2
		glyph_rect := mask.Bounds().Add(offset)
		dr := c.clip.Intersect(glyph_rect)
		if !dr.Empty() {
			mp := image.Point{0, dr.Min.Y - glyph_rect.Min.Y}
			draw.DrawMask(c.dst, dr, c.src, image.ZP, mask, mp, draw.Over)
		}
		prev, has_prev = idx, true
	}
	return pt, nil
}

func (c *Context) recalc() {
	c.scale = int32(c.font_size * c.dpi * (64.0 / 72.0))
	if c.font == nil {
		c.rast.SetBounds(0, 0)
	} else {
		b := c.font.Bounds(c.scale)
		xmin := +int(b.XMin) >> 6
		ymin := -int(b.YMax) >> 6
		xmax := +int(b.XMax+63) >> 6
		ymax := -int(b.YMin-63) >> 6
		c.rast.SetBounds(xmax-xmin, ymax-ymin)
	}

	for i := range c.cache {
		c.cache[i] = cache_entry_t{}
	}
}

func (c *Context) SetDPI(dpi float64) {
	if c.dpi == dpi {
		return
	}

	c.dpi = dpi
	c.recalc()
}

func (c *Context) SetFont(font *Font) {
	if c.font == font {
		return
	}
	c.font = font
	c.recalc()
}

func (c *Context) SetFontSize(font_size float64) {
	if c.font_size == font_size {
		return
	}
	c.font_size = font_size
	c.recalc()
}

func (c *Context) SetDst(dst draw.Image) {
	c.dst = dst
}

func (c *Context) SetSrc(src image.Image) {
	c.src = src
}

func (c *Context) SetClip(clip image.Rectangle) {
	c.clip = clip
}

func NewContext() *Context {
	return &Context{
		rast:      NewRast(0, 0),
		glyph:     NewGlyph(),
		font_size: 12,
		dpi:       72,
		scale:     12 << 6,
	}
}
