// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

import (
	"errors"
	"fmt"
	"log"
)

type Font struct {
	// Tables sliced from the TTF data. The different tables are documented at
	// http://developer.apple.com/fonts/TTRefMan/RM06/Chap6.html
	cmap []byte
	cvt  []byte
	fpgm []byte
	glyf []byte
	head []byte
	hhea []byte
	hmtx []byte
	kern []byte
	loca []byte
	maxp []byte
	prep []byte

	// Cached values derives from the raw ttf data.
	units_per_em       int32
	loca_offset_format int
	cmap_entry_array   []cmap_entry_t
	cmap_index_array   []byte
	glyph_num          int
	hmetric_num        int
	kern_num           int
	bounds             Bounds

	// Values from the maxp section.
	max_twilight_points uint16
	max_storage         uint16
	max_function_defs   uint16
	max_stack_elements  uint16
}

func Parse(ttf_bytes []byte) (*Font, error) {
	return parse_impl(ttf_bytes, 0)
}

func parse_impl(ttf_bytes []byte, index int) (font *Font, err error) {
	if len(ttf_bytes)-index < 12 {
		log.Printf("INVALID: TTF bytes too short.")
		return nil, errors.New("INVALID")
	}

	saved_index := index

	magic, index := octets_to_u32(ttf_bytes, index), index+4
	// log.Printf("magic %x", magic)

	if magic == 0x00010000 {
		// No-op.
	} else if magic == 0x74746366 {
		if saved_index != 0 {
			err = errors.New("INVALID: recursive TTC.")
			return
		}

		version, index := octets_to_u32(ttf_bytes, index), index+4
		if version != 0x00010000 {
			// TODO(needtime): support TTC version 2.0.
			err = errors.New("UNSUPPORT: ttc version 2.0")
			return
		}

		font_num, index := int(octets_to_u32(ttf_bytes, index)), index+4
		if font_num <= 0 {
			err = errors.New("INVALID: bad number of TTC fonts.")
			return
		}
		if len(ttf_bytes[index:])/4 < font_num {
			err = errors.New("INVALID: TTC offset table is too short.")
			return
		}

		// TODO: provide an API to select which font in a TrueType collection to
		// return, not just the first one. This may require an API to parse a
		// TTC's name tables, so users of this package can select the font in a
		// TTC by name.
		index = int(octets_to_u32(ttf_bytes, index))
		if index <= 0 || index > len(ttf_bytes) {
			err = errors.New("INVALID: ttf offset.")
			return
		}

		return parse_impl(ttf_bytes, index)
	} else {
		err = errors.New("INVALID: ttf version.")
		return
	}

	table_num, index := int(octets_to_u16(ttf_bytes, index)), index+2
	if len(ttf_bytes) < table_num*16+12 {
		err = errors.New("INVALID: ttf data is too short.")
		return
	}

	new_font := new(Font)
	for i := 0; i < table_num; i++ {
		table_offset := 16*i + 12

		title := string(ttf_bytes[table_offset : table_offset+4])
		begin := int(octets_to_u32(ttf_bytes, table_offset+8))
		length := int(octets_to_u32(ttf_bytes, table_offset+12))

		switch title {
		case "cmap":
			new_font.cmap, err = read_table(ttf_bytes, begin, length)

		case "head":
			new_font.head, err = read_table(ttf_bytes, begin, length)

		case "kern":
			new_font.kern, err = read_table(ttf_bytes, begin, length)

		case "maxp":
			new_font.maxp, err = read_table(ttf_bytes, begin, length)
		case "cvt ":
			new_font.cvt, err = read_table(ttf_bytes, begin, length)

		case "fpgm":
			new_font.fpgm, err = read_table(ttf_bytes, begin, length)

		case "glyf":
			new_font.glyf, err = read_table(ttf_bytes, begin, length)

		case "hmtx":
			new_font.hmtx, err = read_table(ttf_bytes, begin, length)

		case "loca":
			new_font.loca, err = read_table(ttf_bytes, begin, length)

		case "prep":
			new_font.prep, err = read_table(ttf_bytes, begin, length)

		case "hhea":
			new_font.hhea, err = read_table(ttf_bytes, begin, length)

		}

		if err != nil {
			return
		}
	}

	if err = new_font.parse_head(); err != nil {
		return
	}

	if err = new_font.parse_cmap(); err != nil {
		return
	}

	if err = new_font.parse_maxp(); err != nil {
		return
	}

	if err = new_font.parse_hhea(); err != nil {
		return
	}

	if err = new_font.parse_kern(); err != nil {
		return
	}

	return new_font, nil
}

func read_table(ttf_bytes []byte, begin int, length int) ([]byte, error) {
	if begin < 0 {
		return nil, errors.New("INVALID: begin too large.")
	}

	if length < 0 {
		return nil, errors.New("INVALID: length too large.")
	}

	end := begin + length
	if end < 0 || end > len(ttf_bytes) {
		return nil, errors.New("INVALID: begin + length too large.")
	}
	return ttf_bytes[begin:end], nil
}

type cmap_entry_t struct {
	start_code      uint32
	end_code        uint32
	id_delta        uint32
	id_range_offset uint32
}

// https://developer.apple.com/fonts/TTRefMan/RM06/Chap6cmap.html
func (f *Font) parse_cmap() error {
	if len(f.cmap) < 4 {
		log.Print("Font cmap too short.")
	}

	index := 2
	subtable_num, index := int(octets_to_u16(f.cmap, index)), index+2

	if len(f.cmap) < subtable_num*8+4 {
		log.Print("Font cmap too short.")
	}

	valid := false
	offset := 0
	for i := 0; i < subtable_num; i++ {
		// platform id is platform identifier, platform specific id is platform
		// specific encoding identifier.
		var pid_psid uint32
		pid_psid, index = octets_to_u32(f.cmap, index), index+4
		offset, index = int(octets_to_u32(f.cmap, index)), index+4

		// Unicode encoding.
		if pid_psid == 0x00000003 {
			valid = true
			break
		}

		// Microsoft UCS-2 Encoding or Microsoft UCS-4 Encoding.
		if pid_psid == 0x00030001 || pid_psid == 0x0003000a {
			valid = true
			// Don't break. So that unicode encoding can override ms encoding.
		}

		// TODO(coding): support whole list about pid and psid.
		// https://developer.apple.com/fonts/TTRefMan/RM06/Chap6name.html#ID
	}

	if !valid {
		return errors.New("UNSUPPORT or INVALID: cmap language encoding.")
	}

	cmap_format_version, offset := octets_to_u16(f.cmap, offset), offset+2
	if cmap_format_version == 4 { // cmap format 2
		// uint16, Length of subtable in bytes.
		// length, offset := octets_to_u16(font.cmap, offset), offset+2
		offset = offset + 2

		// uint16, Language code for this encoding subtable, or zero if
		// language-independent.
		lang, offset := octets_to_u16(f.cmap, offset), offset+2
		if lang != 0 {
			return errors.New("UNSUPPORT: cmap language isn't independent.")
		}

		// uint16, 2 * segCount.
		seg_count_x_2, offset := int(octets_to_u16(f.cmap, offset)), offset+2
		seg_count := seg_count_x_2 / 2

		// uint16, 2 * (2**FLOOR(log2(segCount))).
		// search_range, offset := octets_to_u16(font.cmap, offset), offset+2
		offset = offset + 2

		// uint16, log2(searchRange/2)
		// entry_selector, offset := octets_to_u16(font.cmap, offset), offset+2
		offset = offset + 2

		// uint16, (2 * segCount) - searchRange.
		// range_shift, offset := octets_to_u16(font.cmap, offset), offset+2
		offset = offset + 2
		// log.Printf("seg_count %v", seg_count)

		f.cmap_entry_array = make([]cmap_entry_t, seg_count)
		// uint16 * seg_count, Ending character code for each segment,
		// last = 0xFFFF.
		for i := 0; i < seg_count; i++ {
			f.cmap_entry_array[i].end_code, offset =
				uint32(octets_to_u16(f.cmap, offset)), offset+2
			// log.Printf("end_code %v", f.cmap_entry_array[i].end_code)
		}

		// uint16, This value should be zero.
		// reserved_pad, offset := octets_to_u16(font.cmap, offset), offset+2
		offset = offset + 2

		// uint16 * seg_count, Starting character code for each segment.
		for i := 0; i < seg_count; i++ {
			f.cmap_entry_array[i].start_code, offset =
				uint32(octets_to_u16(f.cmap, offset)), offset+2
			// log.Printf("start_code %v", f.cmap_entry_array[i].start_code)
		}

		// uint16 * seg_count, Delta for all character codes in segment.
		for i := 0; i < seg_count; i++ {
			f.cmap_entry_array[i].id_delta, offset =
				uint32(octets_to_u16(f.cmap, offset)), offset+2
			// log.Printf("id_delta %v", f.cmap_entry_array[i].id_delta)
		}

		// uint16 * seg_count, Offset in bytes to glyph indexArray, or 0.
		for i := 0; i < seg_count; i++ {
			f.cmap_entry_array[i].id_range_offset, offset =
				uint32(octets_to_u16(f.cmap, offset)), offset+2
			// log.Printf("id_range_offset %v", f.cmap_entry_array[i].id_range_offset)
		}

		// uint16 * seg_count, Glyph index array.
		f.cmap_index_array = f.cmap[offset:]

		return nil
	} else if cmap_format_version == 12 {
		// Format 12.0 is a bit like format 4, in that it defines segments for
		// sparse representation in 4-byte character space.

		// So, the next two bytes is part of version segment and should be 0.
		expect_zero, offset := octets_to_u16(f.cmap, offset), offset+2
		if expect_zero != 0 {
			msg := fmt.Sprint("UNSUPPORT or INVALID: cmap format version %x",
				f.cmap[offset-4:offset])
			return errors.New(msg)
		}

		// uint32, Byte length of this subtable (including the header).
		// length, offset := octets_to_u32(font.cmap, offset), offset+4
		offset = offset + 4

		// uint32, 0 if don't care.
		// lang, offset := octets_to_u32(font.cmap, offset), offset+4
		offset = offset + 4

		// uint32, Number of groupings which follow.
		group_num, offset := octets_to_u32(f.cmap, offset), offset+4
		// log.Printf("group_num %v", group_num)

		// Here follow the individual groups.
		for i := uint32(0); i < group_num; i++ {
			// uint32, First character code in this group.
			f.cmap_entry_array[i].start_code, offset =
				octets_to_u32(f.cmap, offset), offset+4

			// uint32, Last character code in this group.
			f.cmap_entry_array[i].end_code, offset =
				octets_to_u32(f.cmap, offset), offset+4

			// uint32, Glyph index corresponding to the starting character code.
			f.cmap_entry_array[i].id_delta, offset =
				octets_to_u32(f.cmap, offset), offset+4
		}

		return nil
	} else {
		msg := fmt.Sprintf("UNSUPPORT: cmap format version %v",
			cmap_format_version)
		return errors.New(msg)
	}
}

const (
	kLocaOffsetFormatUnknown int = iota
	kLocaOffsetFormatShort
	kLocaOffsetFormatLong
)

// https://developer.apple.com/fonts/TTRefMan/RM06/Chap6head.html
func (f *Font) parse_head() error {
	if len(f.head) != 54 {
		msg := fmt.Sprintf("INVALID: bad head length %v", len(f.head))
		return errors.New(msg)
	}

	// Range from 64 to 16384
	f.units_per_em = int32(octets_to_u16(f.head, 18))
	// log.Printf("units_per_em %d", f.units_per_em)
	f.bounds.XMin = int32(int16(octets_to_u16(f.head, 36)))
	f.bounds.YMin = int32(int16(octets_to_u16(f.head, 38)))
	f.bounds.XMax = int32(int16(octets_to_u16(f.head, 40)))
	f.bounds.YMax = int32(int16(octets_to_u16(f.head, 42)))

	// 0 for short offsets, 1 for long offsets.
	index_to_loc_format := octets_to_u16(f.head, 50)
	// log.Printf("index_to_loc_format %d", index_to_loc_format)
	if index_to_loc_format == 0 {
		f.loca_offset_format = kLocaOffsetFormatShort
	} else if index_to_loc_format == 1 {
		f.loca_offset_format = kLocaOffsetFormatLong
	} else {
		msg := fmt.Sprintf("INVALID: bad head indexToLocFormat %v",
			index_to_loc_format)
		return errors.New(msg)
	}

	return nil
}

// http://developer.apple.com/fonts/TTRefMan/RM06/Chap6kern.html
func (font *Font) parse_kern() error {
	if len(font.kern) <= 0 {
		if font.kern_num != 0 {
			return errors.New("INVALID: kern length.")
		} else {
			return nil
		}
	}

	index := 0

	// uint16, The version number of the kerning table (0x00010000 for the
	// current version).
	//
	// Upto now, only support the older version. Windows only support the older
	// version. Mac support both, but prefer the newer version.
	//
	// TODO(coding): Support the newer version.
	kern_format_version, index := octets_to_u16(font.kern, index), index+2

	if kern_format_version == 0 {

		// uint16, The number of subtables included in the kerning table.
		table_num, index := octets_to_u16(font.kern, index), index+2
		if table_num != 1 {
			msg := fmt.Sprintf("UNSUPPORT: kern table num %v", table_num)
			return errors.New(msg)
		}

		index = index + 2

		// uint16, The length of this subtable in bytes, including this header.
		length, index := int(octets_to_u16(font.kern, index)), index+2

		// uint16, Circumstances under which this table is used. See below for
		// description.
		coverage, index := octets_to_u16(font.kern, index), index+2
		if coverage != 0x0001 {
			// Upto now, we don't support horizontal kerning.
			// TODO(coding): support the horizontal kerning.
			msg := fmt.Sprintf("UNSUPPORT: kern coverage: 0x%04x", coverage)
			return errors.New(msg)
		}

		// uint16, number of kern.
		font.kern_num, index = int(octets_to_u16(font.kern, index)), index+2
		if font.kern_num*6 != length-14 {
			msg := fmt.Sprintf("INVALID: Bad kern table length")
			return errors.New(msg)
		}

		return nil
	}

	msg := fmt.Sprintf("UNSUPPORT: kern format version %v.",
		kern_format_version)
	return errors.New(msg)
}

// https://developer.apple.com/fonts/TTRefMan/RM06/Chap6maxp.html
func (font *Font) parse_maxp() error {
	if len(font.maxp) != 32 {
		msg := fmt.Sprintf("INVALID: bad maxp length %v", len(font.maxp))
		return errors.New(msg)
	}

	index := 0

	// Fixed 0x00010000, maxp format version.
	version, index := octets_to_u32(font.maxp, index), index+4
	if version != 0x00010000 {
		msg := fmt.Sprintf("UNSUPPORT: font maxp version %v.", version)
		return errors.New(msg)
	}

	// uint16, the number of glyphs in the font
	font.glyph_num, index = int(octets_to_u16(font.maxp, index)), index+2

	// uint16, points in non-compound glyph.
	// max_points, index := octets_to_u16(font.maxp, index), index+2
	index = index + 2

	// uint16, contours in non-compound glyph.
	// max_contours, index := octets_to_u16(font.maxp, index), index+2
	index = index + 2

	// uint16, points in compound glyph.
	// max_component_points, index := octets_to_u16(font.maxp, index), index+2
	index = index + 2

	// uint16, contours in compound glyph.
	// max_component_contours, index := octets_to_u16(font.maxp, index), index+2
	index = index + 2

	// uint16, set to 2.
	// max_zones, index := octets_to_u16(font.maxp, index), index+2
	index = index + 2

	// uint16, points used in Twilight Zone (Z0).
	font.max_twilight_points, index = octets_to_u16(font.maxp, index), index+2

	// uint16, number of Storage Area locations.
	font.max_storage, index = octets_to_u16(font.maxp, index), index+2

	// uint16, number of FDEFs.
	font.max_function_defs, index = octets_to_u16(font.maxp, index), index+2

	// uint16, number of IDEFs.
	// max_instruction_defs, index := octets_to_u16(font.maxp, index), index+2
	index = index + 2

	// uint16, maximum stack depth.
	font.max_stack_elements, index = octets_to_u16(font.maxp, index), index+2

	// uint16, number of glyphs referenced at top level.
	// max_component_elements, index := octets_to_u16(font.maxp, index), index+2
	index = index + 2

	// levels of recursion, set to 0 if font has only simple glyphs.
	// max_component_depth, index := octets_to_u16(font.maxp, index), index+2
	index = index + 2

	return nil
}

// https://developer.apple.com/fonts/TTRefMan/RM06/Chap6hhea.html
func (f *Font) parse_hhea() error {
	if len(f.hhea) != 36 {
		msg := fmt.Sprintf("INVALID: Bad hhea length %v", len(f.hhea))
		return errors.New(msg)
	}

	// TODO(coding): complete this.
	f.hmetric_num = int(octets_to_u16(f.hhea, 34))
	if f.hmetric_num*4+(f.glyph_num-f.hmetric_num)*2 != len(f.hmtx) {
		msg := fmt.Sprintf("INVALID: Bad hmtx length %v", len(f.hmtx))
		return errors.New(msg)
	}

	return nil
}

func (font *Font) scale(x int32) int32 {
	if x >= 0 {
		x += font.units_per_em / 2
	} else {
		x -= font.units_per_em / 2
	}
	return x / font.units_per_em
}

func (font *Font) unscaled_hmetric(idx uint16) HMetric {
	i := int(idx)
	if i >= font.glyph_num {
		return HMetric{}
	}

	if i >= font.hmetric_num {
		p := 4 * (font.hmetric_num - 1)
		return HMetric{
			AdvanceWidth:    int32(octets_to_u16((font.hmtx), p)),
			LeftSideBearing: int32(int16(octets_to_u16(font.hmtx, p+2*(i-font.hmetric_num)+4))),
		}
	}
	return HMetric{
		AdvanceWidth:    int32(octets_to_u16(font.hmtx, 4*i)),
		LeftSideBearing: int32(int16(octets_to_u16(font.hmtx, 4*i+2))),
	}
}

func (f *Font) FUnitsPerEm() int32 {
	return f.units_per_em
}

// Bounds returns the union of a Font's glyphs' bounds.
func (f *Font) Bounds(scale int32) Bounds {
	b := f.bounds
	b.XMin = f.scale(scale * b.XMin)
	b.YMin = f.scale(scale * b.YMin)
	b.XMax = f.scale(scale * b.XMax)
	b.YMax = f.scale(scale * b.YMax)
	return b
}

// HMetric returns the horizontal metrics for the glyph with the given index.
func (f *Font) HMetric(scale int32, i uint16) HMetric {
	h := f.unscaled_hmetric(i)
	h.AdvanceWidth = f.scale(scale * h.AdvanceWidth)
	h.LeftSideBearing = f.scale(scale * h.LeftSideBearing)
	return h
}

func (f *Font) Kerning(scale int32, i0, i1 uint16) int32 {
	if f.kern_num == 0 {
		return 0
	}
	g := uint32(i0)<<16 | uint32(i1)
	lo, hi := 0, f.kern_num
	for lo < hi {
		i := (lo + hi) / 2
		ig := octets_to_u32(f.kern, 18+6*i)
		if ig < g {
			lo = i + 1
		} else if ig > g {
			hi = i
		} else {
			return f.scale(scale * int32(int16(octets_to_u16(f.kern, 22+6*i))))
		}
	}
	return 0
}

// Index returns a Font's index for the given rune.
func (f *Font) Index(r rune) uint16 {
	// log.Printf("r %v", r)
	// log.Printf("len %v", len(f.cmap_entry_array))
	cm_len := len(f.cmap_entry_array)

	x := uint32(r)
	for lo, hi := 0, cm_len; lo < hi; {
		mi := lo + (hi-lo)/2
		// log.Printf("lo hi mi %v %v %v", lo, hi, mi)
		cm := &f.cmap_entry_array[mi]
		// log.Printf("cm %v", cm)
		if x < cm.start_code {
			hi = mi
		} else if cm.end_code < x {
			lo = mi + 1
		} else if cm.id_range_offset == 0 {
			// log.Printf("x id_delta %v %v", x, cm.id_delta)
			return uint16(x + cm.id_delta)
		} else {
			offset := int(cm.id_range_offset) + 2*(mi-cm_len+int(x-cm.start_code))
			return uint16(octets_to_u16(f.cmap_index_array, offset))
		}
	}

	return 0
}
