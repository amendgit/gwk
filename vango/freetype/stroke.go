// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

const kEpsilon = 16384

type Capper interface {
	Cap(adder Adder, half_width Fix32, pivot, pt RastPoint)
}

type CapperFunc func(Adder, Fix32, RastPoint, RastPoint)

func (fn CapperFunc) Cap(adder Adder, half_width Fix32, pivot, pt RastPoint) {
	fn(adder, half_width, pivot, pt)
}

type Joiner interface {
	Join(lhs, rhs Adder, halfwidth Fix32, pivot, pt0, pt1 RastPoint)
}

type JoinerFunc func(lhs, rhs Adder, halfwidth Fix32, pivot, pt0, pt1 RastPoint)

func (fn JoinerFunc) Join(lhs, rhs Adder, halfwidth Fix32, pivot, pt0, pt1 RastPoint) {
	fn(lhs, rhs, halfwidth, pivot, pt0, pt1)
}

var RoundCapper Capper = CapperFunc(round_capper)

func round_capper(adder Adder, half_width Fix32, pivot, pt1 RastPoint) {
	// The cubic Bézier approximation to a circle involves the magic number
	// (√2 - 1) * 4/3, which is approximately 141/256.
	const K = 141
	e := pt1.Rotate(-90)
	side := pivot.Add(e)
	beg, end := pivot.Sub(pt1), pivot.Add(pt1)
	d, e := pt1.Mul(K), e.Mul(K)
	adder.Add3(beg.Add(e), side.Sub(d), side)
	adder.Add3(side.Add(d), end.Add(e), end)
}

// ButtCapper adds butt caps to a stroked path.
var ButtCapper Capper = CapperFunc(butt_capper)

func butt_capper(adder Adder, half_width Fix32, pivot, pt1 RastPoint) {
	adder.Add1(pivot.Add(pt1))
}

var SquareCapper Capper = CapperFunc(square_capper)

func square_capper(adder Adder, half_width Fix32, pivot, pt1 RastPoint) {
	e := pt1.Rotate(-90)
	side := pivot.Add(e)
	adder.Add1(side.Sub(pt1))
	adder.Add1(side.Add(pt1))
	adder.Add1(pivot.Add(pt1))
}

// RoundJoiner adds round joins to a stroked path.
var RoundJoiner Joiner = JoinerFunc(round_joiner)

func round_joiner(lhs, rhs Adder, half_width Fix32, pivot, pt0, pt1 RastPoint) {
	dot := pt0.Rotate(90).Dot(pt1)
	if dot >= 0 {
		add_arc(lhs, pivot, pt0, pt1)
	} else {
		lhs.Add1(pivot.Add(pt1))
		add_arc(rhs, pivot, pt0.Neg(), pt1.Neg())
	}
}

// BevelJoiner adds bevel joins to a stroked path.
var BevelJoiner Joiner = JoinerFunc(bevel_joiner)

func bevel_joiner(lhs, rhs Adder, half_width Fix32, pivot, pt0, pt1 RastPoint) {
	lhs.Add1(pivot.Add(pt1))
	rhs.Add1(pivot.Sub(pt1))
}

// |add_arc| adds a circular arc from pivot + pt0 to pivot + pt1 to p. The shorter of
// the two possible arcs is taken, i.e. the one spanning <= 180 degrees. The two
// vectors pt0 and pt1 must be of equal length.
func add_arc(adder Adder, pivot, pt0, pt1 RastPoint) {
	// r2 is the square of the length of pt0.
	r2 := pt0.Dot(pt0)
	if r2 < kEpsilon {
		// The arc radius is so small that we collapse to a straight line.
		adder.Add1(pivot.Add(pt1))
		return
	}

	// We approximate the arc by 0, 1, 2 or 3 45-degree quadratic segments plus
	// a final quadratic segment from s to n1. Each 45-degree segment has control
	// points {1, 0}, {1, tan(π/8)} and {1/√2, 1/√2} suitably scaled, rotated and
	// translated. tan(π/8) is approximately 106/256.
	const kTpo8 = 106
	var s RastPoint
	// We determine which octant the angle between pt0 and pt1 is in via three dot
	// m0, m1 and m2 are pt0 rotated clockkwise by 45, 90 and 135 degrees.
	m0 := pt0.Rotate(45)
	m1 := pt0.Rotate(90)
	m2 := m0.Rotate(90)

	if m1.Dot(pt1) >= 0 {
		if pt0.Dot(pt1) >= 0 {
			if m2.Dot(pt1) <= 0 {
				s = pt0
			} else {
				adder.Add2(pivot.Add(pt0).Add(m1.Mul(kTpo8)), pivot.Add(m0))
				s = m0
			}
		} else {
			pm1, n0t := pivot.Sub(m1), pt0.Mul(kTpo8)
			adder.Add2(pivot.Add(pt0).Sub(m1.Mul(kTpo8)), pivot.Sub(m2))
			adder.Add2(pm1.Add(n0t), pm1)
			if m2.Dot(pt1) <= 0 {
				s = m1.Neg()
			} else {
				adder.Add2(pm1.Sub(n0t), pivot.Sub(m0))
				s = m0.Neg()
			}
		}
	}
	// The final quadratic segment has two endpoints s and n1 and the middle
	// control point is a multiple of s.Add(n1), i.e. it is on the angle bisector
	// of those two points. The multiple ranges between 128/256 and 150/256 as
	// the angle between s and n1 ranges between 0 and 45 degrees.
	// When the angle is 0 degrees (i.e. s and n1 are coincident) then s.Add(n1)
	// is twice s and so the middle control point of the degenerate quadratic
	// segment should be half s.Add(n1), and half = 128/256.
	// When the angle is 45 degrees then 150/256 is the ratio of the lengths of
	// the two vectors {1, tan(π/8)} and {1 + 1/√2, 1/√2}.
	// d is the normalized dot product between s and n1. Since the angle ranges
	// between 0 and 45 degrees then d ranges between 256/256 and 181/256.
	d := 256 * s.Dot(pt1) / r2
	multiple := Fix32(150 - 22*(d-181)/(256-181))
	adder.Add2(pivot.Add(s.Add(pt1).Mul(multiple)), pivot.Add(pt1))
}

// midpoint returns the midpoint of two RastPoints.
func midpoint(arg0, arg1 RastPoint) RastPoint {
	return RastPoint{(arg0.X + arg1.X) / 2, (arg0.Y + arg1.Y) / 2}
}

// angle_greater_than_45 returns whether the angle between two vectors is more
// than 45 degrees.
func angle_greater_than_45(v0, v1 RastPoint) bool {
	v := v0.Rotate((-45))
	return v.Dot(v1) < 0 || v.Rotate(90).Dot(v1) < 0
}

// interpolate returns the point (1-t)*a + t*b
func interpolate(arg0, arg1 RastPoint, t Fix64) RastPoint {
	s := 65535 - t
	x := s*Fix64(arg0.X) + t*Fix64(arg1.X)
	y := s*Fix64(arg0.Y) + t*Fix64(arg1.Y)
	return RastPoint{Fix32(x >> 16), Fix32(y >> 16)}
}

func curviest2(arg0, arg1, arg2 RastPoint) Fix64 {
	dx := int64(arg1.X - arg0.X)
	dy := int64(arg1.Y - arg0.Y)
	ex := int64(arg2.X - 2*arg1.X + arg0.X)
	ey := int64(arg2.Y - 2*arg1.Y + arg0.Y)
	if ex == 0 && ey == 0 {
		return 32768
	}
	return Fix64(-65535 * (dx*ex + dy*ey) / (ex*ex + ey*ey))
}

type stroke_state_t struct {
	adder      Adder
	half_width Fix32
	capper     Capper
	joiner     Joiner
	path       Path
	recent_pt  RastPoint
	normal_pt  RastPoint
}

func (s stroke_state_t) add_non_curvy2(arg_pt0, arg_pt1 RastPoint) {
	const kMaxDepth = 5
	var depth_stack [kMaxDepth + 1]int
	var point_stack [2*kMaxDepth + 3]RastPoint

	depth_stack[0] = 0
	point_stack[2] = s.recent_pt
	point_stack[1] = arg_pt0
	point_stack[0] = arg_pt1
	normal_pt0 := s.normal_pt
	var normal_pt2 RastPoint

	var top int
	for {
		depth := depth_stack[top]
		pt0 := point_stack[2*top+2]
		pt1 := point_stack[2*top+1]
		pt2 := point_stack[2*top+0]
		v01 := pt1.Sub(pt0)
		v12 := pt2.Sub(pt1)
		is_01_small := v01.Dot(v01) < Fix64(1<<16)
		is_12_small := v12.Dot(v12) < Fix64(1<<16)

		if is_01_small && is_12_small {
			normal_pt2 := v12.Normalize(s.half_width).Rotate(-90)
			mid02 := midpoint(pt0, pt2)
			add_arc(s.adder, mid02, normal_pt0, normal_pt2)
			add_arc(&s.path, mid02, normal_pt0.Neg(), normal_pt2.Neg())
		} else if depth < kMaxDepth && angle_greater_than_45(v01, v12) {
			mid01 := midpoint(pt0, pt1)
			mid12 := midpoint(pt1, pt2)

			top++

			depth_stack[top+0] = depth + 1
			depth_stack[top-1] = depth + 1

			point_stack[2*top+2] = pt0
			point_stack[2*top+1] = mid01
			point_stack[2*top+0] = midpoint(mid01, mid12)
			point_stack[2*top-1] = mid12

			continue
		} else {
			normal_pt1 := pt2.Sub(pt0).Normalize(s.half_width).Rotate(-90)
			normal_pt2 = v12.Normalize(s.half_width).Rotate(-90)
			s.adder.Add2(pt1.Add(normal_pt1), pt2.Add(normal_pt2))
			s.path.Add2(pt1.Sub(normal_pt1), pt2.Sub(normal_pt2))
		}

		if top == 0 {
			s.recent_pt, s.normal_pt = pt2, normal_pt2
			return
		}

		top--

		normal_pt0 = normal_pt2
	}
	panic("FONT unreachable.")
}

func (s *stroke_state_t) Add1(pt RastPoint) {
	normal_pt := pt.Sub(s.recent_pt).Normalize(s.half_width).Rotate(-90)
	if len(s.path) == 0 {
		s.adder.Start(s.recent_pt.Add(normal_pt))
		s.path.Start(s.recent_pt.Sub(normal_pt))
	} else {
		s.joiner.Join(s.adder, &s.path, s.half_width, s.recent_pt, s.normal_pt, normal_pt)
	}
	s.adder.Add1(pt.Add(normal_pt))
	s.path.Add1(pt.Sub(normal_pt))
	s.recent_pt, s.normal_pt = pt, normal_pt
}

func (s *stroke_state_t) Add2(pt1, pt2 RastPoint) {
	v01 := pt1.Sub(s.recent_pt)
	v12 := pt2.Sub(pt1)
	norm01 := v01.Normalize(s.half_width).Rotate(-90)

	if len(s.path) == 0 {
		s.adder.Start(s.recent_pt.Add(norm01))
		s.path.Start(s.recent_pt.Sub(norm01))
	} else {
		s.joiner.Join(s.adder, &s.path, s.half_width, s.recent_pt, s.normal_pt, norm01)
	}

	is_01_small := v01.Dot(v01) < kEpsilon
	is_12_small := v12.Dot(v12) < kEpsilon

	if is_01_small || is_12_small {
		norm02 := pt2.Sub(s.recent_pt).Normalize(s.half_width).Rotate(-90)
		s.adder.Add1(pt2.Add(norm02))
		s.path.Add1(pt2.Sub(norm02))
		s.recent_pt, s.normal_pt = pt2, norm02
		return
	}

	top := curviest2(s.recent_pt, pt1, pt2)
	if top <= 0 || top >= 65536 {
		s.add_non_curvy2(pt1, pt2)
		return
	}

	mid01 := interpolate(s.recent_pt, pt1, top)
	mid12 := interpolate(pt1, pt2, top)
	mid012 := interpolate(mid01, mid12, top)

	norm12 := v12.Normalize(s.half_width).Rotate(-90)
	if norm01.Dot(norm12) < -Fix64(s.half_width)*Fix64(s.half_width)*2047/2048 {
		arc := norm01.Dot(v12) < 0

		s.adder.Add1(mid012.Add(norm01))
		if arc {
			ptz := norm01.Rotate(90)
			add_arc(s.adder, mid012, norm01, ptz)
			add_arc(s.adder, mid012, ptz, norm12)
		}

		s.adder.Add1(mid012.Sub(norm12))
		s.adder.Add1(pt2.Sub(norm12))

		s.recent_pt, s.normal_pt = pt2, norm12
		return
	}

	s.add_non_curvy2(mid01, mid012)
	s.add_non_curvy2(mid12, pt2)
}

func (s *stroke_state_t) Add3(pt2, pt3, pt4 RastPoint) {
	panic("font raster stroke_state_t NOT IMPLEMENT for cubic segments.")
}

func (s *stroke_state_t) stroke(path Path) {
	s.path = Path(make([]Fix32, 0, len(path)))
	s.recent_pt = RastPoint{path[1], path[2]}
	for i := 4; i < len(path); {
		switch path[i] {
		case 1:
			s.Add1(RastPoint{path[i+1], path[i+2]})
			i += 4
		case 2:
			pt0 := RastPoint{path[i+1], path[i+2]}
			pt1 := RastPoint{path[i+3], path[i+4]}
			s.Add2(pt0, pt1)
			i += 6
		case 3:
			pt0 := RastPoint{path[i+1], path[i+2]}
			pt1 := RastPoint{path[i+3], path[i+4]}
			pt2 := RastPoint{path[i+4], path[i+5]}
			s.Add3(pt0, pt1, pt2)
			i += 8
		default:
			panic("FT raster bad path")
		}
	}

	if len(s.path) == 0 {
		return
	}

	s.capper.Cap(s.adder, s.half_width, path.last_point(), s.normal_pt.Neg())
	reverse_add_path(s.adder, s.path)
	pivot := path.first_point()
	s.capper.Cap(s.adder, s.half_width, pivot, pivot.Sub(RastPoint{s.path[1], s.path[2]}))
}

func Stroke(adder Adder, path Path, width Fix32, capper Capper, joiner Joiner) {
	if len(path) == 0 {
		return
	}

	if capper == nil {
		capper = RoundCapper
	}

	if joiner == nil {
		joiner = RoundJoiner
	}

	if path[0] != 0 {
		panic("FONT raster bad path")
	}

	s := stroke_state_t{
		adder:      adder,
		half_width: width,
		capper:     capper,
		joiner:     joiner,
	}
	i := 0

	for j := 4; j < len(path); {
		switch path[j] {
		case 0:
			s.stroke(path[i:j])
			i, j = j, j+4
		case 1:
			j += 4
		case 2:
			j += 6
		case 3:
			j += 8
		default:
			panic("FONT raster bad path.")
		}
	}
	s.stroke(path[i:])
}
