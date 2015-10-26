package ggk_test

import (
	"gwk/ggk"
	"testing"
)

var sizeEqualTests = []struct {
	a       ggk.Size
	b       ggk.Size
	isEqual bool
}{
	{
		ggk.MakeSize(0, 0),
		ggk.MakeSize(0, 0),
		true,
	},
}

func TestSizeEqual(t *testing.T) {
	for _, tt := range sizeEqualTests {
		var isEqual bool = tt.a.Equal(tt.b)
		if isEqual != tt.isEqual {
			t.Errorf("%v.Equal(%v) want %v got %v", tt.a, tt.b, tt.isEqual, isEqual)
		}
	}
}
