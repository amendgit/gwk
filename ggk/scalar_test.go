package ggk_test

import (
	"gwk/ggk"
	"testing"
)

var scalarBasicTests = []struct {
	x        ggk.Scalar
	floor    ggk.Scalar
	floorInt int
	ceil     ggk.Scalar
	ceilInt  int
	round    ggk.Scalar
	roundInt int
}{
	{0.0, 0.0, 0, 0.0, 0, 0.0, 0},
	{2.5, 2.0, 2, 3.0, 3, 3.0, 3},
	{-2.5, -3.0, -3, -2.0, -2, -3.0, -3},
}

func TestBasicScalar(t *testing.T) {
	for _, tt := range scalarBasicTests {
		floor := ggk.ScalarFloor(tt.x)
		if floor != tt.floor {
			t.Errorf("ScalarFloor(%v) want %v got %v", tt.x, tt.floor, floor)
		}

		floorInt := ggk.ScalarFloorToInt(tt.x)
		if floorInt != tt.floorInt {
			t.Errorf("ScalarFloorInt(%v) want %v got %v", tt.x, tt.floorInt, floorInt)
		}

		ceil := ggk.ScalarCeil(tt.x)
		if ceil != tt.ceil {
			t.Errorf("ScalarCeil(%v) want %v got %v", tt.x, tt.ceil, ceil)
		}

		ceilInt := ggk.ScalarCeilToInt(tt.x)
		if ceilInt != tt.ceilInt {
			t.Errorf("ScalarCeilInt(%v) want %v got %v", tt.x, tt.ceilInt, ceilInt)
		}

		round := ggk.ScalarRound(tt.x)
		if round != tt.round {
			t.Errorf("ScalarRound(%v) want %v got %v", tt.x, tt.round, round)
		}

		roundInt := ggk.ScalarRoundToInt(tt.x)
		if roundInt != tt.roundInt {
			t.Errorf("ScalarRound(%v) want %v got %v", tt.x, tt.roundInt, roundInt)
		}
	}
}

var scalarPropertyTests = []struct {
	x     ggk.Scalar
	isNaN bool
	isFin bool
}{
	{0, false, true},
	{1, false, true},
	{ggk.ScalarNaN(), true, false},
	{ggk.ScalarInfinity(), false, false},
}

func TestScalarProperty(t *testing.T) {
	for _, tt := range scalarPropertyTests {
		var (
			isNaN = ggk.ScalarIsNaN(tt.x)
			isFin = ggk.ScalarIsFinite(tt.x)
		)

		if isNaN != tt.isNaN {
			t.Errorf("ScalarIsNaN(%v) want %v get %v", tt.x, tt.isNaN, isNaN)
		}

		if isFin != tt.isFin {
			t.Errorf("ScalarIsFinite(%v) want %v get %v", tt.x, tt.isFin, isFin)
		}
	}
}
