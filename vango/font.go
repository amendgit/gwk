package vango

import (
	"errors"
	"gwk/vango/freetype"
	"image"
	"io/ioutil"
	"log"
	"os"
)

var g_default_font *freetype.Font

func init_font() {
	wd, _ := os.Getwd()
	log.Printf("WD %v", wd)
	bytes, err := ioutil.ReadFile("./resc/luxisr.ttf")
	if err != nil {
		log.Printf("error: load font failed -> %v", err)
		return
	}

	g_default_font, err = freetype.ParseFont(bytes)
	if g_default_font == nil {
		log.Printf("err: load font failed -> %v", err)
		return
	}
}

var g_description_font_map map[string]*Font

func find_font_by_description(description string) *Font {
	return g_description_font_map[description]
}

const (
	kGlyphNum      = 256
	kXFractionsNum = 4
	kYFractionsNum = 1
)

type glyph_cache_t struct {
	valid  bool
	glyph  uint16
	mask   *image.Alpha
	offset image.Point
}

// vango.Font is a wrapper of freetype.Font. It is more like freetype.Context
type Font struct {
	font  *freetype.Font
	size  float64
	cache [kGlyphNum * kXFractionsNum * kYFractionsNum]glyph_cache_t
	glyph *freetype.Glyph
	rast  *freetype.Rast
	dpi   float64
	scale int32
}

func NewFont() *Font {
	f := &Font{
		rast:  freetype.NewRast(0, 0),
		glyph: freetype.NewGlyph(),
		size:  12,
		font:  g_default_font,
		dpi:   72,
	}

	f.recalc()
	return f
}

func (f *Font) GlyphAt(glyph uint16, pt freetype.RastPoint) (*image.Alpha, image.Point, error) {
	ix, fx := int(pt.X>>8), pt.X&0xff
	iy, fy := int(pt.Y>>8), pt.Y&0xff

	tg := int(glyph) % kGlyphNum
	tx := int(fx) / (256 / kXFractionsNum)
	ty := int(fy) / (256 / kYFractionsNum)
	t := ((tg*kXFractionsNum)+tx)*kYFractionsNum + ty

	if f.cache[t].valid && f.cache[t].glyph == glyph {
		return f.cache[t].mask, f.cache[t].offset.Add(image.Point{ix, iy}), nil
	}
	log.Printf("glyph %v", glyph)
	mask, offset, err := f.rasterize(glyph, fx, fy)
	if err != nil {
		return nil, image.ZP, err
	}

	f.cache[t] = glyph_cache_t{true, glyph, mask, offset}
	return mask, offset.Add(image.Point{ix, iy}), nil
}

func (f *Font) Index(ch rune) uint16 {
	return f.font.Index(ch)
}

func (f *Font) Kerning(scale int32, i0, i1 uint16) int32 {
	return f.font.Kerning(scale, i0, i1)
}

func (f *Font) rasterize(glyph uint16, fx, fy freetype.Fix32) (*image.Alpha, image.Point, error) {
	if f.glyph != nil {
		err := f.glyph.Load(f.font, f.scale, glyph, nil)
		if err != nil {
			return nil, image.ZP, err
		}
	}

	xmin := int(fx+freetype.Fix32(f.glyph.Rect.XMin<<2)) >> 8
	ymin := int(fy-freetype.Fix32(f.glyph.Rect.YMax<<2)) >> 8
	xmax := int(fx+freetype.Fix32(f.glyph.Rect.XMax<<2)+0xff) >> 8
	ymax := int(fy-freetype.Fix32(f.glyph.Rect.YMin<<2)+0xff) >> 8

	if xmin > xmax || ymin > ymax {
		return nil, image.ZP, errors.New("vango negative sized glyph")
	}

	fx += freetype.Fix32(-xmin << 8)
	fy += freetype.Fix32(-ymin << 8)

	f.rast.Clear()

	e0 := 0
	for _, e1 := range f.glyph.EndIndexArray {
		f.draw_contour(f.glyph.AllPoints[e0:e1], fx, fy)
		e0 = e1
	}

	a := image.NewAlpha(image.Rect(0, 0, xmax-xmin, ymax-ymin))
	f.rast.Rast(freetype.NewAlphaSrcDrawer(a))

	return a, image.Point{xmin, ymin}, nil
}

func (f *Font) draw_contour(pt_array []freetype.FontPoint, dx, dy freetype.Fix32) {
	if len(pt_array) == 0 {
		return
	}

	start := freetype.RastPoint{
		X: dx + freetype.Fix32(pt_array[0].X<<2),
		Y: dy - freetype.Fix32(pt_array[0].Y<<2),
	}

	f.rast.Start(start)
	q0, on0 := start, true
	for _, pt := range pt_array[1:] {
		q := freetype.RastPoint{
			X: dx + freetype.Fix32(pt.X<<2),
			Y: dy - freetype.Fix32(pt.Y<<2),
		}
		on := pt.Flag&0x01 != 0
		if on {
			if on0 {
				f.rast.Add1(q)
			} else {
				f.rast.Add2(q0, q)
			}
		} else {
			if on0 {
				// empty
			} else {
				mid := freetype.RastPoint{
					X: (q0.X + q.X) / 2,
					Y: (q0.Y + q.Y) / 2,
				}
				f.rast.Add2(q0, mid)
			}
		}
		q0, on0 = q, on
	}
	// close the curve.
	if on0 {
		f.rast.Add1(start)
	} else {
		f.rast.Add2(q0, start)
	}
}

func (f *Font) recalc() {
	f.scale = int32(f.size * f.dpi * (64.0 / 72.0))

	b := f.font.Bounds(f.scale)
	xmin := +int(b.XMin) >> 6
	ymin := -int(b.YMax) >> 6
	xmax := +int(b.XMax+63) >> 6
	ymax := -int(b.YMin-63) >> 6

	f.rast.SetBounds(xmax-xmin, ymax-ymin)
}
