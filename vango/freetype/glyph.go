// Copyright 2010 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// any later version), both of which can be found in the LICENSE file.

package freetype

import (
	"errors"
	// "log"
)

// control point.
type FontPoint struct {
	X, Y int32
	// The Flags' LSB means whether or not this Point is "on" the contour. Other
	// bits are reserved for internal use.
	Flag uint32
}

// An HMetric holds the horizontal metrics of a single glyph.
type HMetric struct {
	AdvanceWidth    int32
	LeftSideBearing int32
}

// A Glyph holds a glyph's contours. A GlyphBuf can be re-used to load a
// series of glyphs from a Font.
type Glyph struct {
	// B is the glyph's bounding box.
	Rect Bounds
	// Point contains all Points from all contours of the glyph. If a execer was
	// used to load a glyph then UnExec contains those Points before they were
	// execed, and Raw contains those Points before they were execed and
	// scaled.
	AllPoints    []FontPoint
	UnexecPoints []FontPoint
	RawRoints    []FontPoint

	// EndIndexArray is the point indexes of the end point of each countour. The
	// Length of EndIndexArray is the number of contours in the glyph. The i'th
	// contour consists of points Point[EndIndexArray[i-1]:EndIndexArray[i]],
	// where EndIndexArray[-1] is interpreted to mean zero.
	EndIndexArray []int

	font  *Font
	exec  *exec_t
	scale int32

	// pp1x is the x co-ordinate of the first phantom point.
	pp1x int32

	is_metrics_set bool
}

// Flags for decoding a glyph's contours. These flags are documented at
// http://developer.apple.com/fonts/TTRefMan/RM06/Chap6glyf.html
const (
	kDecodeOnCurve = 1 << iota
	kDecodeXShortVector
	kDecodeYShortVector
	kDecodeRepeat
	kDecodePositiveXShortVector
	kDecodePositiveYShortVector

	// For internal use.
	kDecodeTouchedX
	kDecodeTouchedY
)

// The same flag bits (0x10 and 0x20) are overloaded to have two meanings,
// dependent on the value of the k{X,Y}ShortVector bits.
const (
	kDecodeThisXIsSame = kDecodePositiveXShortVector
	kDecodeThisYIsSame = kDecodePositiveYShortVector
)

// Load loads a glyph's contours from a Font, overwriting any previously loaded
// contours for this GlyphBuf. scale is the number of 26.6 fixed point units in
// 1 em. The execer is optional; if non-nil, then the resulting glyph will be
// execed by the Font's bytecode instructions.
func (g *Glyph) Load(font *Font, scale int32, idx uint16, exec *exec_t) error {
	g.Rect = Bounds{}
	g.AllPoints = g.AllPoints[:0]
	g.UnexecPoints = g.UnexecPoints[:0]
	g.RawRoints = g.RawRoints[:0]
	g.EndIndexArray = g.EndIndexArray[:0]
	g.font = font
	g.scale = scale
	g.exec = exec
	g.pp1x = 0
	g.is_metrics_set = false

	if exec != nil {
		// log.Printf("F %v", g.AllPoints)
		if err := exec.init(font, scale); err != nil {
			return err
		}
		// log.Printf("G %v", g.AllPoints)
	}
	if err := g.load_impl(0, idx, true); err != nil {
		return err
	}
	if g.pp1x != 0 {
		for i := range g.AllPoints {
			g.AllPoints[i].X -= g.pp1x
		}
		// TODO: also adjust g.Rect?
	}

	return nil
}

func (g *Glyph) load_impl(recursion int32, idx uint16, use_my_metrics bool) (err error) {
	// The recursion limit here is arbitrary, but defends against malformed
	// glyphs.
	if recursion >= 32 {
		return errors.New("UNSUPPORT: excessive compound glyph recursion.")
	}
	// Find the relevant slice of gly.font.glyf
	var g0, g1 uint32
	if g.font.loca_offset_format == kLocaOffsetFormatShort {
		g0 = 2 * uint32(octets_to_u16(g.font.loca, 2*int(idx)))
		g1 = 2 * uint32(octets_to_u16(g.font.loca, 2*int(idx)+2))
	} else {
		g0 = octets_to_u32(g.font.loca, 4*int(idx))
		g1 = octets_to_u32(g.font.loca, 4*int(idx)+4)
	}

	if g0 == g1 {
		return nil
	}
	glyf := g.font.glyf[g0:g1]
	// log.Printf("glyf %v", glyf)
	// Decode the contour end indices.
	contour_num := int(int16(octets_to_u16(glyf, 0)))
	rect := Bounds{
		XMin: int32(int16(octets_to_u16(glyf, 2))),
		YMin: int32(int16(octets_to_u16(glyf, 4))),
		XMax: int32(int16(octets_to_u16(glyf, 6))),
		YMax: int32(int16(octets_to_u16(glyf, 8))),
	}
	mtrc, pp1x := g.font.unscaled_hmetric(idx), int32(0)
	if contour_num < 0 {
		if contour_num != -1 {
			// http://developer.apple.com/fonts/TTRefMan/RM06/Chap6glyf.html
			// says that "the values -2, -3, and so forth, are reserved for
			// future use."
			return errors.New("UNSUPPORT: negative number of contours.")
		}

		pp1x = g.font.scale(g.scale * (rect.XMin - mtrc.LeftSideBearing))

		if err := g.load_compound(recursion, glyf, use_my_metrics); err != nil {
			return err
		}
	} else {
		np0, ne0 := len(g.AllPoints), len(g.EndIndexArray)
		//log.Printf("Points A %v", g.AllPoints)
		pgm := g.load_simple(glyf, contour_num)
		//log.Printf("Points B %v", g.AllPoints)
		// Set the four phantom points. Freetype-Go uses only the first two,
		// but the exec bytecode may expect four.
		g.AllPoints = append(g.AllPoints,
			FontPoint{X: rect.XMin - mtrc.LeftSideBearing},
			FontPoint{X: rect.XMin - mtrc.LeftSideBearing + mtrc.AdvanceWidth},
			FontPoint{},
			FontPoint{})
		// Scale and exec the glyph.
		if g.exec != nil {
			g.RawRoints = append(g.RawRoints, g.AllPoints[np0:]...)
		}
		for i := np0; i < len(g.AllPoints); i++ {
			pt := &g.AllPoints[i]
			pt.X = g.font.scale(g.scale * pt.X)
			pt.Y = g.font.scale(g.scale * pt.Y)
		}

		if g.exec != nil {
			g.UnexecPoints = append(g.UnexecPoints, g.AllPoints[np0:]...)
			if len(pgm) != 0 {
				// log.Printf("B %v", g.AllPoints)
				err := g.exec.exec(pgm, g.AllPoints[np0:],
					g.UnexecPoints[np0:], g.RawRoints[np0:],
					g.EndIndexArray[ne0:])
				// log.Printf("C %v", g.AllPoints)
				if err != nil {
					return err
				}
			}

		}
		// Drop the four phantom points.
		pp1x = g.AllPoints[len(g.AllPoints)-4].X
		g.AllPoints = g.AllPoints[:len(g.AllPoints)-4]
		if g.exec != nil {
			g.RawRoints = g.RawRoints[:len(g.RawRoints)-4]
			g.UnexecPoints = g.UnexecPoints[:len(g.UnexecPoints)-4]
			if dx := (pp1x+32)&^63 - pp1x; dx != 0 {
				for i := np0; i < len(g.AllPoints); i++ {
					g.AllPoints[i].X += dx
				}
			}
		}
		if np0 != 0 {
			// The executing program expects the []EndIndexArray values to be
			// indexed relative to the inner glyph, not the outer glyph, so we
			// delayadding np0 units after the executing program (if any) has
			// run.
			for i := ne0; i < len(g.EndIndexArray); i++ {
				g.EndIndexArray[i] += np0
			}
		}
	}
	if use_my_metrics && !g.is_metrics_set {
		g.is_metrics_set = true
		g.Rect.XMin = g.font.scale(g.scale * rect.XMin)
		g.Rect.YMin = g.font.scale(g.scale * rect.YMin)
		g.Rect.XMax = g.font.scale(g.scale * rect.XMax)
		g.Rect.YMax = g.font.scale(g.scale * rect.YMax)
		g.pp1x = pp1x
	}
	return nil
}

// kLoadOffset is the initial offset for load_simple and load_compound. The
// first 1- bytes are the number of contours and the bounding box.
const kLoadOffset = 10

func (g *Glyph) load_simple(glybuf []byte, ne int) (pgm []byte) {
	offset := kLoadOffset

	for i := 0; i < ne; i++ {
		g.EndIndexArray = append(g.EndIndexArray, int(octets_to_u16(glybuf, offset))+1)
		offset += 2
	}

	// Note the truetype execing instructions.
	instructions_length := int(octets_to_u16(glybuf, offset))
	offset += 2
	pgm = glybuf[offset : offset+instructions_length]
	offset += instructions_length

	np0 := len(g.AllPoints)
	np1 := np0 + int(g.EndIndexArray[len(g.EndIndexArray)-1])

	// Decode the flags.
	for i := np0; i < np1; {
		code := uint32(glybuf[offset])
		offset++

		g.AllPoints = append(g.AllPoints, FontPoint{Flag: code})
		i++

		if code&kDecodeRepeat != 0 {
			count := glybuf[offset]
			offset++
			for ; count > 0; count-- {
				g.AllPoints = append(g.AllPoints, FontPoint{Flag: code})
				i++
			}
		}
	}

	// Decode the co-ordinates.
	var x int16
	for i := np0; i < np1; i++ {
		flag := g.AllPoints[i].Flag
		if flag&kDecodeXShortVector != 0 {
			dx := int16(glybuf[offset])
			offset++
			if flag&kDecodePositiveXShortVector == 0 {
				x -= dx
			} else {
				x += dx
			}
		} else if flag&kDecodeThisXIsSame == 0 {
			x += int16(octets_to_u16(glybuf, offset))
			offset += 2
		}
		g.AllPoints[i].X = int32(x)
	}

	var y int16
	for i := np0; i < np1; i++ {
		flag := g.AllPoints[i].Flag
		if flag&kDecodeYShortVector != 0 {
			dy := int16(glybuf[offset])
			offset++
			if flag&kDecodePositiveYShortVector == 0 {
				y -= dy
			} else {
				y += dy
			}
		} else if flag&kDecodeThisYIsSame == 0 {
			// log.Printf("Point %v Offset %v", g.AllPoints[i], offset)
			y += int16(octets_to_u16(glybuf, offset))
			offset += 2
		}
		g.AllPoints[i].Y = int32(y)
	}

	return pgm
}

func (gly *Glyph) load_compound(recursion int32, glybuf []byte, use_my_metrics bool) error {
	// Flags for decoding a compound glyph. These flags are documented at
	// http://developer.apple.com/fonts/TTRefMan/RM06/Chap6glyf.html
	const (
		kArg1AndArg2AreWords = 1 << iota
		kArgsAreXYValues
		kRoundXYToGrid
		kWeHaveAScale
		kUnunsed
		kMoreComponents
		kWeHaveAnXAndYScale
		kWeHaveATwoByTwo
		kWeHaveInstructions
		kUseMyMetrics
		kOverlavCompound
	)

	for offset := kLoadOffset; ; {
		flag := octets_to_u16(glybuf, offset)
		offset += 2
		component := octets_to_u16(glybuf, offset)
		offset += 2

		dx, dy, transform, has_transform := int32(0), int32(0), [4]int32{}, false
		if flag&kArg1AndArg2AreWords != 0 {
			dx = int32(int16(octets_to_u16(glybuf, offset)))
			dy = int32(int16(octets_to_u16(glybuf, offset+2)))
			offset += 4
		} else {
			dx = int32(int16(int8(glybuf[offset])))
			dy = int32(int16(int8(glybuf[offset+1])))
			offset += 2
		}

		if flag&kArgsAreXYValues == 0 {
			return errors.New("FT UNSUPPORT: compound ")
		}

		if flag&(kWeHaveAScale|kWeHaveAnXAndYScale|kWeHaveATwoByTwo) != 0 {
			has_transform = true
			if flag&kWeHaveAScale != 0 {
				transform[0] = int32(int16(octets_to_u16(glybuf, offset)))
				offset += 2
				transform[3] = transform[0]
			} else if flag&kWeHaveAnXAndYScale != 0 {
				transform[0] = int32(int16(octets_to_u16(glybuf, offset)))
				transform[3] = int32(int16(octets_to_u16(glybuf, offset+2)))
				offset += 4
			} else if flag&kWeHaveATwoByTwo != 0 {
				transform[0] = int32(int16(octets_to_u16(glybuf, offset)))
				transform[1] = int32(int16(octets_to_u16(glybuf, offset+2)))
				transform[2] = int32(int16(octets_to_u16(glybuf, offset+4)))
				transform[3] = int32(int16(octets_to_u16(glybuf, offset+6)))
			}
		}

		np0 := len(gly.AllPoints)
		component_umm := use_my_metrics && (flag&kUseMyMetrics != 0)
		if err := gly.load_impl(recursion+1, component, component_umm); err != nil {
			return err
		}

		if has_transform {
			for i := np0; i < len(gly.AllPoints); i++ {
				pt := &gly.AllPoints[i]
				new_x := int32((int64(pt.X)*int64(transform[0])+1<<13)>>14) +
					int32((int64(pt.Y)*int64(transform[2])+1<<13)>>14)
				new_y := int32((int64(pt.X)*int64(transform[1])+1<<13)>>14) +
					int32((int64(pt.Y)*int64(transform[3])+1<<13)>>14)
				pt.X, pt.Y = new_x, new_y
			}
		}

		dx = gly.font.scale(gly.scale * dx)
		dy = gly.font.scale(gly.scale * dy)
		if flag&kRoundXYToGrid != 0 {
			dx = (dx + 32) &^ 63
			dy = (dy + 32) &^ 63
		}

		for i := np0; i < len(gly.AllPoints); i++ {
			pt := &gly.AllPoints[i]
			pt.X += dx
			pt.Y += dy
		}

		// TODO: also adjust gly.RawPoints and gly.UnexecPoints?
		if flag&kMoreComponents == 0 {
			break
		}
	}

	// TODO: exec the compound glyph.
	return nil
}

// TODO: is this necessary? The zero-valued Glyph is perfectly useable.

func NewGlyph() *Glyph {
	return &Glyph{
		AllPoints:     make([]FontPoint, 0, 256),
		EndIndexArray: make([]int, 0, 32),
	}
}
