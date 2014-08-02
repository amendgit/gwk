// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

import (
	"reflect"
	"strings"
	"testing"
	// "log"
)

func TestBytecode(t *testing.T) {
	testCases := []struct {
		desc   string
		prog   []byte
		want   []int32
		errStr string
	}{
		{
			"underflow",
			[]byte{
				kOpDUP,
			},
			nil,
			"underflow",
		},
		{
			"infinite loop",
			[]byte{
				kOpPUSHW000, // [-1]
				0xff,
				0xff,
				kOpDUP,  // [-1, -1]
				kOpJMPR, // [-1]
			},
			nil,
			"too many steps",
		},
		{
			"unbalanced if/else",
			[]byte{
				kOpPUSHB000, // [0]
				0,
				kOpIF,
			},
			nil,
			"unbalanced",
		},
		{
			"vector set/gets",
			[]byte{
				kOpSVTCA1,   // []
				kOpGPV,      // [0x4000, 0]
				kOpSVTCA0,   // [0x4000, 0]
				kOpGFV,      // [0x4000, 0, 0, 0x4000]
				kOpNEG,      // [0x4000, 0, 0, -0x4000]
				kOpSPVFS,    // [0x4000, 0]
				kOpSFVTPV,   // [0x4000, 0]
				kOpPUSHB000, // [0x4000, 0, 1]
				1,
				kOpGFV,      // [0x4000, 0, 1, 0, -0x4000]
				kOpPUSHB000, // [0x4000, 0, 1, 0, -0x4000, 2]
				2,
			},
			[]int32{0x4000, 0, 1, 0, -0x4000, 2},
			"",
		},
		{
			"jumps",
			[]byte{
				kOpPUSHB001, // [10, 2]
				10,
				2,
				kOpJMPR,     // [10]
				kOpDUP,      // not executed
				kOpDUP,      // [10, 10]
				kOpPUSHB010, // [10, 10, 20, 2, 1]
				20,
				2,
				1,
				kOpJROT,     // [10, 10, 20]
				kOpDUP,      // not executed
				kOpDUP,      // [10, 10, 20, 20]
				kOpPUSHB010, // [10, 10, 20, 20, 30, 2, 1]
				30,
				2,
				1,
				kOpJROF, // [10, 10, 20, 20, 30]
				kOpDUP,  // [10, 10, 20, 20, 30, 30]
				kOpDUP,  // [10, 10, 20, 20, 30, 30, 30]
			},
			[]int32{10, 10, 20, 20, 30, 30, 30},
			"",
		},
		{
			"stack ops",
			[]byte{
				kOpPUSHB010, // [10, 20, 30]
				10,
				20,
				30,
				kOpCLEAR,    // []
				kOpPUSHB010, // [40, 50, 60]
				40,
				50,
				60,
				kOpSWAP,     // [40, 60, 50]
				kOpDUP,      // [40, 60, 50, 50]
				kOpDUP,      // [40, 60, 50, 50, 50]
				kOpPOP,      // [40, 60, 50, 50]
				kOpDEPTH,    // [40, 60, 50, 50, 4]
				kOpCINDEX,   // [40, 60, 50, 50, 40]
				kOpPUSHB000, // [40, 60, 50, 50, 40, 4]
				4,
				kOpMINDEX, // [40, 50, 50, 40, 60]
			},
			[]int32{40, 50, 50, 40, 60},
			"",
		},
		{
			"push ops",
			[]byte{
				kOpPUSHB000, // [255]
				255,
				kOpPUSHW001, // [255, -2, 253]
				255,
				254,
				0,
				253,
				kOpNPUSHB, // [1, -2, 253, 1, 2]
				2,
				1,
				2,
				kOpNPUSHW, // [1, -2, 253, 1, 2, 0x0405, 0x0607, 0x0809]
				3,
				4,
				5,
				6,
				7,
				8,
				9,
			},
			[]int32{255, -2, 253, 1, 2, 0x0405, 0x0607, 0x0809},
			"",
		},
		{
			"store ops",
			[]byte{
				kOpPUSHB011, // [1, 22, 3, 44]
				1,
				22,
				3,
				44,
				kOpWS,       // [1, 22]
				kOpWS,       // []
				kOpPUSHB000, // [3]
				3,
				kOpRS, // [44]
			},
			[]int32{44},
			"",
		},
		{
			"comparison ops",
			[]byte{
				kOpPUSHB001, // [10, 20]
				10,
				20,
				kOpLT,       // [1]
				kOpPUSHB001, // [1, 10, 20]
				10,
				20,
				kOpLTEQ,     // [1, 1]
				kOpPUSHB001, // [1, 1, 10, 20]
				10,
				20,
				kOpGT,       // [1, 1, 0]
				kOpPUSHB001, // [1, 1, 0, 10, 20]
				10,
				20,
				kOpGTEQ, // [1, 1, 0, 0]
				kOpEQ,   // [1, 1, 1]
				kOpNEQ,  // [1, 0]
			},
			[]int32{1, 0},
			"",
		},
		{
			"odd/even",
			// Calculate odd(2+31/64), odd(2+32/64), even(2), even(1).
			[]byte{
				kOpPUSHB000, // [159]
				159,
				kOpODD,      // [0]
				kOpPUSHB000, // [0, 160]
				160,
				kOpODD,      // [0, 1]
				kOpPUSHB000, // [0, 1, 128]
				128,
				kOpEVEN,     // [0, 1, 1]
				kOpPUSHB000, // [0, 1, 1, 64]
				64,
				kOpEVEN, // [0, 1, 1, 0]
			},
			[]int32{0, 1, 1, 0},
			"",
		},
		{
			"if true",
			[]byte{
				kOpPUSHB001, // [255, 1]
				255,
				1,
				kOpIF,
				kOpPUSHB000, // [255, 2]
				2,
				kOpEIF,
				kOpPUSHB000, // [255, 2, 254]
				254,
			},
			[]int32{255, 2, 254},
			"",
		},
		{
			"if false",
			[]byte{
				kOpPUSHB001, // [255, 0]
				255,
				0,
				kOpIF,
				kOpPUSHB000, // [255]
				2,
				kOpEIF,
				kOpPUSHB000, // [255, 254]
				254,
			},
			[]int32{255, 254},
			"",
		},
		{
			"if/else true",
			[]byte{
				kOpPUSHB000, // [1]
				1,
				kOpIF,
				kOpPUSHB000, // [2]
				2,
				kOpELSE,
				kOpPUSHB000, // not executed
				3,
				kOpEIF,
			},
			[]int32{2},
			"",
		},
		{
			"if/else false",
			[]byte{
				kOpPUSHB000, // [0]
				0,
				kOpIF,
				kOpPUSHB000, // not executed
				2,
				kOpELSE,
				kOpPUSHB000, // [3]
				3,
				kOpEIF,
			},
			[]int32{3},
			"",
		},
		{
			"if/else true if/else false",
			// 0x58 is the opcode for opIF. The literal 0x58s below are pushed data.
			[]byte{
				kOpPUSHB010, // [255, 0, 1]
				255,
				0,
				1,
				kOpIF,
				kOpIF,
				kOpPUSHB001, // not executed
				0x58,
				0x58,
				kOpELSE,
				kOpPUSHW000, // [255, 0x5858]
				0x58,
				0x58,
				kOpEIF,
				kOpELSE,
				kOpIF,
				kOpNPUSHB, // not executed
				3,
				0x58,
				0x58,
				0x58,
				kOpELSE,
				kOpNPUSHW, // not executed
				2,
				0x58,
				0x58,
				0x58,
				0x58,
				kOpEIF,
				kOpEIF,
				kOpPUSHB000, // [255, 0x5858, 254]
				254,
			},
			[]int32{255, 0x5858, 254},
			"",
		},
		{
			"if/else false if/else true",
			// 0x58 is the opcode for opIF. The literal 0x58s below are pushed data.
			[]byte{
				kOpPUSHB010, // [255, 1, 0]
				255,
				1,
				0,
				kOpIF,
				kOpIF,
				kOpPUSHB001, // not executed
				0x58,
				0x58,
				kOpELSE,
				kOpPUSHW000, // not executed
				0x58,
				0x58,
				kOpEIF,
				kOpELSE,
				kOpIF,
				kOpNPUSHB, // [255, 0x58, 0x58, 0x58]
				3,
				0x58,
				0x58,
				0x58,
				kOpELSE,
				kOpNPUSHW, // not executed
				2,
				0x58,
				0x58,
				0x58,
				0x58,
				kOpEIF,
				kOpEIF,
				kOpPUSHB000, // [255, 0x58, 0x58, 0x58, 254]
				254,
			},
			[]int32{255, 0x58, 0x58, 0x58, 254},
			"",
		},
		{
			"logical ops",
			[]byte{
				kOpPUSHB010, // [0, 10, 20]
				0,
				10,
				20,
				kOpAND, // [0, 1]
				kOpOR,  // [1]
				kOpNOT, // [0]
			},
			[]int32{0},
			"",
		},
		{
			"arithmetic ops",
			// Calculate abs((-(1 - (2*3)))/2 + 1/64).
			// The answer is 5/2 + 1/64 in ideal numbers, or 161 in 26.6 fixed point math.
			[]byte{
				kOpPUSHB010, // [64, 128, 192]
				1 << 6,
				2 << 6,
				3 << 6,
				kOpMUL,      // [64, 384]
				kOpSUB,      // [-320]
				kOpNEG,      // [320]
				kOpPUSHB000, // [320, 128]
				2 << 6,
				kOpDIV,      // [160]
				kOpPUSHB000, // [160, 1]
				1,
				kOpADD, // [161]
				kOpABS, // [161]
			},
			[]int32{161},
			"",
		},
		{
			"floor, ceiling",
			[]byte{
				kOpPUSHB000, // [96]
				96,
				kOpFLOOR,    // [64]
				kOpPUSHB000, // [64, 96]
				96,
				kOpCEILING, // [64, 128]
			},
			[]int32{64, 128},
			"",
		},
		{
			"rounding",
			// Round 1.40625 (which is 90/64) under various rounding policies.
			// See figure 20 of https://developer.apple.com/fonts/TTRefMan/RM02/Chap2.html#rounding
			[]byte{
				kOpROFF,     // []
				kOpPUSHB000, // [90]
				90,
				kOpROUND00,  // [90]
				kOpRTG,      // [90]
				kOpPUSHB000, // [90, 90]
				90,
				kOpROUND00,  // [90, 64]
				kOpRTHG,     // [90, 64]
				kOpPUSHB000, // [90, 64, 90]
				90,
				kOpROUND00,  // [90, 64, 96]
				kOpRDTG,     // [90, 64, 96]
				kOpPUSHB000, // [90, 64, 96, 90]
				90,
				kOpROUND00,  // [90, 64, 96, 64]
				kOpRUTG,     // [90, 64, 96, 64]
				kOpPUSHB000, // [90, 64, 96, 64, 90]
				90,
				kOpROUND00,  // [90, 64, 96, 64, 128]
				kOpRTDG,     // [90, 64, 96, 64, 128]
				kOpPUSHB000, // [90, 64, 96, 64, 128, 90]
				90,
				kOpROUND00, // [90, 64, 96, 64, 128, 96]
			},
			[]int32{90, 64, 96, 64, 128, 96},
			"",
		},
		{
			"super-rounding",
			// See figure 20 of https://developer.apple.com/fonts/TTRefMan/RM02/Chap2.html#rounding
			// and the sign preservation steps of the "Order of rounding operations" section.
			[]byte{
				kOpPUSHB000, // [0x58]
				0x58,
				kOpSROUND,   // []
				kOpPUSHW000, // [-81]
				0xff,
				0xaf,
				kOpROUND00,  // [-80]
				kOpPUSHW000, // [-80, -80]
				0xff,
				0xb0,
				kOpROUND00,  // [-80, -80]
				kOpPUSHW000, // [-80, -80, -17]
				0xff,
				0xef,
				kOpROUND00,  // [-80, -80, -16]
				kOpPUSHW000, // [-80, -80, -16, -16]
				0xff,
				0xf0,
				kOpROUND00,  // [-80, -80, -16, -16]
				kOpPUSHB000, // [-80, -80, -16, -16, 0]
				0,
				kOpROUND00,  // [-80, -80, -16, -16, 16]
				kOpPUSHB000, // [-80, -80, -16, -16, 16, 16]
				16,
				kOpROUND00,  // [-80, -80, -16, -16, 16, 16]
				kOpPUSHB000, // [-80, -80, -16, -16, 16, 16, 47]
				47,
				kOpROUND00,  // [-80, -80, -16, -16, 16, 16, 16]
				kOpPUSHB000, // [-80, -80, -16, -16, 16, 16, 16, 48]
				48,
				kOpROUND00, // [-80, -80, -16, -16, 16, 16, 16, 80]
			},
			[]int32{-80, -80, -16, -16, 16, 16, 16, 80},
			"",
		},
		{
			"roll",
			[]byte{
				kOpPUSHB010, // [1, 2, 3]
				1,
				2,
				3,
				kOpROLL, // [2, 3, 1]
			},
			[]int32{2, 3, 1},
			"",
		},
		{
			"max/min",
			[]byte{
				kOpPUSHW001, // [-2, -3]
				0xff,
				0xfe,
				0xff,
				0xfd,
				kOpMAX,      // [-2]
				kOpPUSHW001, // [-2, -4, -5]
				0xff,
				0xfc,
				0xff,
				0xfb,
				kOpMIN, // [-2, -5]
			},
			[]int32{-2, -5},
			"",
		},
		{
			"functions",
			[]byte{
				kOpPUSHB011, // [3, 7, 0, 3]
				3,
				7,
				0,
				3,

				kOpFDEF, // Function #3 (not called)
				kOpPUSHB000,
				98,
				kOpENDF,

				kOpFDEF, // Function #0
				kOpDUP,
				kOpADD,
				kOpENDF,

				kOpFDEF, // Function #7
				kOpPUSHB001,
				10,
				0,
				kOpCALL,
				kOpDUP,
				kOpENDF,

				kOpFDEF, // Function #3 (again)
				kOpPUSHB000,
				99,
				kOpENDF,

				kOpPUSHB001, // [2, 0]
				2,
				0,
				kOpCALL,     // [4]
				kOpPUSHB000, // [4, 3]
				3,
				kOpLOOPCALL, // [99, 99, 99, 99]
				kOpPUSHB000, // [99, 99, 99, 99, 7]
				7,
				kOpCALL, // [99, 99, 99, 99, 20, 20]
			},
			[]int32{99, 99, 99, 99, 20, 20},
			"",
		},
	}

	for _, tc := range testCases {
		exe := &exec_t{}
		f := Font{
			max_storage:        32,
			max_stack_elements: 100,
		}
		exe.init(&f, 768)
		// log.Printf("%v B", tc.desc)
		err, errStr := exe.exec(tc.prog, nil, nil, nil, nil), ""
		// log.Printf("%v E", tc.desc)
		if err != nil {
			errStr = err.Error()
		}
		if tc.errStr != "" {
			if errStr == "" {
				t.Errorf("%s: got no error, want %q", tc.desc, tc.errStr)
			} else if !strings.Contains(errStr, tc.errStr) {
				t.Errorf("%s: got error %q, want one containing %q", tc.desc, errStr, tc.errStr)
			}
			continue
		}
		if errStr != "" {
			t.Errorf("%s: got error %q, want none", tc.desc, errStr)
			continue
		}
		got := exe.stack[:len(tc.want)]
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: got %v, want %v", tc.desc, got, tc.want)
			continue
		}
	}
}

// TestMove tests that the Hinter.move method matches the output of the C
// Freetype implementation.
func TestMove(t *testing.T) {
	exe, p := exec_t{}, FontPoint{}
	testCases := []struct {
		pvX, pvY, fvX, fvY f2d14_t
		wantX, wantY       int32
	}{
		{+0x4000, +0x0000, +0x4000, +0x0000, +1000, +0},
		{+0x4000, +0x0000, -0x4000, +0x0000, +1000, +0},
		{-0x4000, +0x0000, +0x4000, +0x0000, -1000, +0},
		{-0x4000, +0x0000, -0x4000, +0x0000, -1000, +0},
		{+0x0000, +0x4000, +0x0000, +0x4000, +0, +1000},
		{+0x0000, +0x4000, +0x0000, -0x4000, +0, +1000},
		{+0x4000, +0x0000, +0x2d41, +0x2d41, +1000, +1000},
		{+0x4000, +0x0000, -0x2d41, +0x2d41, +1000, -1000},
		{+0x4000, +0x0000, +0x2d41, -0x2d41, +1000, -1000},
		{+0x4000, +0x0000, -0x2d41, -0x2d41, +1000, +1000},
		{-0x4000, +0x0000, +0x2d41, +0x2d41, -1000, -1000},
		{-0x4000, +0x0000, -0x2d41, +0x2d41, -1000, +1000},
		{-0x4000, +0x0000, +0x2d41, -0x2d41, -1000, +1000},
		{-0x4000, +0x0000, -0x2d41, -0x2d41, -1000, -1000},
		{+0x376d, +0x2000, +0x2d41, +0x2d41, +732, +732},
		{-0x376d, +0x2000, +0x2d41, +0x2d41, -2732, -2732},
		{+0x376d, +0x2000, +0x2d41, -0x2d41, +2732, -2732},
		{-0x376d, +0x2000, +0x2d41, -0x2d41, -732, +732},
		{-0x376d, -0x2000, +0x2d41, +0x2d41, -732, -732},
		{+0x376d, +0x2000, +0x4000, +0x0000, +1155, +0},
		{+0x376d, +0x2000, +0x0000, +0x4000, +0, +2000},
	}
	for _, tc := range testCases {
		p = FontPoint{}
		gs := &exe.graphic_state
		gs.projection_vector = [2]f2d14_t{tc.pvX, tc.pvY}
		gs.freedom_vector = [2]f2d14_t{tc.fvX, tc.fvY}
		exe.move_pt(&p, 1000, true)
		tx := p.Flag&kDecodeTouchedX != 0
		ty := p.Flag&kDecodeTouchedY != 0
		wantTX := tc.fvX != 0
		wantTY := tc.fvY != 0
		if p.X != tc.wantX || p.Y != tc.wantY || tx != wantTX || ty != wantTY {
			t.Errorf("pv=%v, fv=%v\ngot  %d, %d, %t, %t\nwant %d, %d, %t, %t",
				gs.projection_vector, gs.freedom_vector, p.X, p.Y, tx, ty, tc.wantX, tc.wantY, wantTX, wantTY)
			continue
		}

		// Check that p is aligned with the freedom vector.
		a := int64(p.X) * int64(tc.fvY)
		b := int64(p.Y) * int64(tc.fvX)
		if a != b {
			t.Errorf("pv=%v, fv=%v, p=%v not aligned with fv", gs.projection_vector, gs.freedom_vector, p)
			continue
		}

		// Check that the projected p is 1000 away from the origin.
		dotProd := (int64(p.X)*int64(tc.pvX) + int64(p.Y)*int64(tc.pvY) + 1<<13) >> 14
		if dotProd != 1000 {
			t.Errorf("pv=%v, fv=%v, p=%v not 1000 from origin", gs.projection_vector, gs.freedom_vector, p)
			continue
		}
	}
}

// TestNormalize tests that the normalize function matches the output of the C
// Freetype implementation.
func TestNormalize(t *testing.T) {
	testCases := [][2]f2d14_t{
		{-15895, 3974},
		{-15543, 5181},
		{-14654, 7327},
		{-11585, 11585},
		{0, 16384},
		{11585, 11585},
		{14654, 7327},
		{15543, 5181},
		{15895, 3974},
		{16066, 3213},
		{16161, 2694},
		{16219, 2317},
		{16257, 2032},
		{16284, 1809},
	}
	for i, want := range testCases {
		got := normal(f2d14_t(i)-4, 1)
		if got != want {
			t.Errorf("i=%d: got %v, want %v", i, got, want)
		}
	}
}
