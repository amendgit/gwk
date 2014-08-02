// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

// http://projects.tuxee.net/cl-vectors/section-the-cl-aa-algorithm

import (
	// "log"
	"strconv"
)

type cell_t struct {
	xi    int
	area  int
	cover int
	next  int
}

type Rast struct {
	UseNonZeroWinding bool

	Dx, Dy                       int
	width                        int
	split_scale_2, split_scale_3 int

	curr_pen_pos RastPoint
	xi, yi       int

	area, cover int

	cell_array     []cell_t
	cell_idx_array []int

	cell_buf     [256]cell_t
	cell_idx_buf [64]int
	span_buf     [64]Span
}

func (r *Rast) find_cell() int {
	if r.yi < 0 || r.yi >= len(r.cell_idx_array) {
		return -1
	}

	xi := r.xi
	if xi < 0 {
		xi = -1
	} else if xi > r.width {
		xi = r.width
	}

	idx, prev := r.cell_idx_array[r.yi], -1
	for idx != -1 && r.cell_array[idx].xi <= xi {
		if r.cell_array[idx].xi == xi {
			return idx
		}
		idx, prev = r.cell_array[idx].next, idx
	}

	length := len(r.cell_array)
	if length == cap(r.cell_array) {
		tmp := make([]cell_t, length, 4*length)
		copy(tmp, r.cell_array)
		r.cell_array = tmp[0 : length+1]
	} else {
		r.cell_array = r.cell_array[0 : length+1]
	}

	r.cell_array[length] = cell_t{xi, 0, 0, idx}
	if prev == -1 {
		r.cell_idx_array[r.yi] = length
	} else {
		r.cell_array[prev].next = length
	}

	return length
}

func (r *Rast) save_cell() {
	if r.area != 0 || r.cover != 0 {
		idx := r.find_cell()
		if idx != -1 {
			r.cell_array[idx].area += r.area
			r.cell_array[idx].cover += r.cover
		}
		r.area = 0
		r.cover = 0
	}
}

func (r *Rast) set_cell(xi, yi int) {
	if r.xi != xi || r.yi != yi {
		r.save_cell()
		r.xi, r.yi = xi, yi
	}
}

// scan accumulates area/coverage for the yi'th scanline, going from
// x0 to x1 in the horizontal direction (in 24.8 fixed point co-ordinates)
// and from y0f to y1f fractional vertical units within that scanline.
func (r *Rast) scan(yi int, x0, y0f, x1, y1f Fix32) {
	// Break the 24.8 fixed point X co-ordinates into integral and fractional
	// parts.
	x0i := int(x0) / 256
	x0f := x0 - Fix32(256*x0i)
	x1i := int(x1) / 256
	x1f := x1 - Fix32(256*x1i)

	// A perfectly horizontal scan.
	if y0f == y1f {
		r.set_cell(x1i, yi)
		return
	}

	dx, dy := x1-x0, y1f-y0f
	// A single cell scan.
	if x0i == x1i {
		r.area += int((x0f + x1f) * dy)
		r.cover += int(dy)
		return
	}

	// There are at least two cells. Apart from the first and last cells,
	// all intermediate cells go through the full width of the cell,
	// or 256 units in 24.8 fixed point format.
	var p, q, edge0, edge1 Fix32
	var xi_delta int

	if dx > 0 {
		p, q = (256-x0f)*dy, dx
		edge0, edge1, xi_delta = 0, 256, 1
	} else {
		p, q = x0f*dy, -dx
		edge0, edge1, xi_delta = 256, 0, -1
	}

	y_delta, y_rem := p/q, p%q
	if y_rem < 0 {
		y_delta -= 1
		y_rem += q
	}

	// Do the first cell.
	xi, y := x0i, y0f
	r.area += int((x0f + edge1) * y_delta)
	r.cover += int(y_delta)
	xi, y = xi+xi_delta, y+y_delta
	r.set_cell(xi, yi)

	// Do all the intermediate cells.
	if xi != x1i {
		p = 256 * (y1f - y + y_delta)
		full_delta, full_rem := p/q, p%q

		if full_rem < 0 {
			full_delta -= 1
			full_rem += q
		}

		y_rem -= q

		for xi != x1i {
			y_delta = full_delta
			y_rem += full_rem

			if y_rem >= 0 {
				y_delta += 1
				y_rem -= q
			}

			r.area += int(256 * y_delta)
			r.cover += int(y_delta)
			xi, y = xi+xi_delta, y+y_delta
			r.set_cell(xi, yi)
		}
	}

	// Do the last cell.
	y_delta = y1f - y
	r.area += int((edge0 + x1f) * y_delta)
	r.cover += int(y_delta)
}

func (r *Rast) Start(pt RastPoint) {
	r.set_cell(int(pt.X/256), int(pt.Y/256))
	r.curr_pen_pos = pt
}

// Add1 adds a linear segment to the current curve.
func (r *Rast) Add1(pt RastPoint) {
	x0, y0 := r.curr_pen_pos.X, r.curr_pen_pos.Y
	x1, y1 := pt.X, pt.Y
	dx, dy := x1-x0, y1-y0
	// Break the 24.8 fixed point Y co-ordinates into integral and fractional parts.
	y0i := int(y0) / 256
	y0f := y0 - Fix32(256*y0i)
	y1i := int(y1) / 256
	y1f := y1 - Fix32(256*y1i)

	if y0i == y1i {
		// There is only one scanline.
		r.scan(y0i, x0, y0f, x1, y1f)
	} else if dx == 0 {
		// This is a vertical line segment. We avoid calling r.scan and instead
		// manipulate r.area and r.cover directly.
		var edge0, edge1 Fix32
		var yi_delta int

		if dy > 0 {
			edge0, edge1, yi_delta = 0, 256, 1
		} else {
			edge0, edge1, yi_delta = 256, 0, -1
		}

		x0i, yi := int(x0)/256, y0i
		x0f_times_2 := (int(x0) - (256 * x0i)) * 2
		// Do the first pixel.
		dcover := int(edge1 - y0f)
		darea := int(x0f_times_2 * dcover)
		r.area += darea
		r.cover += dcover
		yi += yi_delta
		r.set_cell(x0i, yi)
		// Do all the intermediate pixels.
		dcover = int(edge1 - edge0)
		darea = int(x0f_times_2 * dcover)

		for yi != y1i {
			r.area += darea
			r.cover += dcover
			yi += yi_delta
			r.set_cell(x0i, yi)
		}
		// Do the last pixel.
		dcover = int(y1f - edge0)
		darea = int(x0f_times_2 * dcover)
		r.area += darea
		r.cover += dcover
	} else {
		// There are at least two scanlines. Apart from the first and last
		// scanlines, all intermediate scanlines go through the full height of
		// the row, or 256 units in 24.8 fixed point format.
		var p, q, edge0, edge1 Fix32
		var yi_delta int

		if dy > 0 {
			p, q = (256-y0f)*dx, dy
			edge0, edge1, yi_delta = 0, 256, 1
		} else {
			p, q = y0f*dx, -dy
			edge0, edge1, yi_delta = 256, 0, -1
		}

		x_delta, x_rem := p/q, p%q
		if x_rem < 0 {
			x_delta -= 1
			x_rem += q
		}
		// Do the first scanline.
		x, yi := x0, y0i
		r.scan(yi, x, y0f, x+x_delta, edge1)
		x, yi = x+x_delta, yi+yi_delta
		r.set_cell(int(x)/256, yi)
		if yi != y1i {
			// Do all the intermediate scanlines.
			p = 256 * dx
			full_delta, full_rem := p/q, p%q
			if full_rem < 0 {
				full_delta -= 1
				full_rem += q
			}

			x_rem -= q
			for yi != y1i {
				x_delta = full_delta
				x_rem += full_rem
				if x_rem >= 0 {
					x_delta += 1
					x_rem -= q
				}
				r.scan(yi, x, edge0, x+x_delta, edge1)
				x, yi = x+x_delta, yi+yi_delta
				r.set_cell(int(x)/256, yi)
			}
		}
		// Do the last scanline.
		r.scan(yi, x, edge0, x1, y1f)
	}
	// The next lineTo starts from b.
	r.curr_pen_pos = pt
}

// Add2 adds a quadratic segment to the current curve.
func (r *Rast) Add2(pt0, pt1 RastPoint) {
	// Calculate nSplit (the number of recursive decompositions) based on how `curvy' it is.
	// Specifically, how much the middle point b deviates from (a+c)/2.
	pos := r.curr_pen_pos
	dev := max_abs(pos.X-2*pt0.X+pt1.X, pos.Y-2*pt0.Y+pt1.Y) /
		Fix32(r.split_scale_2)
	split_num := 0
	for dev > 0 {
		dev /= 4
		split_num++
	}
	// dev is 32-bit, and nsplit++ every time we shift off 2 bits, so maxNsplit
	// is 16.
	const kMaxSplitNum = 16
	if split_num > kMaxSplitNum {
		panic("FONT raster Add2 split num too large" + strconv.Itoa(split_num))
	}
	// Recursively decompose the curve nSplit levels deep.
	var point_stack [2*kMaxSplitNum + 3]RastPoint
	var split_stack [kMaxSplitNum + 1]int
	var idx int

	split_stack[0] = split_num
	point_stack[0] = pt1
	point_stack[1] = pt0
	point_stack[2] = r.curr_pen_pos

	for idx >= 0 {
		st := split_stack[idx]
		pt := point_stack[2*idx:]
		if st > 0 {
			// Split the quadratic curve p[0:3] into an equivalent set of two
			// shorter curves: p[0:3] and p[2:5]. The new p[4] is the old p[2],
			// and p[0] is unchanged.
			x, y := pt[1].X, pt[1].Y
			pt[4].X, pt[4].Y = pt[2].X, pt[2].Y
			pt[3].X, pt[3].Y = (pt[4].X+x)/2, (pt[4].Y+y)/2
			pt[1].X, pt[1].Y = (pt[0].X+x)/2, (pt[0].Y+y)/2
			pt[2].X, pt[2].Y = (pt[1].X+pt[3].X)/2, (pt[1].Y+pt[3].Y)/2
			// The two shorter curves have one less split to do.
			split_stack[idx] = st - 1
			idx++
			split_stack[idx] = st - 1
		} else {
			// Replace the level-0 quadratic with a two-linear-piece
			// approximation.
			mid_x := (pt[0].X + 2*pt[1].X + pt[2].X) / 4
			mid_y := (pt[0].Y + 2*pt[1].Y + pt[2].Y) / 4
			// log.Printf("AddQuad %v %v", mid_x, mid_y)
			r.Add1(RastPoint{mid_x, mid_y})
			// log.Printf("AddQuad %v %v", pt[0].X, pt[0].Y)
			r.Add1(pt[0])
			idx--
		}
	}
}

func (r *Rast) Add3(pt0, pt1, pt2 RastPoint) {
	pos := r.curr_pen_pos

	dev_2 := max_abs(pos.X-3*(pt0.X+pt1.X)+pt2.X, pos.Y-3*(pt0.Y+pt1.Y)+pt2.Y) /
		Fix32(r.split_scale_2)
	dev_3 := max_abs(pos.X-2*pt0.X+pt1.X, pos.Y-2*pt0.Y+pt1.Y) /
		Fix32(r.split_scale_3)

	split_num := 0
	for dev_2 > 0 || dev_3 > 0 {
		dev_2 /= 8
		dev_3 /= 4
		split_num++
	}

	const kMaxSplitNum = 16
	if split_num > kMaxSplitNum {
		panic("FONT raster Add3 split num too large.")
	}

	var point_stack [3*kMaxSplitNum + 4]RastPoint
	var split_stack [kMaxSplitNum + 1]int
	var idx int

	split_stack[0] = split_num
	point_stack[0] = pt2
	point_stack[1] = pt1
	point_stack[2] = pt0
	point_stack[3] = r.curr_pen_pos

	for idx >= 0 {
		st := split_stack[idx]
		pt := point_stack[3*idx:]
		if st > 0 {
			x01, y01 := (pt[0].X+pt[1].X)/2, (pt[0].Y+pt[1].Y)/2
			x12, y12 := (pt[1].X+pt[2].X)/2, (pt[1].Y+pt[2].Y)/2
			x23, y23 := (pt[2].X+pt[3].X)/2, (pt[2].X+pt[3].X)/2
			pt[6].X, pt[6].Y = pt[3].X, pt[3].Y
			pt[5].X, pt[5].Y = x23, y23
			pt[1].X, pt[1].Y = x01, y01
			pt[2].X, pt[2].Y = (x01+x12)/2, (y01+y12)/2
			pt[4].X, pt[4].Y = (x12+x23)/2, (y12+y23)/2
			pt[3].X, pt[3].Y = (x12+x23)/2, (y12+y23)/2

			split_stack[idx] = st - 1
			idx++
			split_stack[idx] = st - 1
		} else {
			//mid_x := (pt[0].X + 3*(pt[1].X+pt[2].X) + pt[3].X) / 8
			//mid_y := (pt[1].Y + 3*(pt[1].Y+pt[2].Y) + pt[3].Y) / 8
			//r.Add1(RastPoint{mid_x, mid_y})
			//r.Add1(pt[0])
			idx--
		}
	}
}

func (r *Rast) AddPath(path Path) {
	for i := 0; i < len(path); {
		switch path[i] {
		case 0:
			pt := RastPoint{path[i+1], path[i+2]}
			r.Start(pt)
		case 1:
			pt := RastPoint{path[i+1], path[i+2]}
			r.Add1(pt)
			i += 4
		case 2:
			pt0 := RastPoint{path[i+1], path[i+2]}
			pt1 := RastPoint{path[i+2], path[i+3]}
			r.Add2(pt0, pt1)
			i += 6
		case 3:
			pt0 := RastPoint{path[i+1], path[i+2]}
			pt1 := RastPoint{path[i+2], path[i+3]}
			pt2 := RastPoint{path[i+5], path[i+6]}
			r.Add3(pt0, pt1, pt2)
			i += 8
		default:
			panic("FONT raster bad path")
		}
	}
}

func (r *Rast) AddStroke(p Path, width Fix32, capper Capper, joiner Joiner) {
	Stroke(r, p, width, capper, joiner)
}

func (r *Rast) area_to_alpha(area int) uint32 {
	//log.Printf("area_to_alpha")
	tmp := (area + 1) >> 1
	if tmp < 0 {
		tmp = -tmp
	}
	alpha := uint32(tmp)

	if r.UseNonZeroWinding {
		if alpha > 0xffff {
			alpha = 0xffff
		}
	} else {
		alpha &= 0x1ffff

		if alpha > 0x10000 {
			alpha = 0x20000 - alpha
		} else if alpha == 0x10000 {
			alpha = 0x0ffff
		}
	}

	alpha |= alpha << 16

	return alpha
}

// Rastize converts r's accumulated curves into Spans for p. The Spans
// passed to p are non-overlapping, and sorted by Y and then X. They all
// have non-zero width (and 0 <= X0 < X1 <= r.width) and non-zero A, except
// for the final Span, which has Y, X0, X1 and A all equal to zero.
func (r *Rast) Rast(draw Drawer) {
	r.save_cell()
	span_idx := 0
	for yi := 0; yi < len(r.cell_idx_array); yi++ {
		xi, cover := 0, 0
		for idx := r.cell_idx_array[yi]; idx != -1; idx = r.cell_array[idx].next {
			if cover != 0 && r.cell_array[idx].xi > xi {
				alpha := r.area_to_alpha(cover * 256 * 2)
				if alpha != 0 {
					xi0, xi1 := xi, r.cell_array[idx].xi
					if xi0 < 0 {
						xi0 = 0
					}
					if xi1 >= r.width {
						xi1 = r.width
					}
					if xi0 < xi1 {
						r.span_buf[span_idx] =
							Span{yi + r.Dy, xi0 + r.Dx, xi1 + r.Dx, alpha}
						span_idx++
					}
				}
			}
			cover += r.cell_array[idx].cover
			alpha := r.area_to_alpha(cover*256*2 - r.cell_array[idx].area)
			xi = r.cell_array[idx].xi + 1
			if alpha != 0 {
				xi0, xi1 := r.cell_array[idx].xi, xi
				if xi0 < 0 {
					xi0 = 0
				}
				if xi1 > r.width {
					xi1 = r.width
				}
				if xi0 < xi1 {
					r.span_buf[span_idx] =
						Span{yi + r.Dy, xi0 + r.Dx, xi1 + r.Dx, alpha}
					span_idx++
				}
			}
			if span_idx > len(r.span_buf)-2 {
				draw.Draw(r.span_buf[0:span_idx], false)
				span_idx = 0
			}
		}
	}
	draw.Draw(r.span_buf[0:span_idx], true)
}

func (r *Rast) Clear() {
	r.curr_pen_pos = RastPoint{0, 0}
	r.xi = 0
	r.yi = 0
	r.area = 0
	r.cover = 0
	r.cell_array = r.cell_array[0:0]
	for i := 0; i < len(r.cell_idx_array); i++ {
		r.cell_idx_array[i] = -1
	}
}

func (r *Rast) SetBounds(width, height int) {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}

	// Use the same ssN heuristic as the C Freetype implementation.
	// The C implementation uses the vales 32, 16, but those are in
	// 26.6 fixed point units, and use 24.8 fixed point everywhere.
	ss2, ss3 := 128, 64
	if width > 24 || height > 24 {
		ss2, ss3 = 2*ss2, 2*ss3
		if width > 120 || height > 120 {
			ss2, ss3 = 2*ss2, 2*ss3
		}
	}
	r.width = width
	r.split_scale_2 = ss2
	r.split_scale_3 = ss3
	r.cell_array = r.cell_buf[0:0]
	if height > len(r.cell_idx_buf) {
		r.cell_idx_array = make([]int, height)
	} else {
		r.cell_idx_array = r.cell_idx_buf[0:height]
	}

	r.Clear()
}

func NewRast(width, height int) *Rast {
	r := new(Rast)
	r.SetBounds(width, height)
	return r
}
