// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

import "math"

type f26d6_t int32
type f2d14_t int16

func normal(a, b f2d14_t) [2]f2d14_t {
	f0, f1 := float64(a), float64(b)
	l := 0x4000 / math.Hypot(f0, f1)

	f0 = f0 * l
	if f0 >= 0 {
		f0 = f0 + 0.5
	} else {
		f0 = f0 - 0.5
	}

	f1 = f1 * l
	if f1 >= 0 {
		f1 = f1 + 0.5
	} else {
		f1 = f1 - 0.5
	}

	return [2]f2d14_t{f2d14_t(f0), f2d14_t(f1)}
}

func f26d6_abs(f f26d6_t) f26d6_t {
	if f < 0 {
		f = -f
	}
	return f
}

func f26d6_div(a, b f26d6_t) f26d6_t {
	return f26d6_t((int64(a) << 6) / int64(b))
}

func f26d6_mul(a, b f26d6_t) f26d6_t {
	return f26d6_t(int64(a) * int64(b) >> 6)
}

func dot_X(a, b f26d6_t, q [2]f2d14_t) f26d6_t {
	x0 := int64(a)
	y0 := int64(b)
	x1 := int64(q[0])
	y1 := int64(q[1])
	return f26d6_t((x0*x1 + y0*y1 + 1<<13) >> 14)
}
