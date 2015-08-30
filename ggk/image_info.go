package ggk

// Alpha types
// Describe how to interpret the alpha compoent of a pixel.
const (
	AlphaUnknown = iota

	// All pixels are stored as opaque. This differs slightly from kIgnore in
	// that kOpaque has correct "Opaque" values stored in the pixels, while
	// kIgnore may not, but in both cases the caller should treat the pixels
	// as opaque.
	AlphaOpaque

	// All pixels have their alpha premultiplied in their color components.
	// This is the natural format for the rendering target pixels.
	AlphaPremul

	// All pixels have their color components stroed without any regard to the
	// alpha. e.g. this is the default configuration for PNG images.
	//
	// This alpha-type is ONLY supported for input images. Rendering cannot
	// generate this on output.
	AlphaUnpremul

	AlphaLast = AlphaUnpremul
)

// Color types
// Describes how to interpret the components of a pixel.
// ColorTypeN32 is an alias for whichever 32bit ARGB format is the "native"
// form for blitters. Use this if you don't hava a swizzle preference
// for 32bit pixels.
const (
	ColorTypeUnknown = iota
	ColorTypeAlpha8
	ColorTypeRGB565
	ColorTypeARGB4444
	ColorTypeRGBA8888
	ColorTypeBGRA8888
	ColorTypeIndex8
	ColorTypeGray8
	ColorTypeLast = ColorTypeGray8
)

func ColorTypeBytesPerPixel(colorType uint) uint {
	// TODO: const?
	var bytesPerPixel = []uint{
		0, // Unknown
		1, // Alpha8
		2, // RGB565
		2, // ARGB4444
		4, // RGBA8888
		4, // BGRA8888
		1, // Index8
		1, // Gray8
	}

	if colorType >= uint(len(bytesPerPixel)) {
		return 0
	}

	return bytesPerPixel[colorType]
}

func ColorTypeMinRowBytes(colorType uint, width uint) uint {
	return width * ColorTypeBytesPerPixel(colorType)
}

func ColorTypeIsVaild(value uint) bool {
	return value <= ColorTypeLast
}

func AlphaTypeIsOpaque(alphaType uint) bool {
	return alphaType == AlphaOpaque
}

func AlphaTypeIsValid(value uint) bool {
	return value <= AlphaLast
}

func ColorTypeComputeOffset(colorType uint, x, y int, rowBytes uint) uint {
	var shift uint = 0

	switch ColorTypeBytesPerPixel(colorType) {
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
func ColorTypeValidateAlphaType(colorType uint, alphaType uint) (canonical uint) {
	return 0
}

// YUV color space
// Describes the color space a YUV pixel
const (
	// Standard JPEG color space.
	YuvColorSpaceJPEG = iota
	// SDTV standard Rec. 601 color space. Uses "studio swing" [16, 245] color
	// range. See http://en.wikipedia.org/wiki/Rec._601 for details.
	YUVColorSpaceRec601
	// HDTV standard Rec. 709 color space. Uses "studio swing" [16, 235] color
	// range. See http://en.wikipedia.org/wiki/Rec._709 for details.
	YUVColorSpaceRec709

	YUVColorSpaceLast = YUVColorSpaceRec709
)

// Color profile type
const (
	ColorProfileTypeLinear = iota
	ColorProfileTypeSRGB
	ColorProfileTypeLast = ColorProfileTypeSRGB
)

type ImageInfo struct {
	width  int
	height int

	colorType        uint
	alphaType        uint
	colorProfileType uint
}

func (imageInfo *ImageInfo) Width() int {
	return imageInfo.width
}

func (imageInfo *ImageInfo) Height() int {
	return imageInfo.height
}
