package ggk

// Alpha types
// Describe how to interpret the alpha compoent of a pixel.
type AlphaType int

const (
	AlphaTypeUnknown AlphaType = iota

	// All pixels are stored as opaque. This differs slightly from kIgnore in
	// that kOpaque has correct "Opaque" values stored in the pixels, while
	// kIgnore may not, but in both cases the caller should treat the pixels
	// as opaque.
	AlphaTypeOpaque

	// All pixels have their alpha premultiplied in their color components.
	// This is the natural format for the rendering target pixels.
	AlphaTypePremul

	// All pixels have their color components stroed without any regard to the
	// alpha. e.g. this is the default configuration for PNG images.
	//
	// This alpha-type is ONLY supported for input images. Rendering cannot
	// generate this on output.
	AlphaTypeUnpremul

	AlphaTypeLastEnum = AlphaTypeUnpremul
)

func AlphaTypeIsOpaque(alphaType AlphaType) bool {
	return alphaType == AlphaTypeOpaque
}

func AlphaTypeIsValid(value AlphaType) bool {
	return value >= 0 && value <= AlphaTypeLastEnum
}

// Color types
// Describes how to interpret the components of a pixel.
// ColorTypeN32 is an alias for whichever 32bit ARGB format is the "native"
// form for blitters. Use this if you don't hava a swizzle preference
// for 32bit pixels.
type ColorType int

const (
	ColorTypeUnknown ColorType = iota
	ColorTypeAlpha8
	ColorTypeRGB565
	ColorTypeARGB4444
	ColorTypeRGBA8888
	ColorTypeBGRA8888
	ColorTypeIndex8
	ColorTypeGray8

	ColorTypeLastEnum = ColorTypeGray8
)

func ColorTypeBytesPerPixel(colorType ColorType) int {
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

	if colorType >= ColorType(len(bytesPerPixel)) {
		return 0
	}

	return bytesPerPixel[colorType]
}

func ColorTypeMinRowBytes(colorType ColorType, width int) int {
	return width * ColorTypeBytesPerPixel(colorType)
}

func ColorTypeIsVaild(value ColorType) bool {
	return value <= ColorTypeLastEnum
}

func ColorTypeComputeOffset(colorType ColorType, x, y int, rowBytes int) int {
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

	return y*rowBytes + x<<shift
}

// Return true if alphaType is supported by colorType. If there is a canonical
// alphaType for this colorType, return it in canonical.
func ColorTypeValidateAlphaType(colorType ColorType, alphaType AlphaType) AlphaType {
	return 0
}

// YUV color space
// Describes the color space a YUV pixel
type YUVColorSpace int

const (
	// Standard JPEG color space.
	YUVColorSpaceJPEG YUVColorSpace = iota
	// SDTV standard Rec. 601 color space. Uses "studio swing" [16, 245] color
	// range. See http://en.wikipedia.org/wiki/Rec._601 for details.
	YUVColorSpaceRec601
	// HDTV standard Rec. 709 color space. Uses "studio swing" [16, 235] color
	// range. See http://en.wikipedia.org/wiki/Rec._709 for details.
	YUVColorSpaceRec709

	YUVColorSpaceLastEnum = YUVColorSpaceRec709
)

// Color profile type
type ColorProfileType int

const (
	ColorProfileTypeLinear ColorProfileType = iota
	ColorProfileTypeSRGB

	ColorProfileTypeLastEnum = ColorProfileTypeSRGB
)

// Describe an image's dimensions and pixel type.
// Used for both src images and render-targets (surfaces).
type ImageInfo struct {
	Width  int
	Height int

	ColorType   ColorType
	AlphaType   AlphaType
	ProfileType ColorProfileType
}

func (imageInfo *ImageInfo) BytesPerPixel() int {
	return ColorTypeBytesPerPixel(imageInfo.ColorType)
}

func New(width, height int, colorType ColorType, alphaType AlphaType,
	profileType ColorProfileType) *ImageInfo {
	var imageInfo = &ImageInfo{
		Width:       width,
		Height:      height,
		ColorType:   colorType,
		AlphaType:   alphaType,
		ProfileType: profileType,
	}

	return imageInfo
}

func NewN32() {

}

func NewN32Premul() {

}

func NewA8() {

}
