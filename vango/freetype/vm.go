// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

import (
	"errors"
	// "log"
	// "os"
	// "runtime/debug"
)

const (
	kZoneTwilight = 0
	kZoneGlyph    = 1
	kZoneNum      = 2
)

type point_type_t uint32

const (
	kPointTypeExec   point_type_t = 0 // Point that execing.
	kPointTypeUnexec point_type_t = 1 // Point that not execed.
	kPointTypeRaw    point_type_t = 2 // Point that not execed and scaled.
	kPointTypeNum                 = 3
)

type call_entry_t struct {
	pgm        []byte
	pgm_count  int
	loop_count int32
}

// Execute the program bytes in font.
type exec_t struct {
	stack []int32
	store []int32

	// A map from function number to function byte code.
	id_func_map map[int32][]byte

	font  *Font
	scale int32

	graphic_state         graphic_state_t
	default_graphic_state graphic_state_t

	// |points| and |ends| are the twilight zone's points, glyph's points
	// and glyph's contour boundaries.
	points [kZoneNum][kPointTypeNum][]FontPoint
	ends   []int

	// scaled_cvt is the lazily initialized scaled control value table.
	is_scaled_cvt_init bool
	scaled_cvt         []f26d6_t
}

// https://developer.apple.com/fonts/TTRefMan/RM04/Chap4.html
type graphic_state_t struct {
	projection_vector [2]f2d14_t
	freedom_vector    [2]f2d14_t
	dual_vector       [2]f2d14_t

	// reference ponit and zone points.
	rp0, rp1, rp2 int32
	zp0, zp1, zp2 int32

	ctrl_val_cut_in     f26d6_t
	single_width_cut_in f26d6_t
	single_width        f26d6_t

	delta_base  int32
	delta_shift int32

	min_dist f26d6_t

	loop_count int32

	round_period    f26d6_t
	round_phase     f26d6_t
	round_threshold f26d6_t

	auto_flip bool
}

var g_defualt_graphic_state = graphic_state_t{
	projection_vector: [2]f2d14_t{0x4000, 0},
	freedom_vector:    [2]f2d14_t{0x4000, 0},
	dual_vector:       [2]f2d14_t{0x4000, 0},
	zp0:               1,
	zp1:               1,
	zp2:               1,
	ctrl_val_cut_in:   (17 << 6) / 16, // 17/16 as an f26d6_t
	delta_base:        9,
	delta_shift:       3,
	min_dist:          1 << 6, // 1 as an f26d6_t
	loop_count:        1,
	round_period:      1 << 6,
	round_threshold:   1 << 5, // 1/2 as an f26d6_t
	auto_flip:         true,
}

func (exec *exec_t) init(font *Font, scale int32) error {
	reserve := func(point_array []FontPoint) []FontPoint {
		if num := int(font.max_twilight_points) + 4; num <= cap(point_array) {
			point_array = point_array[:num]
			for i := range point_array {
				point_array[i] = FontPoint{}
			}
		} else {
			point_array = make([]FontPoint, num)
		}
		return point_array
	}

	exec.points[kZoneTwilight][0] = reserve(exec.points[kZoneTwilight][0])
	exec.points[kZoneTwilight][1] = reserve(exec.points[kZoneTwilight][1])
	exec.points[kZoneTwilight][2] = reserve(exec.points[kZoneTwilight][2])

	scale_changed := exec.scale != scale

	if exec.font != font {
		exec.font, scale_changed = font, true

		if exec.id_func_map == nil {
			exec.id_func_map = make(map[int32][]byte)
		} else {
			for id := range exec.id_func_map {
				delete(exec.id_func_map, id)
			}
		}

		if num := int(font.max_stack_elements); num > len(exec.stack) {
			num = num + 255
			num = num &^ 255
			exec.stack = make([]int32, num)
		}

		if num := int(font.max_storage); num > len(exec.store) {
			num = num + 15
			num = num &^ 15
			exec.store = make([]int32, num)
		}

		if len(font.fpgm) != 0 {
			if err := exec.exec(font.fpgm, nil, nil, nil, nil); err != nil {
				return err
			}
		}
	}

	if scale_changed {
		exec.scale = scale
		exec.is_scaled_cvt_init = false
		exec.default_graphic_state = g_defualt_graphic_state

		if len(font.prep) != 0 {
			if err := exec.exec(font.prep, nil, nil, nil, nil); err != nil {
				return err
			}

			exec.default_graphic_state = exec.graphic_state

			// The MS rasterizer doesn't allow the following graphics state
			// variables to be modified by the CVT program.
			gs0, gs1 := &(exec.default_graphic_state), &(g_defualt_graphic_state)
			gs0.projection_vector = gs1.projection_vector
			gs0.freedom_vector = gs1.freedom_vector
			gs0.dual_vector = gs1.dual_vector
			gs0.rp0 = gs1.rp0
			gs0.rp1 = gs1.rp1
			gs0.rp2 = gs1.rp2
			gs0.zp0 = gs1.zp0
			gs0.zp1 = gs1.zp1
			gs0.zp2 = gs1.zp2
			gs0.loop_count = gs1.loop_count
		}
	}

	return nil
}

func (exec *exec_t) exec(pgm []byte, cur_point_array []FontPoint,
	unhint_point_array []FontPoint,
	inunit_point_array []FontPoint, ends []int) error {
	exec.graphic_state = exec.default_graphic_state
	exec.points[kZoneGlyph][kPointTypeExec] = cur_point_array
	exec.points[kZoneGlyph][kPointTypeUnexec] = unhint_point_array
	exec.points[kZoneGlyph][kPointTypeRaw] = inunit_point_array
	exec.ends = ends
	// log.Printf("pgm %v", pgm)
	// debug.PrintStack()
	if len(pgm) > 50000 {
		return errors.New("Freetype: pgm too many instructions.")
	}

	var call_stack [32]call_entry_t
	var call_stack_top int
	var pgm_count int
	var opcode byte
	var top int

	skip_instruction_playload := func(pgm []byte, pgm_count int) (int, bool) {
		switch pgm[pgm_count] {
		case kOpNPUSHB:
			pgm_count++
			if pgm_count >= len(pgm) {
				return 0, false
			}
			pgm_count += int(pgm[pgm_count])
		case kOpNPUSHW:
			pgm_count++
			if pgm_count >= len(pgm) {
				return 0, false
			}
			pgm_count += 2 * int(pgm[pgm_count])
		case kOpPUSHB000, kOpPUSHB001, kOpPUSHB010, kOpPUSHB011, kOpPUSHB100,
			kOpPUSHB101, kOpPUSHB110, kOpPUSHB111:
			pgm_count += int(pgm[pgm_count] - (kOpPUSHB000 - 1))
		case kOpPUSHW000, kOpPUSHW001, kOpPUSHW010, kOpPUSHW011, kOpPUSHW100,
			kOpPUSHW101, kOpPUSHW110, kOpPUSHW111:
			pgm_count += 2 * int(pgm[pgm_count]-(kOpPUSHW000-1))
		}
		return pgm_count, true
	}

	ifelse := func() error {
		// Skip past bytecode until the next ELSE (if opcode == 0) or the
		// next EIF (for all opcodes). Opcode == 0 means that we have come
		// from an IF. Opcode == 1 means taht we have come from an ELSE.
	loop:
		for depth := 0; ; {
			pgm_count++
			if pgm_count >= len(pgm) {
				return errors.New("Freetype: exec unbalanced IF or ELSE.")
			}
			switch pgm[pgm_count] {
			case kOpIF:
				depth++
			case kOpELSE:
				if depth == 0 && opcode == 0 {
					break loop
				}
			case kOpEIF:
				depth--
				if depth < 0 {
					break loop
				}
			default:
				ok := false
				pgm_count, ok = skip_instruction_playload(pgm, pgm_count)
				if !ok {
					return errors.New("Freetype: exec unbalanced IF or ELSE")
				}
			}
		}
		pgm_count++

		return nil
	}

	push := func() error {
		// Push n elements from the program to the stack, where n is low 7 bits
		// of opcode. If the low 7 bits are zero, then n is the next byte from
		// The high bit being 0 means that elements are zero-exetended bytes.
		// The high bit being 1 means that the elements are sign-extended words.
		width := 1
		if opcode&0x80 != 0 {
			opcode &^= 0x80
			width = 2
		}

		if opcode == 0 {
			pgm_count++
			if pgm_count >= len(pgm) {
				return errors.New("Freetype: exec insufficient data.")
			}
			opcode = pgm[pgm_count]
		}

		pgm_count++

		if top+width*int(opcode) > len(exec.stack) {
			return errors.New("Freetype: stack overflow.")
		}

		if pgm_count+width*int(opcode) > len(pgm) {
			return errors.New("Freetype: exec insufficient data.")
		}
		for ; opcode > 0; opcode-- {
			if width == 1 {
				exec.stack[top] = int32(pgm[pgm_count])
			} else {
				exec.stack[top] =
					int32(int8(pgm[pgm_count]))<<8 | int32(pgm[pgm_count+1])
			}
			top = top + 1
			pgm_count = pgm_count + width
		}

		return nil
	}

	deltap := func() error {
		num := f26d6_t(exec.stack[top-1])
		top = top - 1

		gs := &(exec.graphic_state)

		if top < 2*int(gs.loop_count) {
			return errors.New("Freetype: exec stack overflow.")
		}

		for ; num > 0; num-- {
			arg := exec.stack[top-1]
			top = top - 1

			pt := exec.point_at(gs.zp0, kPointTypeExec, arg)
			if pt == nil {
				return errors.New("Freetype: exec point out of range.")
			}

			lo := exec.stack[top-1]
			top = top - 1

			hi := (lo & 0xf0) >> 4
			switch opcode {
			case kOpDELTAP2:
				hi = hi + 16
			case kOpDELTAP3:
				hi = hi + 32
			}
			hi = hi + gs.delta_base
			if pixel_per_em := (exec.scale + 1<<5) >> 6; pixel_per_em != hi {
				continue
			}

			lo = (lo & 0x0f) - 8
			if lo >= 0 {
				lo++
			}
			lo = lo * 64 / (1 << uint32(gs.delta_shift))
			exec.move_pt(pt, f26d6_t(lo), true)
		}

		pgm_count++

		return nil
	}

	var step_count = 0

	// fd, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// defer fd.Close()
	// log.SetOutput(fd)
	// log.Printf("pgm %v", pgm)
	for 0 <= pgm_count && pgm_count < len(pgm) {
		opcode = pgm[pgm_count]
		// log.Printf("stack %v", exec.stack[0:top])
		// // log.Printf("loop %v", exec.graphic_state.loop_count)
		// log.Printf("Points %v", exec.points[kZoneGlyph][kPointTypeExec])
		// log.Printf("pgm_count %v", pgm_count)
		// log.Printf("opcode %x", opcode)

		step_count++

		if step_count == 10000 {
			return errors.New("Freetype: too many steps.")
		}

		if kPopNum[opcode] == kOpN {
			return errors.New("Freetype: pgm unimplement instruction.")
		}

		if top < int(kPopNum[opcode]) {
			return errors.New("Freetype: pgm stack underflow.")
		}

		switch opcode {
		case kOpSVTCA0:
			gs := &(exec.graphic_state)
			gs.projection_vector = [2]f2d14_t{0, 0x4000}
			gs.freedom_vector = [2]f2d14_t{0, 0x4000}
			gs.dual_vector = [2]f2d14_t{0, 0x4000}

		case kOpSVTCA1:
			gs := &(exec.graphic_state)
			gs.projection_vector = [2]f2d14_t{0x4000, 0}
			gs.freedom_vector = [2]f2d14_t{0x4000, 0}
			gs.dual_vector = [2]f2d14_t{0x4000, 0}

		case kOpSPVTCA0:
			gs := &(exec.graphic_state)
			gs.projection_vector = [2]f2d14_t{0, 0x4000}
			gs.dual_vector = [2]f2d14_t{0, 0x4000}

		case kOpSPVTCA1:
			gs := &(exec.graphic_state)
			gs.projection_vector = [2]f2d14_t{0x4000, 0}
			gs.dual_vector = [2]f2d14_t{0x4000, 0}

		case kOpSFVTCA0:
			exec.graphic_state.freedom_vector = [2]f2d14_t{0, 0x4000}

		case kOpSFVTCA1:
			exec.graphic_state.freedom_vector = [2]f2d14_t{0x4000, 0}

		case kOpSPVTL0, kOpSPVTL1, kOpSFVTL0, kOpSFVTL1:
			// log.Printf("!!!!! %v", exec.stack[0:top])
			gs := &(exec.graphic_state)
			p0 := exec.point_at(gs.zp0, kPointTypeExec, exec.stack[top-2])
			p1 := exec.point_at(gs.zp0, kPointTypeExec, exec.stack[top-1])
			top = top - 2

			if p0 == nil || p1 == nil {
				return errors.New("Freetype: exec point out of range.")
			}

			dx := f2d14_t(p0.X - p1.X)
			dy := f2d14_t(p0.Y - p1.Y)

			if dx == 0 && dy == 0 {
				dx = 0x4000
			} else if opcode&1 != 0 {
				// Counter-clockwise rotation.
				dx, dy = -dy, dx
			}

			vec := normal(dx, dy)

			if opcode < kOpSFVTL0 {
				exec.graphic_state.projection_vector = vec
				exec.graphic_state.dual_vector = vec
			} else {
				exec.graphic_state.freedom_vector = vec
			}
			// log.Printf("!!!!! %v", exec.stack[0:top])
		case kOpSPVFS:
			exec.graphic_state.projection_vector[0] = f2d14_t(exec.stack[top-2])
			exec.graphic_state.projection_vector[1] = f2d14_t(exec.stack[top-1])
			top = top - 2
			// TODO: normalize gs.pv ??
			// TODO: exec.gs.dv = exec.gs.pv ??

		case kOpSFVFS:
			exec.graphic_state.freedom_vector[0] = f2d14_t(exec.stack[top-2])
			exec.graphic_state.freedom_vector[1] = f2d14_t(exec.stack[top-1])
			top = top - 2
			// TODO: normalize exec.gs.pv ??

		case kOpGPV:
			if top >= len(exec.stack)-1 {
				return errors.New("Freetype: exec stack overflow.")
			}

			exec.stack[top] = int32(exec.graphic_state.freedom_vector[0])
			exec.stack[top+1] = int32(exec.graphic_state.freedom_vector[1])

			top = top + 2

		case kOpGFV:
			if top >= len(exec.stack)-1 {
				return errors.New("Freetype: exec stack overflow.")
			}

			exec.stack[top] = int32(exec.graphic_state.projection_vector[0])
			exec.stack[top+1] = int32(exec.graphic_state.projection_vector[1])

			top = top + 2

		case kOpSFVTPV:
			gs := &(exec.graphic_state)
			gs.freedom_vector = gs.projection_vector

		case kOpSRP0:
			exec.graphic_state.rp0 = exec.stack[top-1]
			top = top - 1

		case kOpSRP1:
			exec.graphic_state.rp1 = exec.stack[top-1]
			top = top - 1

		case kOpSRP2:
			exec.graphic_state.rp2 = exec.stack[top-1]
			top = top - 1

		case kOpSZP0:
			exec.graphic_state.zp0 = exec.stack[top-1]
			top = top - 1

		case kOpSZP1:
			exec.graphic_state.zp1 = exec.stack[top-1]
			top = top - 1

		case kOpSZP2:
			exec.graphic_state.zp2 = exec.stack[top-1]
			top = top - 1

		case kOpSZPS:
			gs := &(exec.graphic_state)
			gs.zp0 = exec.stack[top-1]
			gs.zp1 = exec.stack[top-1]
			gs.zp2 = exec.stack[top-1]
			top = top - 1

		case kOpSLOOP:
			exec.graphic_state.loop_count = exec.stack[top-1]
			top = top - 1

		case kOpRTG:
			gs := &(exec.graphic_state)
			gs.round_period = 1 << 6
			gs.round_phase = 0
			gs.round_threshold = 1 << 5

		case kOpRTHG:
			gs := &(exec.graphic_state)
			gs.round_period = 1 << 6
			gs.round_phase = 1 << 5
			gs.round_threshold = 1 << 5

		case kOpSMD:
			exec.graphic_state.min_dist = f26d6_t(exec.stack[top])

		case kOpELSE:
			opcode = 1
			if err := ifelse(); err != nil {
				return err
			}

			continue

		case kOpJMPR:
			// log.Printf("JMPR")
			pgm_count = pgm_count + int(exec.stack[top-1])
			top = top - 1

			continue

		case kOpSCVTCI:
			exec.graphic_state.ctrl_val_cut_in = f26d6_t(exec.stack[top-1])
			top = top - 1

		case kOpSSWCI:
			exec.graphic_state.single_width_cut_in = f26d6_t(exec.stack[top-1])
			top = top - 1

		case kOpSSW:
			sw := exec.font.scale(exec.scale * exec.stack[top-1])
			exec.graphic_state.single_width = f26d6_t(sw)
			top = top - 1

		case kOpDUP:
			if top >= len(exec.stack) {
				return errors.New("Freetype: exec stack overflow.")
			}
			exec.stack[top] = exec.stack[top-1]
			top = top + 1

		case kOpPOP:
			top = top - 1

		case kOpCLEAR:
			top = 0

		case kOpSWAP:
			exec.stack[top-1], exec.stack[top-2] =
				exec.stack[top-2], exec.stack[top-1]

		case kOpDEPTH:
			if top >= len(exec.stack) {
				return errors.New("Freetype: exec stack overflow.")
			}
			exec.stack[top] = int32(top)
			top = top + 1

		case kOpCINDEX, kOpMINDEX:
			offset := int(exec.stack[top-1])
			if offset <= 0 || offset >= top {
				return errors.New("Freetype: exec stack overflow.")
			}

			exec.stack[top-1] = exec.stack[top-1-offset]

			if opcode == kOpMINDEX {
				copy(exec.stack[top-1-offset:top-1], exec.stack[top-offset:top])
				top = top - 1
			}

		case kOpLOOPCALL, kOpCALL:
			if call_stack_top >= len(call_stack) {
				return errors.New("Freetype: exec call stack overflow.")
			}

			id := exec.stack[top-1]
			top = top - 1
			func_bytes, ok := exec.id_func_map[id]

			if !ok {
				return errors.New("Freetype: exec undefined function.")
			}

			call_stack[call_stack_top] = call_entry_t{pgm, pgm_count, 1}

			if opcode == kOpLOOPCALL {
				val := exec.stack[top-1]
				top = top - 1

				if val == 0 {
					break
				}

				call_stack[call_stack_top].loop_count = val
			}

			call_stack_top = call_stack_top + 1
			pgm, pgm_count = func_bytes, 0

			continue

		case kOpFDEF:
			// Save all byte code up until the next ENDF.
			begin_pgm_count := pgm_count + 1

		loop:
			for {
				pgm_count++
				if pgm_count >= len(pgm) {
					return errors.New("Freetype: exec unbalanced FDEF.")
				}

				switch pgm[pgm_count] {
				case kOpFDEF:
					return errors.New("Freetype: exec nested FDEF.")
				case kOpENDF:
					id := exec.stack[top-1]
					exec.id_func_map[id] = pgm[begin_pgm_count : pgm_count+1]
					top = top - 1
					break loop
				default:
					ok := false
					pgm_count, ok = skip_instruction_playload(pgm, pgm_count)
					if !ok {
						return errors.New("Freetype: exec unbalanced FDEF.")
					}
				}
			}

		case kOpENDF:
			if call_stack_top == 0 {
				return errors.New("Freetype: exec call stack underflow.")
			}

			call_stack_top = call_stack_top - 1
			call_stack[call_stack_top].loop_count--

			if call_stack[call_stack_top].loop_count != 0 {
				call_stack_top++
				pgm_count = 0
				continue
			}

			entry := &(call_stack[call_stack_top])
			pgm, pgm_count = entry.pgm, entry.pgm_count

		case kOpMDAP0, kOpMDAP1:
			index := exec.stack[top-1]
			top = top - 1
			gs := &(exec.graphic_state)
			pt := exec.point_at(gs.zp0, kPointTypeExec, index)
			if pt == nil {
				return errors.New("Freetype: exec point out of range.")
			}

			dist := f26d6_t(0)
			if opcode == kOpMDAP1 {
				dist = dot_X(f26d6_t(pt.X), f26d6_t(pt.Y),
					exec.graphic_state.projection_vector)
				// TODO: metrics compensation.
				dist = exec.round(dist) - dist
			}

			exec.move_pt(pt, dist, true)
			exec.graphic_state.rp0 = index
			exec.graphic_state.rp1 = index

		case kOpIUP0, kOpIUP1:

			mask := uint32(kDecodeTouchedX)
			y := false

			if opcode == kOpIUP0 {
				mask = kDecodeTouchedY
				y = true
			}

			prev := 0

			for _, end := range exec.ends {
				for i := prev; i < end; i++ {
					glyph_cur_pts := exec.points[kZoneGlyph][kPointTypeExec]
					for i < end && glyph_cur_pts[i].Flag&mask == 0 {
						i++
					}

					if i == end {
						break
					}

					first, cur := i, i
					i++

					for ; i < end; i++ {
						if glyph_cur_pts[i].Flag&mask != 0 {
							exec.iup_interp(y, cur+1, i-1, cur, i)
							cur = i
						}
					}

					if cur == first {
						exec.iup_shift(y, prev, end, cur)
					} else {
						exec.iup_interp(y, cur+1, end-1, cur, first)
						if first > 0 {
							exec.iup_interp(y, prev, first-1, cur, first)
						}
					}

				}
				prev = end
			}

		case kOpSHP0, kOpSHP1:
			if top < int(exec.graphic_state.loop_count) {
				return errors.New("Freetype: exec stack overflow.")
			}

			_, _, dist, ok := exec.displace(opcode&1 == 0)
			if !ok {
				return errors.New("Freetype: exec point out of range.")
			}

			gs := &(exec.graphic_state)
			loop := exec.graphic_state.loop_count
			for ; loop != 0; loop-- {
				pt := exec.point_at(gs.zp0, kPointTypeExec, exec.stack[top-1])
				top = top - 1
				if pt == nil {
					return errors.New("Freetype: exec point out of range.")
				}
				exec.move_pt(pt, dist, true)
			}

			gs.loop_count = 1

		case kOpSHZ0, kOpSHZ1:
			zpi, index, dist, ok := exec.displace(opcode&1 == 0)
			if !ok {
				return errors.New("Freetype: exec point out of range.")
			}

			// As per C Freetype, SHZ doesn't move_pt the phantom points, or mark
			// the points as touched.
			zp2 := exec.graphic_state.zp2
			limit := int32(len(exec.points[zp2][kPointTypeExec]))
			if exec.graphic_state.zp2 == kZoneGlyph {
				limit = limit - 4
			}

			for i := int32(0); i < limit; i++ {
				for index != i || zpi != 2 {
					pt := exec.point_at(zp2, kPointTypeExec, index)
					exec.move_pt(pt, dist, false)
				}
			}
			top = top - 1

		case kOpSHPIX:
			dist := f26d6_t(exec.stack[top-1])
			top = top - 1
			if top < int(exec.graphic_state.loop_count) {
				return errors.New("Freetype: exec stack overflow.")
			}

			loop := exec.graphic_state.loop_count
			zp2 := exec.graphic_state.zp2
			for ; loop != 0; loop-- {
				pt := exec.point_at(zp2, kPointTypeExec, exec.stack[top-1])
				top = top - 1
				if pt == nil {
					return errors.New("Freetype: exec point out of range.")
				}
				exec.move_pt(pt, dist, true)
			}

			exec.graphic_state.loop_count = 1

		case kOpIP:
			gs := &(exec.graphic_state)

			if top < int(gs.loop_count) {
				return errors.New("Freetype: exec stack overflow.")
			}

			point_type := kPointTypeRaw
			twilight := gs.zp0 == 0 || gs.zp1 == 0 || gs.zp2 == 0
			if twilight {
				point_type = kPointTypeUnexec
			}

			zp0 := gs.zp0
			zp1 := gs.zp1
			zp2 := gs.zp2

			rp1 := gs.rp1
			rp2 := gs.rp2

			pt := exec.point_at(zp1, point_type, rp2)
			old_pt := exec.point_at(zp0, point_type, rp1)
			old_range := dot_X(f26d6_t(pt.X-old_pt.X), f26d6_t(pt.Y-old_pt.Y),
				gs.dual_vector)

			pt = exec.point_at(zp1, kPointTypeExec, rp2)
			cur_pt := exec.point_at(zp0, kPointTypeExec, rp1)
			cur_range := dot_X(f26d6_t(pt.X-cur_pt.X), f26d6_t(pt.Y-cur_pt.Y),
				gs.projection_vector)

			loop := gs.loop_count
			for ; loop != 0; loop-- {
				index := exec.stack[top-1]
				top = top - 1

				pt = exec.point_at(zp2, point_type, index)
				old_dist := dot_X(f26d6_t(pt.X-old_pt.X),
					f26d6_t(pt.Y-old_pt.Y), gs.dual_vector)

				pt = exec.point_at(zp2, kPointTypeExec, index)
				cur_dist := dot_X(f26d6_t(pt.X-cur_pt.X),
					f26d6_t(pt.Y-cur_pt.Y), gs.projection_vector)

				new_dist := f26d6_t(0)

				if old_dist != 0 {
					if old_range != 0 {
						rst := mul_div(int64(old_dist), int64(cur_range),
							int64(old_range))
						new_dist = f26d6_t(rst)
					} else {
						new_dist = -old_dist
					}
				}

				exec.move_pt(pt, new_dist-cur_dist, true)
			}

			gs.loop_count = 1

		case kOpMSIRP0, kOpMSIRP1:
			index := exec.stack[top-1]
			top = top - 1
			dist := f26d6_t(exec.stack[top-1])
			top = top - 1

			gs := &(exec.graphic_state)

			// TODO: special case exec.graphic_state.zp1 == 0 in C Freetype.
			ref := exec.point_at(gs.zp0, kPointTypeExec, gs.rp0)
			pt := exec.point_at(gs.zp1, kPointTypeExec, index)

			if ref == nil || pt == nil {
				return errors.New("Freetype: exec out of range.")
			}

			cur_dist := dot_X(f26d6_t(pt.X-ref.X), f26d6_t(pt.Y-ref.Y),
				gs.dual_vector)

			// Set rp0 bit
			if opcode == kOpMSIRP1 {
				gs.rp0 = index
			}

			gs.rp1 = gs.rp0
			gs.rp2 = index

			exec.move_pt(pt, dist-cur_dist, true)

		case kOpALIGNRP:
			// log.Printf("###%v", exec.stack[0:top])
			gs := &(exec.graphic_state)

			if top < int(gs.loop_count) {
				return errors.New("Freetype: exec stack overflow.")
			}

			ref := exec.point_at(gs.zp0, kPointTypeExec, gs.rp0)
			if ref == nil {
				return errors.New("Freetype: exec point out of range.")
			}
			// log.Printf("###%v", gs.loop_count)
			loop := gs.loop_count
			for ; loop != 0; loop-- {
				pt := exec.point_at(gs.zp1, kPointTypeExec, exec.stack[top-1])
				top = top - 1
				if pt == nil {
					return errors.New("Freetype: exec point out of range.")
				}
				dist := dot_X(f26d6_t(pt.X-ref.X), f26d6_t(pt.Y-ref.Y),
					gs.projection_vector)
				exec.move_pt(pt, -dist, true)
			}

			gs.loop_count = 1
			// log.Printf("###%v", exec.stack[0:top])
		case kOpRTDG:
			gs := &(exec.graphic_state)
			gs.round_period = 1 << 5
			gs.round_phase = 0
			gs.round_threshold = 1 << 4

		case kOpMIAP0, kOpMIAP1:
			// log.Printf("!!!!!!!Points %v", exec.points[kZoneGlyph][kPointTypeExec])

			dist := exec.get_scaled_cvt(exec.stack[top-1])

			top = top - 1
			index := exec.stack[top-1]
			top = top - 1

			gs := &(exec.graphic_state)
			if gs.zp0 == 0 {
				pt0 := exec.point_at(gs.zp0, kPointTypeUnexec, index)
				pt0.X = int32(int64(dist) * int64(gs.freedom_vector[0]) >> 14)
				pt0.Y = int32(int64(dist) * int64(gs.freedom_vector[1]) >> 14)
				pt1 := exec.point_at(gs.zp0, kPointTypeExec, index)
				*pt1 = *pt0
			}
			pt := exec.point_at(gs.zp0, kPointTypeExec, index)
			old_dist := dot_X(f26d6_t(pt.X), f26d6_t(pt.Y),
				gs.projection_vector)

			if opcode == kOpMIAP1 {
				if f26d6_abs(dist-old_dist) > gs.ctrl_val_cut_in {
					dist = old_dist
				}

				// TODO: metrics compensation.
				dist = exec.round(dist)
			}

			exec.move_pt(pt, dist-old_dist, true)
			gs.rp0 = index
			gs.rp1 = index
			// log.Printf("!!!!!!!Points %v", exec.points[kZoneGlyph][kPointTypeExec])
		case kOpNPUSHB:
			// log.Printf("kOpNPUSHB")
			opcode = 0
			if err := push(); err != nil {
				return err
			}

			continue

		case kOpNPUSHW:
			opcode = 0x80
			if err := push(); err != nil {
				return err
			}

			continue

		case kOpWS:
			data := int32(exec.stack[top-1])
			index := int(exec.stack[top-2])
			top = top - 2
			if index < 0 || len(exec.store) <= index {
				return errors.New("Freetype: exec invalid data.")
			}

			exec.store[index] = data

		case kOpRS:
			index := int(exec.stack[top-1])
			if index < 0 || len(exec.store) <= index {
				return errors.New("Freetype: exec invalid data.")
			}

			exec.stack[top-1] = exec.store[index]

		case kOpWCVTP:
			exec.set_scaled_cvt(exec.stack[top-2], f26d6_t(exec.stack[top-1]))
			top = top - 2

		case kOpRCVT:
			exec.stack[top-1] = int32(exec.get_scaled_cvt(exec.stack[top-1]))
			top = top - 1

		case kOpGC0, kOpGC1:
			index := exec.stack[top-1]
			top = top - 1

			gs := &(exec.graphic_state)

			if opcode == kOpGC0 {
				pt := exec.point_at(gs.zp2, kPointTypeExec, index)
				exec.stack[top-1] = int32(dot_X(f26d6_t(pt.X),
					f26d6_t(pt.Y), gs.projection_vector))
			} else {
				pt := exec.point_at(gs.zp2, kPointTypeUnexec, index)
				exec.stack[top-1] = int32(dot_X(f26d6_t(pt.X),
					f26d6_t(pt.Y), gs.projection_vector))
			}

		case kOpMD0, kOpMD1:
			idx0 := exec.stack[top-1]
			top = top - 1
			idx1 := exec.stack[top-1]
			top = top - 1

			gs := &(exec.graphic_state)

			if opcode == kOpMD1 {
				pt0 := exec.point_at(gs.zp0, kPointTypeExec, idx0)
				pt1 := exec.point_at(gs.zp1, kPointTypeExec, idx1)
				exec.stack[top] = int32(dot_X(f26d6_t(pt0.X-pt1.X),
					f26d6_t(pt0.Y-pt1.Y), gs.projection_vector))
				top = top + 1
			} else {
				// TODO: do we need to check (h.gs.zp[0] == 0 || h.gs.zp[1] == 0)
				// as C Freetype does, similar to the MDRP instructions?
				pt0 := exec.point_at(gs.zp0, kPointTypeUnexec, idx0)
				pt1 := exec.point_at(gs.zp1, kPointTypeUnexec, idx1)
				exec.stack[top] = int32(dot_X(f26d6_t(pt0.X-pt1.X),
					f26d6_t(pt0.Y-pt1.Y), gs.projection_vector))
				top = top + 1
			}

		case kOpMPPEM, kOpMPS:
			if top >= len(exec.stack) {
				return errors.New("Freetype: exec stack overflow.")
			}

			// For MPS, point size should be irrelevant; we return the PPEM.
			exec.stack[top] = exec.scale >> 6
			top = top + 1

		case kOpFLIPON, kOpFLIPOFF:
			exec.graphic_state.auto_flip = (opcode == kOpFLIPON)

		case kOpDEBUG:
			// No-op.

		case kOpLT:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val1 < val0 {
				exec.stack[top] = 1
			} else {
				exec.stack[top] = 0
			}

			top = top + 1

		case kOpLTEQ:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val1 <= val0 {
				exec.stack[top] = 1
			} else {
				exec.stack[top] = 0
			}

			top = top + 1

		case kOpGT:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val1 > val0 {
				exec.stack[top] = 1
			} else {
				exec.stack[top] = 0
			}

			top = top + 1

		case kOpGTEQ:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val1 >= val0 {
				exec.stack[top] = 1
			} else {
				exec.stack[top] = 0
			}

			top = top + 1

		case kOpEQ:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val1 == val0 {
				exec.stack[top] = 1
			} else {
				exec.stack[top] = 0
			}

			top = top + 1

		case kOpNEQ:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val1 != val0 {
				exec.stack[top] = 1
			} else {
				exec.stack[top] = 0
			}

			top = top + 1

		case kOpODD, kOpEVEN:
			val := exec.round(f26d6_t(exec.stack[top-1])) >> 6
			exec.stack[top-1] = int32(val&1) ^ int32(opcode-kOpODD)

		case kOpIF:
			condition := exec.stack[top-1]
			top = top - 1
			if condition == 0 {
				opcode = 0
				if err := ifelse(); err != nil {
					return err
				}
				continue
			}

		case kOpEIF:
			// No-op.

		case kOpAND:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val1 == 0 || val0 == 0 {
				exec.stack[top] = 0
			} else {
				exec.stack[top] = 1
			}

			top = top + 1

		case kOpOR:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val1 == 0 && val0 == 0 {
				exec.stack[top] = 0
			} else {
				exec.stack[top] = 1
			}

			top = top + 1

		case kOpNOT:
			val0 := exec.stack[top-1]
			top = top - 1

			if val0 == 0 {
				exec.stack[top] = 1
			} else {
				exec.stack[top] = 0
			}

		case kOpDELTAP1:
			if err := deltap(); err != nil {
				return err
			}

		case kOpSDB:
			exec.graphic_state.delta_base = exec.stack[top-1]
			top = top - 1

		case kOpSDS:
			exec.graphic_state.delta_shift = exec.stack[top-1]
			top = top - 1

		case kOpADD:
			exec.stack[top-2] = exec.stack[top-2] + exec.stack[top-1]
			top = top - 1

		case kOpSUB:
			exec.stack[top-2] = exec.stack[top-2] - exec.stack[top-1]
			top = top - 1

		case kOpDIV:
			val0 := exec.stack[top-1]
			top = top - 1
			if val0 == 0 {
				return errors.New("Freetype: exec division by zero")
			}

			val1 := exec.stack[top-1]
			top = top - 1

			rst := f26d6_div(f26d6_t(val1), (f26d6_t(val0)))
			exec.stack[top] = int32(rst)
			top = top + 1

		case kOpMUL:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			rst := f26d6_mul(f26d6_t(val1), f26d6_t(val0))
			exec.stack[top] = int32(rst)
			top = top + 1

		case kOpABS:
			val := exec.stack[top-1]
			if val < 0 {
				exec.stack[top-1] = -val
			}

		case kOpNEG:
			exec.stack[top-1] = -exec.stack[top-1]

		case kOpFLOOR:
			exec.stack[top-1] = exec.stack[top-1] &^ 63

		case kOpCEILING:
			val := exec.stack[top-1]
			val = val + 63
			exec.stack[top-1] = val &^ 63

		case kOpROUND00, kOpROUND01, kOpROUND10, kOpROUND11:
			// The four flavors of opROUND are equivalent. See the comment below
			// on opNROUND for the rationale.
			exec.stack[top-1] = int32(exec.round(f26d6_t(exec.stack[top-1])))

		case kOpNROUND00, kOpNROUND01, kOpNROUND10, kOpNROUND11:
			// No-op. The spec says to add one of four "compensations for the
			// engine characteristics", to cater for things like "different
			// dot-size printers".
			// https://developer.apple.com/fonts/TTRefMan/RM02/Chap2.html#engine_compensation
			// This code does not implement engine compensation, as we don't
			// expect to be used to output on dot-matrix printers.

		case kOpDELTAP2, kOpDELTAP3:
			if err := deltap(); err != nil {
				return err
			}

		case kOpSROUND, kOpS45ROUND:
			val := exec.stack[top-1]
			top = top - 1

			gs := &(exec.graphic_state)

			switch (val >> 6) & 0x03 {
			case 0:
				gs.round_period = 1 << 5
			case 1, 3:
				gs.round_period = 1 << 6
			case 2:
				gs.round_period = 1 << 7
			}

			if opcode == kOpS45ROUND {
				// // The spec says to multiply by √2, but the C Freetype code
				// says 1/√2. We go with 1/√2.
				gs.round_period = gs.round_period * 46341
				gs.round_period = gs.round_period / 65536
			}

			gs.round_phase = gs.round_period * f26d6_t((val>>4)&0x03) / 4
			val = val & 0x0f
			if val != 0 {
				gs.round_threshold = gs.round_period * f26d6_t(val-4) / 8
			} else {
				gs.round_threshold = gs.round_period - 1
			}

		case kOpJROT:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val0 != 0 {
				pgm_count = pgm_count + int(val1)
				continue
			}

		case kOpJROF:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val0 == 0 {
				pgm_count = pgm_count + int(val1)
				continue
			}

		case kOpROFF:
			gs := &(exec.graphic_state)
			gs.round_period = 0
			gs.round_phase = 0
			gs.round_threshold = 0

		case kOpRUTG:
			gs := &(exec.graphic_state)
			gs.round_period = 1 << 6
			gs.round_phase = 0
			gs.round_threshold = 1<<6 - 1

		case kOpRDTG:
			gs := &(exec.graphic_state)
			gs.round_period = 1 << 6
			gs.round_phase = 0
			gs.round_threshold = 0

		case kOpSANGW, kOpAA:
			// These ops are "anachronistic" and no longer used.
			top = top - 1

		case kOpSCANCTRL:
			// We do not support dropout control, as we always rasterize
			// grayscale glyphs.
			top = top - 1

		case kOpGETINFO:
			rst := int32(0)
			val := exec.stack[top-1]

			if val&(1<<0) != 0 {
				// Set the engine version. We hard-code this to 35, the same as
				// the C freetype code, which says that "Version~35 corresponds
				// to MS rasterizer v.1.7 as used e.g. in Windows~98".
				rst = rst | 35
			}

			if val&(1<<5) != 0 {
				// Set that we support grayscale.
				rst = rst | 1<<12
			}

			exec.stack[top-1] = rst
		case kOpIDEF:
			// IDEF is for ancient versions of the bytecode interpreter, and is
			// no longer used.
			return errors.New("Freetype: exec unsupport IDEF instruction.")

		case kOpROLL:
			exec.stack[top-1], exec.stack[top-2], exec.stack[top-3] =
				exec.stack[top-3], exec.stack[top-1], exec.stack[top-2]

		case kOpMAX:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val1 < val0 {
				exec.stack[top] = val0
			}

			top = top + 1

		case kOpMIN:
			val0 := exec.stack[top-1]
			val1 := exec.stack[top-2]
			top = top - 2

			if val1 > val0 {
				exec.stack[top] = val0
			}

			top = top + 1

		case kOpSCANTYPE:
			// We do not support dropout control, as we always rasterize
			// grayscale glyphs.
			top = top - 1

		case kOpPUSHB000, kOpPUSHB001, kOpPUSHB010, kOpPUSHB011,
			kOpPUSHB100, kOpPUSHB101, kOpPUSHB110, kOpPUSHB111:
			opcode = opcode - kOpPUSHB000 + 1
			if err := push(); err != nil {
				return err
			}

			continue

		case kOpPUSHW000, kOpPUSHW001, kOpPUSHW010, kOpPUSHW011,
			kOpPUSHW100, kOpPUSHW101, kOpPUSHW110, kOpPUSHW111:
			opcode = opcode - kOpPUSHW000 + 1
			// log.Printf("PUSHW%v", opcode)
			opcode = opcode + 0x80
			if err := push(); err != nil {
				return err
			}

			continue
		case kOpMDRP00000, kOpMDRP00001, kOpMDRP00010, kOpMDRP00011,
			kOpMDRP00100, kOpMDRP00101, kOpMDRP00110, kOpMDRP00111,
			kOpMDRP01000, kOpMDRP01001, kOpMDRP01010, kOpMDRP01011,
			kOpMDRP01100, kOpMDRP01101, kOpMDRP01110, kOpMDRP01111,
			kOpMDRP10000, kOpMDRP10001, kOpMDRP10010, kOpMDRP10011,
			kOpMDRP10100, kOpMDRP10101, kOpMDRP10110, kOpMDRP10111,
			kOpMDRP11000, kOpMDRP11001, kOpMDRP11010, kOpMDRP11011,
			kOpMDRP11100, kOpMDRP11101, kOpMDRP11110, kOpMDRP11111:

			i := exec.stack[top-1]
			top = top - 1
			gs := &exec.graphic_state
			ref := exec.point_at(gs.zp0, kPointTypeExec, gs.rp0)
			p := exec.point_at(gs.zp1, kPointTypeExec, i)
			if ref == nil || p == nil {
				return errors.New("exec: point_at out of range")
			}

			oldDist := f26d6_t(0)
			if gs.zp0 == 0 || gs.zp1 == 0 {
				p0 := exec.point_at(gs.zp1, kPointTypeUnexec, i)
				p1 := exec.point_at(gs.zp0, kPointTypeUnexec, gs.rp0)
				oldDist = dot_X(f26d6_t(p0.X-p1.X), f26d6_t(p0.Y-p1.Y), gs.dual_vector)
			} else {
				p0 := exec.point_at(gs.zp1, kPointTypeRaw, i)
				p1 := exec.point_at(gs.zp0, kPointTypeRaw, gs.rp0)
				oldDist = dot_X(f26d6_t(p0.X-p1.X), f26d6_t(p0.Y-p1.Y), gs.dual_vector)
				oldDist = f26d6_t(exec.font.scale(exec.scale * int32(oldDist)))
			}

			// Single-width cut-in test.
			if x := f26d6_abs(oldDist - gs.single_width); x < gs.single_width_cut_in {
				if oldDist >= 0 {
					oldDist = +gs.single_width
				} else {
					oldDist = -gs.single_width
				}
			}

			// Rounding bit.
			// TODO: metrics compensation.
			distance := oldDist
			if opcode&0x04 != 0 {
				distance = exec.round(oldDist)
			}

			// Minimum distance bit.
			if opcode&0x08 != 0 {
				if oldDist >= 0 {
					if distance < gs.min_dist {
						distance = gs.min_dist
					}
				} else {
					if distance > -gs.min_dist {
						distance = -gs.min_dist
					}
				}
			}

			// Set-RP0 bit.
			gs.rp1 = gs.rp0
			gs.rp2 = i
			if opcode&0x10 != 0 {
				gs.rp0 = i
			}

			// Move the point.
			oldDist = dot_X(f26d6_t(p.X-ref.X), f26d6_t(p.Y-ref.Y), gs.projection_vector)
			exec.move_pt(p, distance-oldDist, true)

		case kOpMIRP00000, kOpMIRP00001, kOpMIRP00010, kOpMIRP00011,
			kOpMIRP00100, kOpMIRP00101, kOpMIRP00110, kOpMIRP00111,
			kOpMIRP01000, kOpMIRP01001, kOpMIRP01010, kOpMIRP01011,
			kOpMIRP01100, kOpMIRP01101, kOpMIRP01110, kOpMIRP01111,
			kOpMIRP10000, kOpMIRP10001, kOpMIRP10010, kOpMIRP10011,
			kOpMIRP10100, kOpMIRP10101, kOpMIRP10110, kOpMIRP10111,
			kOpMIRP11000, kOpMIRP11001, kOpMIRP11010, kOpMIRP11011,
			kOpMIRP11100, kOpMIRP11101, kOpMIRP11110, kOpMIRP11111:
			scale := exec.stack[top-1]
			top = top - 1

			cvt_dist := exec.get_scaled_cvt(scale)

			index := exec.stack[top-1]
			top = top - 1

			gs := &(exec.graphic_state)
			if f26d6_abs(cvt_dist-gs.single_width) < gs.single_width_cut_in {
				if cvt_dist >= 0 {
					cvt_dist = gs.single_width
				} else {
					cvt_dist = -gs.single_width
				}
			}

			if gs.zp1 == 0 {
				// TODO: implement once we have a .ttf file that triggers this.
				// So that we can step throungh C's freetype.
				return errors.New("Freetype: exec unimplement twilight point agjustment")
			}

			ref := exec.point_at(gs.zp0, kPointTypeUnexec, gs.rp0)
			pt := exec.point_at(gs.zp1, kPointTypeUnexec, index)
			if ref == nil || pt == nil {
				return errors.New("Freetype: exec point out of range.")
			}

			old_dist := dot_X(f26d6_t(pt.X-ref.X), f26d6_t(pt.Y-ref.Y),
				gs.dual_vector)

			ref = exec.point_at(gs.zp0, kPointTypeExec, gs.rp0)
			pt = exec.point_at(gs.zp1, kPointTypeExec, index)
			if ref == nil || pt == nil {
				return errors.New("Freetype: exec point out of range.")
			}

			cur_dist := dot_X(f26d6_t(pt.X-ref.X), f26d6_t(pt.Y-ref.Y),
				gs.projection_vector)

			if gs.auto_flip && old_dist^cvt_dist < 0 {
				cvt_dist = -cvt_dist
			}

			// Rounding bit.
			// TODO: metrics compension.
			dist := cvt_dist
			if opcode&0x04 != 0 {
				// The CVT value is only used if close enough to old_dist.
				if f26d6_abs(cvt_dist-old_dist) > gs.ctrl_val_cut_in {
					dist = old_dist
				}
				dist = exec.round(dist)
			}

			// Minimum distance bit.
			if opcode&0x08 != 0 {
				if old_dist >= 0 {
					if dist < gs.min_dist {
						dist = gs.min_dist
					}
				} else {
					if dist > -gs.min_dist {
						dist = -gs.min_dist
					}
				}
			}

			// Set-RP0 bit.
			gs.rp1 = gs.rp0
			gs.rp2 = index

			if opcode&0x10 != 0 {
				gs.rp0 = index
			}

			// Move the point.
			exec.move_pt(pt, dist-cur_dist, true)

		default:
			// log.Printf("%x", opcode)
			return errors.New("exec unrecognized instruction.")
		}

		pgm_count++

	} // for 0 <= pgm_count && pgm_count < len(pgm)

	return nil
}

func (exec *exec_t) point_at(zp int32, point_type point_type_t,
	index int32) *FontPoint {
	point_array := exec.points[zp][point_type]

	if index < 0 || len(point_array) <= int(index) {
		return nil
	}

	return &point_array[index]
}

func (exec *exec_t) move_pt(pt *FontPoint, dist f26d6_t, touch bool) {
	x0 := int64(exec.graphic_state.freedom_vector[0])
	x1 := int64(exec.graphic_state.projection_vector[0])

	if x0 == 0x4000 && x1 == 0x4000 {
		pt.X += int32(dist)
		if touch {
			pt.Flag |= kDecodeTouchedX
		}
		return
	}

	y0 := int64(exec.graphic_state.freedom_vector[1])
	y1 := int64(exec.graphic_state.projection_vector[1])

	if y0 == 0x4000 && y1 == 0x4000 {
		pt.Y += int32(dist)
		if touch {
			pt.Flag |= kDecodeTouchedY
		}
		return
	}

	fv_dot_pv := (x0*x1 + y0*y1) >> 14

	if x0 != 0 {
		pt.X += int32(mul_div(x0, int64(dist), fv_dot_pv))
		if touch {
			pt.Flag |= kDecodeTouchedX
		}
	}

	if y0 != 0 {
		pt.Y += int32(mul_div(y0, int64(dist), fv_dot_pv))
		if touch {
			pt.Flag |= kDecodeTouchedY
		}
	}
}

func (exec *exec_t) init_scaled_cvt() {
	exec.is_scaled_cvt_init = true

	if num := len(exec.font.cvt) / 2; num < cap(exec.scaled_cvt) {
		exec.scaled_cvt = exec.scaled_cvt[:num]
	} else {
		if num < 32 {
			num = 32
		}
		exec.scaled_cvt = make([]f26d6_t, len(exec.font.cvt)/2, num)
	}

	for i := range exec.scaled_cvt {
		unscaled := uint16(exec.font.cvt[2*i])<<8 | uint16(exec.font.cvt[2*i+1])
		scale := exec.font.scale(exec.scale * int32(int16(unscaled)))
		exec.scaled_cvt[i] = f26d6_t(scale)
	}
}

// https://developer.apple.com/fonts/TTRefMan/RM02/Chap2.html#rounding
func (exec *exec_t) round(f f26d6_t) f26d6_t {
	gs := &(exec.graphic_state)

	if gs.round_period == 0 {
		// Rounding is off.
		return f
	}

	if f >= 0 {
		r := (f - gs.round_phase + gs.round_threshold) & (-gs.round_period)
		if f != 0 && r < 0 {
			r = 0
		}
		return r + gs.round_phase
	}

	r := -((-f - gs.round_phase + gs.round_threshold) & (-gs.round_period))
	if r > 0 {
		r = 0
	}

	return r - gs.round_phase
}

func (exec *exec_t) iup_interp(is_y bool, p0, p1, r0, r1 int) {
	if p0 > p1 {
		return
	}

	if r0 >= len(exec.points[kZoneGlyph][kPointTypeExec]) ||
		r1 >= len(exec.points[kZoneGlyph][kPointTypeExec]) {
		return
	}

	var ifu0, ifu1 int32
	if is_y {
		ifu0 = exec.points[kZoneGlyph][kPointTypeRaw][r0].Y
		ifu1 = exec.points[kZoneGlyph][kPointTypeRaw][r1].Y
	} else {
		ifu0 = exec.points[kZoneGlyph][kPointTypeRaw][r0].X
		ifu1 = exec.points[kZoneGlyph][kPointTypeRaw][r1].X
	}

	if ifu0 > ifu1 {
		ifu0, ifu1 = ifu1, ifu0
		r0, r1 = r1, r0
	}

	var unh0, unh1, delta0, delta1 int32
	if is_y {
		unh0 = exec.points[kZoneGlyph][kPointTypeUnexec][r0].Y
		unh1 = exec.points[kZoneGlyph][kPointTypeUnexec][r1].Y
		delta0 = exec.points[kZoneGlyph][kPointTypeExec][r0].Y - unh0
		delta1 = exec.points[kZoneGlyph][kPointTypeExec][r1].Y - unh1
	} else {
		unh0 = exec.points[kZoneGlyph][kPointTypeUnexec][r0].X
		unh1 = exec.points[kZoneGlyph][kPointTypeUnexec][r1].X
		delta0 = exec.points[kZoneGlyph][kPointTypeExec][r0].X - unh0
		delta1 = exec.points[kZoneGlyph][kPointTypeExec][r1].X - unh1
	}

	var xy, ifu_xy int32
	if ifu0 == ifu1 {
		for i := p0; i <= p1; i++ {
			if is_y {
				xy = exec.points[kZoneGlyph][kPointTypeUnexec][i].Y
			} else {
				xy = exec.points[kZoneGlyph][kPointTypeUnexec][i].X
			}

			if xy <= unh0 {
				xy += delta0
			} else {
				xy += delta1
			}

			if is_y {
				exec.points[kZoneGlyph][kPointTypeExec][i].Y = xy
			} else {
				exec.points[kZoneGlyph][kPointTypeExec][i].X = xy
			}
		}
		return
	}

	scale, scale_ok := int64(0), false
	for i := p0; i <= p1; i++ {
		if is_y {
			xy = exec.points[kZoneGlyph][kPointTypeUnexec][i].Y
			ifu_xy = exec.points[kZoneGlyph][kPointTypeRaw][i].Y
		} else {
			xy = exec.points[kZoneGlyph][kPointTypeUnexec][i].X
			ifu_xy = exec.points[kZoneGlyph][kPointTypeRaw][i].X
		}

		if xy <= unh0 {
			xy += delta0
		} else if xy >= unh1 {
			xy += delta1
		} else {
			if !scale_ok {
				scale_ok = true
				scale = mul_div(int64(unh1+delta1-unh0-delta0), 0x10000,
					int64(ifu1-ifu0))
			}
			num := int64(ifu_xy-ifu0) * scale
			if num >= 0 {
				num += 0x8000
			} else {
				num -= 0x8000
			}
			xy = unh0 + delta0 + int32(num/0x10000)
		}

		if is_y {
			exec.points[kZoneGlyph][kPointTypeExec][i].Y = xy
		} else {
			exec.points[kZoneGlyph][kPointTypeExec][i].X = xy
		}

	}
}

func (exec *exec_t) iup_shift(is_y bool, p0, p1, p int) {
	var delta int32
	if is_y {
		delta = exec.points[kZoneGlyph][kPointTypeExec][p].Y -
			exec.points[kZoneGlyph][kPointTypeUnexec][p].Y
	} else {
		delta = exec.points[kZoneGlyph][kPointTypeExec][p].X -
			exec.points[kZoneGlyph][kPointTypeUnexec][p].X
	}

	if delta == 0 {
		return
	}

	for i := p0; i < p1; i++ {
		if i == p {
			continue
		}
		if is_y {
			exec.points[kZoneGlyph][kPointTypeExec][i].Y += delta
		} else {
			exec.points[kZoneGlyph][kPointTypeExec][i].X += delta
		}
	}
}

func (exec *exec_t) displace(flag_use_zp1 bool) (zp int32, i int32, d f26d6_t, ok bool) {
	gs := &(exec.graphic_state)
	zp, i = gs.zp0, gs.rp1
	if flag_use_zp1 {
		zp, i = gs.zp1, gs.rp2
	}
	p := exec.point_at(zp, kPointTypeExec, i)
	q := exec.point_at(zp, kPointTypeUnexec, i)
	if p == nil || q == nil {
		return gs.zp0, 0, 0, false
	}
	d = dot_X(f26d6_t(p.X-q.X), f26d6_t(p.Y-q.Y), gs.projection_vector)
	return zp, i, d, true
}

// mulDiv returns x*y/z, rounded to the nearest integer.
func mul_div(x, y, z int64) int64 {
	xy := x * y
	if z < 0 {
		xy, z = -xy, -z
	}
	if xy >= 0 {
		xy += z / 2
	} else {
		xy -= z / 2
	}
	return xy / z
}

func (exec *exec_t) get_scaled_cvt(idx int32) f26d6_t {
	if !exec.is_scaled_cvt_init {
		exec.init_scaled_cvt()
	}

	if idx < 0 || len(exec.scaled_cvt) <= int(idx) {
		return 0
	}

	return exec.scaled_cvt[idx]
}

func (exec *exec_t) set_scaled_cvt(idx int32, val f26d6_t) {
	if !exec.is_scaled_cvt_init {
		exec.init_scaled_cvt()
	}

	if idx < 0 || len(exec.scaled_cvt) <= int(idx) {
		return
	}

	exec.scaled_cvt[idx] = val
}
