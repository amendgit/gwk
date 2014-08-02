// Copyright 2012 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package freetype

func octets_to_u16(o []byte, i int) uint16 {
	return uint16(o[i])<<8 | uint16(o[i+1])
}

func octets_to_u32(o []byte, i int) uint32 {
	return uint32(o[i])<<24 | uint32(o[i+1])<<16 |
		uint32(o[i+2])<<8 | uint32(o[i+3])
}
