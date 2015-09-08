package ggk_test

import (
	"gwk/ggk"
	"testing"
)

var ctTests = []struct {
	ct       ggk.ColorType
	numBytes int
}{
	{ggk.ColorType_Unknown, 0},
	{ggk.ColorType_Alpha8, 1},
	{ggk.ColorType(1000), 0},
	{ggk.ColorType(-1), 0},
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
			t.Errorf("ColorType.BytesPerPixel(%d) -> %d expect %d", ct, numBytes, tt.numBytes)
		}
	}
}
