package ggk_test

import (
	"gwk/ggk"
	"testing"
)

var atTests = []struct {
	at      ggk.AlphaType
	isValid bool
}{
	{ggk.KAlphaTypeUnknown, true},
	{ggk.KAlphaTypeOpaque, true},
	{ggk.KAlphaTypePremul, true},
	{ggk.KAlphaTypeUnpremul, true},
	{ggk.KAlphaTypeLastEnum, true},
	{ggk.AlphaType(-1), false},
	{ggk.AlphaType(777), false},
}

func TestAlphaType(t *testing.T) {
	for _, tt := range atTests {
		var isValid bool = tt.at.IsValid()
		if (!isValid && tt.isValid) || (isValid && !tt.isValid) {
			t.Errorf("AlphaType(%v).IsValid() want %v get %v", tt.at, tt.isValid, isValid)
		}
	}
}

var ctBytesPerPixelTests = []struct {
	ct       ggk.ColorType
	numBytes int
}{
	{ggk.KColorTypeUnknown, 0},
	{ggk.KColorTypeAlpha8, 1},
	{ggk.KColorTypeRGB565, 2},
	{ggk.KColorTypeARGB4444, 2},
	{ggk.KColorTypeRGBA8888, 4},
	{ggk.KColorTypeBGRA8888, 4},
	{ggk.KColorTypeIndex8, 1},
	{ggk.KColorTypeGray8, 1},
	{ggk.ColorType(1000), 0},
	{ggk.ColorType(-1), 0},
}

var ctComputeOffsetTests = []struct {
	ct       ggk.ColorType
	x, y     int
	rowBytes uint
	offset   uint
}{
	{ggk.KColorTypeRGBA8888, 0, 0, 0, 0},
	{ggk.KColorTypeRGBA8888, 0, 1, 4, 4},
	{ggk.KColorTypeRGBA8888, 0, 1, 8, 8},
	{ggk.KColorTypeRGBA8888, 1, 1, 8, 12},
	{ggk.KColorTypeRGBA8888, 1, 0, 8, 4},
	{ggk.KColorTypeRGBA8888, 0, 0, 0, 0},
	{ggk.KColorTypeRGBA8888, -1, 1, 8, 0},
	{ggk.ColorType(-1), 1, 1, 8, 0},
	{ggk.KColorTypeUnknown, 1, 1, 8, 0},
	{ggk.KColorTypeRGBA8888, 1, 1, 7, 0},
}

func TestColorType(t *testing.T) {
	for _, tt := range ctBytesPerPixelTests {
		var numBytes int = tt.ct.BytesPerPixel()
		if numBytes != tt.numBytes {
			t.Errorf("ColorType(%v).BytesPerPixel() want %v get %v", tt.ct, numBytes, tt.numBytes)
		}
	}

	for _, tt := range ctComputeOffsetTests {
		var offset uint = tt.ct.ComputeOffset(tt.x, tt.y, tt.rowBytes)
		if offset != tt.offset {
			t.Errorf("ColorType(%v).CoputeOffset(x:%v, y:%v, rowBytes:%v) want %v get %v",
				tt.ct, tt.x, tt.y, tt.rowBytes, tt.offset, offset)
		}
	}
}

var imageInfoEqTests = []struct {
	a, b    *ggk.ImageInfo
	isEqual bool
}{
	{
		ggk.NewImageInfo(100, 100, ggk.KColorTypeRGBA8888, ggk.KAlphaTypeOpaque, ggk.KColorProfileTypeLinear),
		ggk.NewImageInfo(100, 100, ggk.KColorTypeRGBA8888, ggk.KAlphaTypeOpaque, ggk.KColorProfileTypeLinear),
		true,
	},
	{
		ggk.NewImageInfo(100, 100, ggk.KColorTypeRGBA8888, ggk.KAlphaTypeOpaque, ggk.KColorProfileTypeLinear),
		ggk.NewImageInfo(100, 100, ggk.KColorTypeBGRA8888, ggk.KAlphaTypeOpaque, ggk.KColorProfileTypeLinear),
		false,
	},
}

var imageInfoMinRowBytesTests = []struct {
	imageInfo     *ggk.ImageInfo
	minRowBytes64 int64
	minRowBytes   int
}{
	{
		ggk.NewImageInfo(100, 100, ggk.KColorTypeN32, ggk.KAlphaTypeOpaque, ggk.KColorProfileTypeLinear),
		400,
		400,
	},
	{
		ggk.NewImageInfo(5000, 100, ggk.KColorTypeN32, ggk.KAlphaTypeOpaque, ggk.KColorProfileTypeLinear),
		20000,
		20000,
	},
}

var imageInfoSafeSizeTests = []struct {
	imageInfo  *ggk.ImageInfo
	rowBytes   int
	safeSize   uint
	safeSize64 uint64
}{
	{
		ggk.NewImageInfo(900, 601, ggk.KColorTypeN32, ggk.KAlphaTypeOpaque, ggk.KColorProfileTypeLinear),
		5000,
		3003600,
		3003600,
	},
}

func TestImageInfo(t *testing.T) {
	for _, tt := range imageInfoEqTests {
		var isEqual bool = tt.a.Equal(tt.b)
		if isEqual != tt.isEqual {
			t.Errorf("ImageInfo(%v).Equal(%v) want %v get %v", tt.a, tt.b, tt.isEqual, isEqual)
		}
	}

	for _, tt := range imageInfoMinRowBytesTests {
		var minRowBytes int = tt.imageInfo.MinRowBytes()

		if minRowBytes != tt.minRowBytes {
			t.Errorf("ImageInfo(%v).MinRowBytes() want %v get %v", tt.imageInfo, tt.minRowBytes, minRowBytes)
		}

		var minRowBytes64 int64 = tt.imageInfo.MinRowBytes64()
		if minRowBytes64 != tt.minRowBytes64 {
			t.Errorf("ImageInfo(%v).MinRowBytes64() want %v get %v", tt.imageInfo, tt.minRowBytes64, minRowBytes64)
		}
	}

	for _, tt := range imageInfoSafeSizeTests {
		var safeSize uint = tt.imageInfo.SafeSize(tt.rowBytes)
		if safeSize != tt.safeSize {
			t.Errorf("ImageInfo(%v).SafeSize(%v) want %v get %v", tt.imageInfo, tt.rowBytes, tt.safeSize, safeSize)
		}

		var safeSize64 uint64 = tt.imageInfo.SafeSize64(tt.rowBytes)
		if safeSize64 != tt.safeSize64 {
			t.Errorf("ImageInfo(%v).SafeSize64(%v) want %v get %v", tt.imageInfo, tt.rowBytes, tt.safeSize64, safeSize64)
		}
	}
}
