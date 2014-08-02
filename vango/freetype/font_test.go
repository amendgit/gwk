// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	// "log"
	"os"
	"strconv"
	"strings"
	"testing"
)

func parseTestdataFont(name string) (font *Font, testdataIsOptional bool, err error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("./exp/data/%s.ttf", name))
	if err != nil {
		// The "x-foo" fonts are optional tests, as they are not checked
		// in for copyright or file size reasons.
		return nil, strings.HasPrefix(name, "x-"), fmt.Errorf("%s: ReadFile: %v", name, err)
	}
	font, err = Parse(b)
	if err != nil {
		return nil, true, fmt.Errorf("%s: Parse: %v", name, err)
	}
	return font, false, nil
}

// TestParse tests that the luxisr.ttf metrics and glyphs are parsed correctly.
// The numerical values can be manually verified by examining luxisr.ttx.
func TestParse(t *testing.T) {
	font, _, err := parseTestdataFont("luxisr")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := font.FUnitsPerEm(), int32(2048); got != want {
		t.Errorf("FUnitsPerEm: got %v, want %v", got, want)
	}
	fupe := font.FUnitsPerEm()
	if got, want := font.Bounds(fupe), (Bounds{-441, -432, 2024, 2033}); got != want {
		t.Errorf("Bounds: got %v, want %v", got, want)
	}

	i0 := font.Index('A')
	i1 := font.Index('V')
	if i0 != 36 || i1 != 57 {
		t.Fatalf("Index: i0, i1 = %d, %d, want 36, 57", i0, i1)
	}
	if got, want := font.HMetric(fupe, i0), (HMetric{1366, 19}); got != want {
		t.Errorf("HMetric: got %v, want %v", got, want)
	}
	if got, want := font.Kerning(fupe, i0, i1), int32(-144); got != want {
		t.Errorf("Kerning: got %v, want %v", got, want)
	}

	g := NewGlyph()
	err = g.Load(font, fupe, i0, nil)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	g0 := &Glyph{
		Rect:          g.Rect,
		AllPoints:     g.AllPoints,
		EndIndexArray: g.EndIndexArray,
	}
	g1 := &Glyph{
		Rect: Bounds{19, 0, 1342, 1480},
		AllPoints: []FontPoint{
			{19, 0, 51},
			{581, 1480, 1},
			{789, 1480, 51},
			{1342, 0, 1},
			{1116, 0, 35},
			{962, 410, 3},
			{368, 410, 33},
			{214, 0, 3},
			{428, 566, 19},
			{904, 566, 33},
			{667, 1200, 3},
		},
		EndIndexArray: []int{8, 11},
	}
	if got, want := fmt.Sprint(g0), fmt.Sprint(g1); got != want {
		t.Errorf("GlyphBuf:\ngot  %v\nwant %v", got, want)
	}
}

func TestIndex(t *testing.T) {
	testCases := map[string]map[rune]uint16{
		"luxisr": {
			' ':      3,
			'!':      4,
			'A':      36,
			'V':      57,
			'É':      101,
			'ﬂ':      193,
			'\u22c5': 385,
			'中':      0,
		},
		"x-arial-bold": {
			' ':      3,
			'+':      14,
			'0':      19,
			'_':      66,
			'w':      90,
			'~':      97,
			'Ä':      98,
			'ﬂ':      192,
			'½':      242,
			'σ':      305,
			'λ':      540,
			'ỹ':      1275,
			'\u04e9': 1319,
			'中':      0,
		},
		"x-deja-vu-sans-oblique": {
			' ':      3,
			'*':      13,
			'Œ':      276,
			'ω':      861,
			'‡':      2571,
			'⊕':      3109,
			'ﬂ':      4560,
			'\ufb03': 4561,
			'\ufffd': 4645,
			// TODO: '\U0001f640': ???,
			'中': 0,
		},
		"x-droid-sans-japanese": {
			' ':      0,
			'\u3000': 3,
			'\u3041': 25,
			'\u30fe': 201,
			'\uff61': 202,
			'\uff67': 208,
			'\uff9e': 263,
			'\uff9f': 264,
			'\u4e00': 265,
			'\u557e': 1000,
			'\u61b6': 2024,
			'\u6ede': 3177,
			'\u7505': 3555,
			'\u81e3': 4602,
			'\u81e5': 4603,
			'\u81e7': 4604,
			'\u81e8': 4605,
			'\u81ea': 4606,
			'\u81ed': 4607,
			'\u81f3': 4608,
			'\u81f4': 4609,
			'\u91c7': 5796,
			'\u9fa0': 6620,
			'\u203e': 12584,
		},
		"x-times-new-roman": {
			' ':      3,
			':':      29,
			'ﬂ':      192,
			'Ŀ':      273,
			'♠':      388,
			'Ŗ':      451,
			'Σ':      520,
			'\u200D': 745,
			'Ẽ':      1216,
			'\u04e9': 1319,
			'中':      0,
		},
	}
	for name, wants := range testCases {
		font, testdataIsOptional, err := parseTestdataFont(name)
		if err != nil {
			if testdataIsOptional {
				t.Log(err)
			} else {
				t.Fatal(err)
			}
			continue
		}
		for r, want := range wants {
			if got := font.Index(r); got != want {
				t.Errorf("%s: Index of %q, aka %U: got %d, want %d", name, r, r, got, want)
			}
		}
	}
}

// scalingTestParse parses a line of points like
// -22 -111 1, 178 555 1, 236 555 1, 36 -111 1
// The line will not have a trailing "\n".
func scalingTestParse(line string) []FontPoint {
	if line == "" {
		return nil
	}
	points := make([]FontPoint, 0, 1+strings.Count(line, ","))
	for len(line) > 0 {
		s := line
		if i := strings.Index(line, ","); i != -1 {
			s, line = line[:i], line[i+1:]
			for len(line) > 0 && line[0] == ' ' {
				line = line[1:]
			}
		} else {
			line = ""
		}
		i := strings.Index(s, " ")
		if i == -1 {
			break
		}
		x, _ := strconv.Atoi(s[:i])
		s = s[i+1:]
		i = strings.Index(s, " ")
		if i == -1 {
			break
		}
		y, _ := strconv.Atoi(s[:i])
		s = s[i+1:]
		f, _ := strconv.Atoi(s)
		points = append(points, FontPoint{
			X:    int32(x),
			Y:    int32(y),
			Flag: uint32(f),
		})
	}
	return points
}

// scalingTestEquals is equivalent to, but faster than, calling
// reflect.DeepEquals(a, b). It also treats a nil []Point and an empty non-nil
// []Point as equal.
func scalingTestEquals(a, b []FontPoint) bool {
	if len(a) != len(b) {
		return false
	}
	for i, p := range a {
		if p != b[i] {
			return false
		}
	}
	return true
}

var scalingTestCases = []struct {
	name string
	size int32
	// hintingBrokenAt, if non-negative, is the glyph index n for which
	// only the first n glyphs are known to be correctly hinted.
	// TODO: remove this field, when hinting is completely implemented.
	hintingBrokenAt int
}{
	{"luxisr", 12, -1},
	{"x-arial-bold", 11, 0},
	{"x-deja-vu-sans-oblique", 17, 229},
	{"x-droid-sans-japanese", 9, 0},
	{"x-times-new-roman", 13, 0},
}

// TODO: also test bounding boxes, not just points.

func testScaling(t *testing.T, exec *exec_t) {
	for _, tc := range scalingTestCases {
		font, testdataIsOptional, err := parseTestdataFont(tc.name)
		if err != nil {
			if testdataIsOptional {
				t.Log(err)
			} else {
				t.Error(err)
			}
			continue
		}

		hinting := "sans"
		if exec != nil {
			hinting = "with"
		}
		f, err := os.Open(fmt.Sprintf("./exp/data/%s-%dpt-%s-hinting.txt", tc.name, tc.size, hinting))
		if err != nil {
			t.Errorf("%s: Open: %v", tc.name, err)
			continue
		}
		defer f.Close()

		wants := [][]FontPoint{}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			wants = append(wants, scalingTestParse(scanner.Text()))
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			t.Errorf("%s: Scanner: %v", tc.name, err)
			continue
		}

		glyphBuf := NewGlyph()
		for i, want := range wants {
			// TODO: completely implement hinting. For now, only the first
			// tc.hintingBrokenAt glyphs of the test case's font are correctly hinted.
			if exec != nil && i == tc.hintingBrokenAt {
				break
			}

			if err = glyphBuf.Load(font, tc.size*64, uint16(i), exec); err != nil {
				t.Errorf("%s: glyph #%d: Load: %v", tc.name, i, err)
				continue
			}
			got := glyphBuf.AllPoints
			for i := range got {
				got[i].Flag &= 0x01
			}
			if !scalingTestEquals(got, want) {
				t.Errorf("%s: glyph #%d:\ngot  %v\nwant %v\n", tc.name, i, got, want)
				return
			}
		}
	}
}

func TestScalingSansHinting(t *testing.T) {
	testScaling(t, nil)
}

func TestScalingWithHinting(t *testing.T) {
	os.Remove("./log.txt")
	testScaling(t, &exec_t{})
}
