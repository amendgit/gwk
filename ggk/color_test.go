package ggk_test

import (
	"gwk/ggk"
	"testing"
)

var argbTests = []struct {
	a, r, g, b byte
	color      ggk.Color
}{
	{0xff, 0xff, 0xff, 0xff, 0xffffffff},
	{0x00, 0x00, 0x00, 0x00, 0x00000000},
	{0x11, 0x22, 0x33, 0x44, 0x11223344},
	{0xff, 0x00, 0x00, 0x00, 0xff000000},
	{0x00, 0xff, 0xff, 0xff, 0x00ffffff},
	{0x10, 0x20, 0x30, 0x40, 0x10203040},
}

var rgbTests = []struct {
	r, g, b byte
	color   ggk.Color
}{
	{0xff, 0xff, 0xff, 0xffffffff},
	{0x11, 0x22, 0x33, 0xff112233},
	{0x00, 0x00, 0x00, 0xff000000},
	{0x10, 0x10, 0x10, 0xff101010},
}

func TestColorInitAndComponents(t *testing.T) {
	// Test ColorWithRGB
	for _, tt := range rgbTests {
		var color = ggk.ColorWithRGB(tt.r, tt.g, tt.b)
		if color != tt.color {
			t.Errorf(`ColorWithRGB(0x%x, 0x%x, 0x%x)->0x%x want 0x%x`,
				tt.r, tt.g, tt.b, color, tt.color)
		}
	}

	// Test get color components.
	for _, tt := range argbTests {
		// Test ColorWithARGB.
		var color = ggk.ColorWithARGB(tt.a, tt.r, tt.g, tt.b)
		if color != tt.color {
			t.Errorf(`ColorWithARGB(0x%x, 0x%x, 0x%x, 0x%x)->0x%x want 0x%x`,
				tt.a, tt.r, tt.g, tt.b, color, tt.color)
		}

		// Test get color components.
		var a, r, g, b uint8

		a = tt.color.Alpha()
		if a != tt.a {
			t.Errorf(`0x%x.Alpha()->0x%x want 0x%x`, tt.color, a, tt.a)
		}

		r = tt.color.Red()
		if r != tt.r {
			t.Errorf(`0x%x.Red()->0x%x want 0x%x`, tt.color, r, tt.r)
		}

		g = tt.color.Green()
		if g != tt.g {
			t.Errorf(`0x%x.Green()->0x%x want 0x%x`, tt.color, g, tt.g)
		}

		b = tt.color.Blue()
		if b != tt.b {
			t.Errorf(`0x%x.Blue()->0x%x want 0x%x`, tt.color, b, tt.b)
		}

		// Test get ARGB method.
		a, r, g, b = tt.color.ARGB()
		if a != tt.a || r != tt.r || g != tt.g || b != tt.b {
			t.Errorf(`0x%x.ARGB()->(0x%x, 0x%x, 0x%x, 0x%x) want (0x%x, 0x%x, 0x%x, 0x%x)`,
				tt.color, a, r, g, b, tt.a, tt.r, tt.g, tt.b)
		}
	}
}

var setAlphaTests = []struct {
	a, b  ggk.Color
	alpha uint8
}{
	{0x00000000, 0x00000000, 0x00},
	{0x00000000, 0xff000000, 0xff},
	{0x00112233, 0x44112233, 0x44},
}

func TestColorSetAlphaTests(t *testing.T) {
	for _, tt := range setAlphaTests {
		var a = tt.a
		tt.a.SetAlpha(tt.alpha)
		if tt.a != tt.b {
			t.Errorf(`0x%x.SetAlpha(0x%x) -> 0x%x, want 0x%x`,
				a, tt.alpha, tt.a, tt.b)
		}
	}
}
