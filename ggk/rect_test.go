package ggk_test

import (
	"gwk/ggk"
	"testing"
)

var rectEuqalTest = []struct {
	a       ggk.Rect
	b       ggk.Rect
	isEqual bool
}{
	{
		ggk.MakeRect(0.0, 0.0, 0.0, 0.0),
		ggk.MakeRect(0.0, 0.0, 0.0, 0.0),
		true,
	},
	{
		ggk.MakeRect(1.0, 1.0, 100.0, 100.0),
		ggk.MakeRect(1.0, 1.0, 100.0, 100.0),
		true,
	},
	{
		ggk.MakeRect(1.0, 1.0, 100.0, 99.0),
		ggk.MakeRect(1.0, 1.0, 100.0, 100.0),
		false,
	},
}

func TestRectEqual(t *testing.T) {
	for _, tt := range rectEuqalTest {
		var isEqual bool = tt.a.Equal(tt.b)
		if isEqual != tt.isEqual {
			t.Errorf("%v.Equal(%v) want %v got %b", tt.a, tt.b, tt.isEqual, isEqual)
		}
	}
}
