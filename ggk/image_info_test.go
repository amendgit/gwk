package ggk_test

import (
	"gwk/ggk"
	"testing"
)

var ctTests = []struct {
	ct       ggk.ColorType
	numBytes int
}{
	{ggk.ColorTypeUnknown, 0},
	{ggk.ColorTypeAlpha8, 1},
}

func TestColorType(t *testing.T) {
	var (
		ct       ggk.ColorType
		numBytes int
	)

	for _, tt := range ctTests {
		ct = tt.ct
		numBytes = ct.BytesPerPixel()
		if numBytes != tt.numBytes {
			t.Errorf("ColorTypeBytesPerPixel(%d) -> %d expect %d", ct, numBytes, tt.numBytes)
		}
	}
}
