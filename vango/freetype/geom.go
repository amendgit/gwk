// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

import (
	"fmt"
	"math"
)

type Bounds struct {
	XMin, YMin, XMax, YMax int32
}

// A Fix32 is a 24.8 fixed point number.
type Fix32 int32

// A Fix64 us a 48.16 fixed point number.
type Fix64 int64

func (f Fix32) String() string {
	if f < 0 {
		f = -f
		return fmt.Sprintf("-%d:%03d", int32(f/256), int32(f%256))
	}
	return fmt.Sprintf("%d:%03d", int32(f/256), int32(f%256))
}

func (f Fix64) String() string {
	if f < 0 {
		f = -f
		return fmt.Sprintf("-%d:%05d", int64(f/65536), int64(f%65536))
	}
	return fmt.Sprintf("%d:%05d", int64(f/65536), int64(f%65536))
}

func max_abs(arg0, arg1 Fix32) Fix32 {
	if arg0 < 0 {
		arg0 = -arg0
	}

	if arg1 < 0 {
		arg1 = -arg1
	}

	if arg0 < arg1 {
		return arg1
	} else {
		return arg0
	}
}

type RastPoint struct {
	X Fix32
	Y Fix32
}

func (pt RastPoint) String() string {
	return "(" + pt.X.String() + ", " + pt.Y.String() + ")"
}

func (pt RastPoint) Add(arg RastPoint) RastPoint {
	return RastPoint{pt.X + arg.X, pt.Y + arg.Y}
}

func (pt RastPoint) Sub(arg RastPoint) RastPoint {
	return RastPoint{pt.X - arg.X, pt.Y - arg.Y}
}

func (pt RastPoint) Mul(k Fix32) RastPoint {
	return RastPoint{pt.X * k / 256, pt.Y * k / 256}
}

func (pt RastPoint) Neg() RastPoint {
	return RastPoint{-pt.X, -pt.Y}
}

func (pt RastPoint) Dot(arg RastPoint) Fix64 {
	x0, y0 := int64(pt.X), int64(pt.Y)
	x1, y1 := int64(arg.X), int64(arg.Y)
	return Fix64(x0*x1 + y0*y1)
}

func (pt RastPoint) Len() Fix32 {
	x := float64(pt.X)
	y := float64(pt.Y)
	return Fix32(math.Sqrt(x*x + y*y))
}

func (pt RastPoint) Normalize(length Fix32) RastPoint {
	l := pt.Len()
	if l == 0 {
		return RastPoint{0, 0}
	}
	s, t := int64(length), int64(l)
	x := int64(pt.X) * s / t
	y := int64(pt.Y) * s / t
	return RastPoint{Fix32(x), Fix32(y)}
}

func (pt RastPoint) Rotate(r int) RastPoint {
	x0, y0 := int64(pt.X), int64(pt.Y)

	var x1, y1 int64

	switch r {
	case 45:
		x1 = (x0 - y0) * 181 / 256
		y1 = (x0 + y0) * 181 / 256
	case 90:
		x1 = int64(-pt.Y)
		y1 = int64(+pt.X)
	case 135:
		x1 = (-x0 - y0) * 181 / 256
		y1 = (+x0 + y0) * 181 / 256
	case -45:
		x1 = (+x0 + y0) * 181 / 256
		y1 = (-x0 + y0) * 181 / 256
	case -90:
		x1 = int64(+pt.Y)
		y1 = int64(-pt.X)
	case -135:
		x1 = (-x0 + y0) * 181 / 256
		y1 = (-x0 - y0) * 181 / 256
	}

	return RastPoint{Fix32(x1), Fix32(y1)}
}

type Adder interface {
	Start(pt RastPoint)
	Add1(arg0 RastPoint)
	Add2(arg0, arg1 RastPoint)
	Add3(arg0, arg1, arg2 RastPoint)
}

type Path []Fix32

func (p Path) String() string {
	s := ""
	for i := 0; i < len(p); {
		if i != 0 {
			s += " "
		}
		switch p[i] {
		case 0:
			s += "S0" + fmt.Sprint([]Fix32(p[i+1:i+3]))
			i += 4
		case 1:
			s += "A1" + fmt.Sprint([]Fix32(p[i+1:i+3]))
			i += 4
		case 2:
			s += "A2" + fmt.Sprint([]Fix32(p[i+1:i+5]))
			i += 6
		case 3:
			s += "A3" + fmt.Sprint([]Fix32(p[i+1:i+7]))
			i += 8
		default:
			panic("FONT bad path")
		}
	}
	return s
}

func (p *Path) grow(n int) {
	n = n + len(*p)
	if n > cap(*p) {
		old := *p
		*p = make([]Fix32, n, 2*n+8)
		copy(*p, old)
		return
	}
	*p = (*p)[0:n]
}

func (p *Path) Clear() {
	*p = (*p)[0:0]
}

func (p *Path) Start(pt RastPoint) {
	n := len(*p)
	p.grow(4)
	(*p)[n+0] = 0
	(*p)[n+1] = pt.X
	(*p)[n+2] = pt.Y
	(*p)[n+3] = 0
}

func (p *Path) Add1(pt0 RastPoint) {
	n := len(*p)
	p.grow(4)
	(*p)[n+0] = 1
	(*p)[n+1] = pt0.X
	(*p)[n+2] = pt0.Y
	(*p)[n+3] = 1
}

func (p *Path) Add2(pt0, pt1 RastPoint) {
	n := len(*p)
	p.grow(6)
	(*p)[n+0] = 2
	(*p)[n+1] = pt0.X
	(*p)[n+2] = pt0.Y
	(*p)[n+3] = pt1.X
	(*p)[n+4] = pt1.Y
	(*p)[n+5] = 2
}

func (p *Path) Add3(pt0, pt1, pt2 RastPoint) {
	n := len(*p)
	p.grow(8)
	(*p)[n+0] = 3
	(*p)[n+1] = pt0.X
	(*p)[n+2] = pt0.Y
	(*p)[n+3] = pt1.X
	(*p)[n+4] = pt1.Y
	(*p)[n+5] = pt2.X
	(*p)[n+6] = pt2.Y
	(*p)[n+7] = 3
}

func (p *Path) AddPath(p0 Path) {
	n0, n1 := len(*p), len(p0)
	p.grow(n1)
	copy((*p)[n0:n0+n1], p0)
}

func (p *Path) AddStroke(path Path, width Fix32, capper Capper, joiner Joiner) {
	Stroke(p, path, width, capper, joiner)
}

func (p Path) first_point() RastPoint {
	return RastPoint{p[1], p[2]}
}

func (p Path) last_point() RastPoint {
	return RastPoint{p[len(p)-3], p[len(p)-2]}
}

func reverse_add_path(adder Adder, path Path) {
	if len(path) == 0 {
		return
	}

	i := len(path) - 1
	for {
		switch path[i] {
		case 0:
			return
		case 1:
			i -= 4
			path.Add1(RastPoint{path[i-2], path[i-1]})
		case 2:
			i -= 6
			pt0 := RastPoint{path[i+2], path[i+3]}
			pt1 := RastPoint{path[i-2], path[i-1]}
			path.Add2(pt0, pt1)
		case 3:
			i -= 8
			pt0 := RastPoint{path[i+4], path[i+5]}
			pt1 := RastPoint{path[i+2], path[i+3]}
			pt2 := RastPoint{path[i-2], path[i-1]}
			path.Add3(pt0, pt1, pt2)
		default:
			panic("FONT geom bad path")
		}
	}
}
