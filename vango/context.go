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

func (c *Context) DrawText(text string, rect image.Rectangle) (freetype.RastPoint, error) {
	if c.font == nil {
		return freetype.RastPoint{}, errors.New("vango DrawString called with nil font.")
	}

	pt := freetype.Point(rect.Min.X, rect.Min.Y)

	prev, has_prev := uint16(0), false
	for _, rune := range text {
		idx := c.font.Index(rune)
		if has_prev {
			pt.X += freetype.Fix32(c.font.Kerning(1, prev, idx)) << 2
		}

		mask, offset, _ := c.font.GlyphAt(idx, pt)
		// if err != nil {
		// 	return RastPoint{}, err
		// }

		pt.X += freetype.Fix32(c.font.font.HMetric(1.0, idx).AdvanceWidth) << 2
		glyph_rect := mask.Bounds().Add(offset)
		log.Printf("%v", glyph_rect)
		// dr := c.clip.Intersect(glyph_rect)
		//if !dr.Empty() {
		// mp := image.Point{0, dr.Min.Y - glyph_rect.Min.Y}
		// draw.DrawMask(c.dst, dr, c.src, image.ZP, mask, mp, draw.Over)
		// c.DrawImage(rect.Min.X, rect.Min.Y, mask, glyph_rect)
		//}
		prev, has_prev = idx, true
	}
	return pt, nil
}
