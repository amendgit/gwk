package ggk

import (
	"fmt"
)

// Alpha types
// Describe how to interpret the alpha compoent of a pixel.
type AlphaType int

const (
	AlphaType_Unknown AlphaType = iota

	// All pixels are stored as opaque. This differs slightly from kIgnore in
	// that kOpaque has correct "Opaque" values stored in the pixels, while
	// kIgnore may not, but in both cases the caller should treat the pixels
	// as opaque.
	AlphaType_Opaque

	// All pixels have their alpha premultiplied in their color components.
	// This is the natural format for the rendering target pixels.
	AlphaType_Premul

	// All pixels have their color components stroed without any regard to the
	// alpha. e.g. this is the default configuration for PNG images.
	//
	// This alpha-type is ONLY supported for input images. Rendering cannot
	// generate this on output.
	AlphaType_Unpremul

	AlphaType_LastEnum = AlphaType_Unpremul
)

func (at AlphaType) IsOpaque() bool {
	return at == AlphaType_Opaque
}

func (at AlphaType) IsValid() bool {
	return at >= 0 && at <= AlphaType_LastEnum
}

// Color types
// Describes how to interpret the components of a pixel.
// ColorTypeN32 is an alias for whichever 32bit ARGB format is the "native"
// form for blitters. Use this if you don't hava a swizzle preference
// for 32bit pixels.
type ColorType int

const (
	ColorType_Unknown ColorType = iota
	ColorType_Alpha8
	ColorType_RGB565
	ColorType_ARGB4444
	ColorType_RGBA8888
	ColorType_BGRA8888
	ColorType_Index8
	ColorType_Gray8

	ColorTypeLastEnum = ColorType_Gray8
)

func (ct ColorType) BytesPerPixel() int {
	// TODO: const?
	var bytesPerPixel = []int{
		0, // Unknown
		1, // Alpha8
		2, // RGB565
		2, // ARGB4444
		4, // RGBA8888
		4, // BGRA8888
		1, // Index8
		1, // Gray8
	}

	if ct < 0 || int(ct) >= len(bytesPerPixel) {
		return 0
	}

	return bytesPerPixel[ct]
}

func (ct ColorType) MinRowBytes(width int) int {
	return width * ct.BytesPerPixel()
}

func (ct ColorType) IsVaild() bool {
	return ct >= 0 && ct <= ColorTypeLastEnum
}

func (ct ColorType) ComputeOffset(x, y int, rowBytes uint) uint {
	var shift uint = 0

	switch ct.BytesPerPixel() {
	case 4:
		shift = 2
	case 2:
		shift = 1
	case 1:
		shift = 0
	default:
		return 0
	}

	return uint(y)*rowBytes + uint(x)<<shift
}

// Return true if alphaType is supported by colorType. If there is a canonical
// alphaType for this colorType, return it in canonical.
func (ct ColorType) ValidateAlphaType(alphaType AlphaType) AlphaType {
	return 0
}

// YUV color space
// Describes the color space a YUV pixel
type YUVColorSpace int

const (
	// Standard JPEG color space.
	YUVColorSpace_JPEG YUVColorSpace = iota
	// SDTV standard Rec. 601 color space. Uses "studio swing" [16, 245] color
	// range. See http://en.wikipedia.org/wiki/Rec._601 for details.
	YUVColorSpace_Rec601
	// HDTV standard Rec. 709 color space. Uses "studio swing" [16, 235] color
	// range. See http://en.wikipedia.org/wiki/Rec._709 for details.
	YUVColorSpace_Rec709

	YUVColorSpace_LastEnum = YUVColorSpace_Rec709
)

// Color profile type
type ColorProfileType int

const (
	ColorProfileType_Linear ColorProfileType = iota
	ColorProfileType_SRGB
	ColorProfileType_LastEnum = ColorProfileType_SRGB
)

func (pt ColorProfileType) IsValid() bool {
	return pt >= 0 && pt <= ColorProfileType_LastEnum
}

// Describe an image's dimensions and pixel type.
// Used for both src images and render-targets (surfaces).
type ImageInfo struct {
	width  int
	height int

	colorType   ColorType
	alphaType   AlphaType
	profileType ColorProfileType
}

func NewImageInfo(width, height int, colorType ColorType, alphaType AlphaType, profileType ColorProfileType) *ImageInfo {
	var imageInfo = &ImageInfo{
		width:       width,
		height:      height,
		colorType:   colorType,
		alphaType:   alphaType,
		profileType: profileType,
	}

	return imageInfo
}

func NewImageInfoN32(width, height int, alphaType AlphaType, profileType ColorProfileType) *ImageInfo {
	return NewImageInfo(width, height, ColorType_RGBA8888, alphaType, profileType)
}

func NewImageInfoA8(width, height int) *ImageInfo {
	return NewImageInfo(width, height, ColorType_Alpha8, AlphaType_Premul, ColorProfileType_Linear)
}

func NewImageInfoUnknown(width, height int) *ImageInfo {
	return NewImageInfo(width, height, ColorType_Unknown, AlphaType_Unknown, ColorProfileType_Linear)
}

func (imageInfo *ImageInfo) Width() int {
	return imageInfo.width
}

func (imageInfo *ImageInfo) Height() int {
	return imageInfo.height
}

func (imageInfo *ImageInfo) ColorType() ColorType {
	return imageInfo.colorType
}

func (imageInfo *ImageInfo) AlphaType() AlphaType {
	return imageInfo.alphaType
}

func (imageInfo *ImageInfo) ProfileType() ColorProfileType {
	return imageInfo.profileType
}

func (imageInfo *ImageInfo) IsEmpty() bool {
	return imageInfo.width <= 0 || imageInfo.height <= 0
}

func (imageInfo *ImageInfo) IsOpaque() bool {
	return imageInfo.alphaType.IsOpaque()
}

func (imageInfo *ImageInfo) IsLinear() bool {
	return imageInfo.profileType == ColorProfileType_Linear
}

func (imageInfo *ImageInfo) IsSRGB() bool {
	return imageInfo.profileType == ColorProfileType_SRGB
}

func (imageInfo *ImageInfo) ComputeOffset(x, y int, rowBytes uint) (uint, error) {
	if uint(x) >= uint(imageInfo.width) || uint(y) >= uint(imageInfo.height) {
		return 0, fmt.Errorf("OOR: ggk.ImageInfo(0x%x).ComputeOffset(%d, %d, %d)", imageInfo, x, y, rowBytes)
	}

	return imageInfo.colorType.ComputeOffset(x, y, rowBytes), nil
}

func (imageInfo *ImageInfo) Equal(other *ImageInfo) bool {
	var equal = false

	equal = (imageInfo.colorType == other.colorType)
	equal = equal && (imageInfo.alphaType == other.alphaType)
	equal = equal && (imageInfo.profileType == other.profileType)
	equal = equal && (imageInfo.width == other.width)
	equal = equal && (imageInfo.height == other.height)

	return equal
}

func (imageInfo *ImageInfo) BytesPerPixel() int {
	return imageInfo.colorType.BytesPerPixel()
}

func (imageInfo *ImageInfo) MinRowBytes64() uint64 {
	var minRowBytes64 uint64 = uint64(imageInfo.width) * uint64(imageInfo.BytesPerPixel())
	return minRowBytes64
}

func (imageInfo *ImageInfo) MinRowBytes() uint {
	return uint(imageInfo.MinRowBytes64())
}

func (imageInfo *ImageInfo) ValidRowBytes(rowBytes uint) bool {
	return uint64(rowBytes) >= imageInfo.MinRowBytes64()
}

func (imageInfo *ImageInfo) SafeSize64(rowBytes uint) uint64 {
	if imageInfo.height == 0 {
		return 0
	}

	return uint64(imageInfo.height-1)*uint64(rowBytes) + uint64(imageInfo.width*imageInfo.BytesPerPixel())
}

func (imageInfo *ImageInfo) SafeSize(rowBytes uint) uint {
	var size uint64 = imageInfo.SafeSize64(rowBytes)
	if size != uint64(uint(size)) {
		return 0
	}
	return uint(size)
}
