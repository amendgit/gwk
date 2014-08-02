// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

import (
	"image"
	"image/color"
	"image/draw"

	//"log"
	"math"
)

// A Span is a horizontal segment of pixels with constant alpha. X0 is an
// inclusive bound and X1 is exclusive, the same as for slices. A fully
// opaque Span has A == 1<<32 - 1.
type Span struct {
	Y  int
	X0 int
	X1 int
	A  uint32
}

// A Drawer knows how to draw a batch of Spans. Rasterization may involve
// Drawing multple batches, and done will be true for the final batch.
// The Spans' Y values are monotonically increasing during a rasterizatio.
// Draw may use all of ss as scratch space during the call.
type Drawer interface {
	Draw(span []Span, done bool)
}

type DrawFunc func(s []Span, done bool)

func (fn DrawFunc) Draw(s []Span, done bool) {
	fn(s, done)
}

type AlphaOverDrawer struct {
	Image *image.Alpha
}

func (r AlphaOverDrawer) Draw(span_array []Span, done bool) {
	b := r.Image.Bounds()
	for _, s := range span_array {
		if s.Y < b.Min.Y {
			continue
		}

		if s.Y >= b.Max.Y {
			return
		}

		if s.X0 < b.Min.X {
			s.X0 = b.Min.X
		}

		if s.X1 > b.Max.X {
			s.X1 = b.Max.X
		}

		if s.X0 >= s.X1 {
			continue
		}

		base := (s.Y-r.Image.Rect.Min.Y)*r.Image.Stride - r.Image.Rect.Min.X
		p := r.Image.Pix[base+s.X0 : base+s.X1]
		a := int(s.A >> 24)

		for i, c := range p {
			v := int(c)
			p[i] = uint8((v*255 + (255-v)*a) / 255)
		}
	}
}

func NewAlphaOverDrawer(m *image.Alpha) AlphaOverDrawer {
	return AlphaOverDrawer{m}
}

type AlphaSrcDrawer struct {
	Image *image.Alpha
}

func (r AlphaSrcDrawer) Draw(span_array []Span, done bool) {
	b := r.Image.Bounds()
	for _, s := range span_array {
		//log.Printf("Alpha Src Draw Span %v %v %v %v", s.X0, s.X1, s.A, s.Y)
		if s.Y < b.Min.Y {
			continue
		}

		if s.Y >= b.Max.Y {
			return
		}

		if s.X0 < b.Min.X {
			s.X0 = b.Min.X
		}

		if s.X1 > b.Max.X {
			s.X1 = b.Max.X
		}

		if s.X0 >= s.X1 {
			continue
		}

		base := (s.Y-r.Image.Rect.Min.Y)*r.Image.Stride - r.Image.Rect.Min.X
		p := r.Image.Pix[base+s.X0 : base+s.X1]
		color := uint8(s.A >> 24)
		for i := range p {
			p[i] = color
		}
	}
}

func NewAlphaSrcDrawer(m *image.Alpha) AlphaSrcDrawer {
	return AlphaSrcDrawer{m}
}

type RGBADrawer struct {
	// The image to compose onto.
	Image *image.RGBA
	// The Porter-Duff composition operator.
	Op draw.Op
	// The 16-bit color to draw the spans.
	r, g, b, a uint32
}

func (r *RGBADrawer) Draw(span_array []Span, done bool) {
	b := r.Image.Bounds()

	for _, s := range span_array {
		// log.Printf("%v", s)
		if s.Y < b.Min.Y {
			continue
		}

		if s.Y >= b.Max.Y {
			return
		}

		if s.X0 < b.Min.X {
			s.X0 = b.Min.X
		}

		if s.X1 > b.Max.X {
			s.X1 = b.Max.X
		}

		if s.X0 >= s.X1 {
			continue
		}

		ma := s.A >> 16
		const kM = 1<<16 - 1
		i0 := (s.Y-r.Image.Rect.Min.Y)*r.Image.Stride + (s.X0-r.Image.Rect.Min.X)*4
		i1 := i0 + (s.X1-s.X0)*4
		if r.Op == draw.Over {
			for i := i0; i < i1; i += 4 {
				dr := uint32(r.Image.Pix[i+0])
				dg := uint32(r.Image.Pix[i+1])
				db := uint32(r.Image.Pix[i+2])
				da := uint32(r.Image.Pix[i+3])
				a := (kM - (r.a * ma / kM)) * 0x101
				r.Image.Pix[i+0] = uint8((dr*a + r.r*ma) / kM >> 8)
				r.Image.Pix[i+1] = uint8((dg*a + r.g*ma) / kM >> 8)
				r.Image.Pix[i+2] = uint8((db*a + r.b*ma) / kM >> 8)
				r.Image.Pix[i+3] = uint8((da*a + r.a*ma) / kM >> 8)
			}
		} else {
			for i := i0; i < i1; i += 4 {
				r.Image.Pix[i+0] = uint8(r.r * ma / kM >> 8)
				r.Image.Pix[i+1] = uint8(r.g * ma / kM >> 8)
				r.Image.Pix[i+2] = uint8(r.b * ma / kM >> 8)
				r.Image.Pix[i+3] = uint8(r.a * ma / kM >> 8)
			}
		}
	}
}

func (r *RGBADrawer) SetColor(c color.Color) {
	r.r, r.g, r.b, r.a = c.RGBA()
}

func NewRGBADrawer(m *image.RGBA) *RGBADrawer {
	return &RGBADrawer{Image: m}
}

type MonochromeDrawer struct {
	Drawer    Drawer
	y, x0, x1 int
}

func (m *MonochromeDrawer) Draw(span_array []Span, done bool) {
	j := 0
	for _, s := range span_array {
		if s.A >= 1<<31 {
			if m.y == s.Y && m.x1 == s.X0 {
				m.x1 = s.X1
			} else {
				span_array[j] = Span{m.y, m.x0, m.x1, 1<<32 - 1}
				j++
				m.y, m.x0, m.x1 = s.Y, s.X0, s.X1
			}
		}
	}
	if done {
		final_span := Span{m.y, m.x0, m.x1, 1<<32 - 1}
		if j < len(span_array) {
			span_array[j] = final_span
			j++
			m.Drawer.Draw(span_array[0:j], true)
		} else if j == len(span_array) {
			m.Drawer.Draw(span_array, false)
			if cap(span_array) > 0 {
				span_array = span_array[0:1]
			} else {
				span_array = make([]Span, 1)
			}
			span_array[0] = final_span
			m.Drawer.Draw(span_array, true)
		} else {
			panic("ft unreachable")
		}

		m.y, m.x0, m.x1 = 0, 0, 0
	} else {
		m.Drawer.Draw(span_array[0:j], false)
	}
}

func NewMonochromeDrawer(d Drawer) *MonochromeDrawer {
	return &MonochromeDrawer{Drawer: d}
}

type GammaCorrectionDrawer struct {
	Drawer       Drawer
	alpha_table  [256]uint16
	gamma_is_one bool
}

func (g *GammaCorrectionDrawer) Draw(span_array []Span, done bool) {
	if !g.gamma_is_one {
		const (
			kM = 0x1010101 // 255*M == 1<<32-1
			kN = 0x8080    // N = M>>9, and N < 1<<16-1
		)

		for i, _ := range span_array {
			if span_array[i].A == 0 || span_array[i].A == 1<<32-1 {
				continue
			}
			p, q := span_array[i].A/kM, (span_array[i].A%kM)>>9
			a := uint32(g.alpha_table[p])*(kN-q) + uint32(g.alpha_table[p+1])*q
			a = (a + kN/2) / kN
			a |= a << 16
			span_array[i].A = a
		}
	}
	g.Drawer.Draw(span_array, done)
}

func (g *GammaCorrectionDrawer) SetGamma(gamma float64) {
	if gamma == 1.0 {
		g.gamma_is_one = true
		return
	}
	g.gamma_is_one = false
	for i := 0; i < 256; i++ {
		a := float64(i) / 0xff
		a = math.Pow(a, gamma)
		g.alpha_table[i] = uint16(0xfff * a)
	}
}

func NewGammaCorrectionDrawer(d Drawer, gamma float64) *GammaCorrectionDrawer {
	g := &GammaCorrectionDrawer{Drawer: d}
	g.SetGamma(gamma)
	return g
}
