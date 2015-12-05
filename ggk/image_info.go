package ggk

import (
	"errors"
	"fmt"
)

// AlphaType describe how to interpret the alpha component of a pixel.
type AlphaType int

const (
	// AlphaTypeUnknown represent unknown alpha type value.
	KAlphaTypeUnknown AlphaType = iota

	// AlphaTypeOpaque all pixels are stored as opaque. This differs slightly from
	// kIgnore in that kOpaque has correct "Opaque" values stored in the pixels,
	// while kIgnore may not, but in both cases the caller should treat the pixels
	// as opaque.
	KAlphaTypeOpaque

	// AlphaTypePremul all pixels have their alpha premultiplied in their color
	// components. This is the natural format for the rendering target pixels.
	KAlphaTypePremul

	// AlphaTypeUnpremul pixels have their color components stroed without any
	// regard to the alpha. e.g. this is the default configuration for PNG images.
	//
	// This alpha-type is ONLY supported for input images. Rendering cannot
	// generate this on output.
	KAlphaTypeUnpremul

	// AlphaTypeLastEnum is the
	KAlphaTypeLastEnum = KAlphaTypeUnpremul
)

// IsOpaque return true if AlphaType value is opaque.
func (at AlphaType) IsOpaque() bool {
	return at == KAlphaTypeOpaque
}

// IsValid return true if AlphaType value is vaild.
func (at AlphaType) IsValid() bool {
	return at >= 0 && at <= KAlphaTypeLastEnum
}

// ColorType describes how to interpret the components of a pixel.
// ColorTypeN32 is an alias for whichever 32bit ARGB format is the "native"
// form for blitters. Use this if you don't hava a swizzle preference
// for 32bit pixels.
type ColorType int

const (
	KColorTypeUnknown ColorType = iota
	KColorTypeAlpha8
	KColorTypeRGB565
	KColorTypeARGB4444
	KColorTypeRGBA8888
	KColorTypeBGRA8888
	KColorTypeIndex8
	KColorTypeGray8

	KColorTypeLastEnum = KColorTypeGray8
)

func (ct ColorType) BytesPerPixel() int {
	var bytesPerPixel = [...]int{
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
	return ct >= 0 && ct <= KColorTypeLastEnum
}

func (ct ColorType) ComputeOffset(x, y int, rowBytes uint) uint {
	if x < 0 || y < 0 || (!ct.IsVaild()) || (ct == KColorTypeUnknown) ||
		(rowBytes%uint(ct.BytesPerPixel()) != 0) {
		return 0
	}

	return uint(y)*rowBytes + uint(x*ct.BytesPerPixel())
}

var ErrAlphaTypeCanNotCanonical = errors.New("color type can't be canonical")

// Return true if alphaType is supported by colorType. If there is a canonical
// alphaType for this colorType, return it in canonical.
func (ct ColorType) ValidateAlphaType(alphaType AlphaType) (canonical AlphaType, err error) {
	switch ct {
	case KColorTypeUnknown:
		alphaType = KAlphaTypeUnknown

	case KColorTypeAlpha8:
		if alphaType == KAlphaTypeUnpremul {
			alphaType = KAlphaTypePremul
		}

		fallthrough

	case KColorTypeIndex8, KColorTypeARGB4444, KColorTypeRGBA8888,
		KColorTypeBGRA8888:
		if alphaType == KAlphaTypeUnknown {
			return KAlphaTypeUnknown, ErrAlphaTypeCanNotCanonical
		}

	case KColorTypeGray8, KColorTypeRGB565:
		alphaType = KAlphaTypeOpaque

	default:
		return KAlphaTypeUnknown, ErrAlphaTypeCanNotCanonical
	}

	return alphaType, nil
}

// YUVColorSpace describes the color space a YUV pixel
type YUVColorSpace int

const (
	// Standard JPEG color space.
	KYUVColorSpaceJPEG YUVColorSpace = iota
	// SDTV standard Rec. 601 color space. Uses "studio swing" [16, 245] color
	// range. See http://en.wikipedia.org/wiki/Rec._601 for details.
	KYUVColorSpaceRec601
	// HDTV standard Rec. 709 color space. Uses "studio swing" [16, 235] color
	// range. See http://en.wikipedia.org/wiki/Rec._709 for details.
	KYUVColorSpaceRec709

	KYUVColorSpaceLastEnum = KYUVColorSpaceRec709
)

// Color profile type
type ColorProfileType int

const (
	KColorProfileTypeLinear ColorProfileType = iota
	KColorProfileTypeSRGB
	KColorProfileTypeLastEnum = KColorProfileTypeSRGB
)

func (pt ColorProfileType) IsValid() bool {
	return pt >= 0 && pt <= KColorProfileTypeLastEnum
}

// Describe an image's dimensions and pixel type.
// Used for both src images and render-targets (surfaces).
type ImageInfo struct {
	width  Scalar
	height Scalar

	colorType   ColorType
	alphaType   AlphaType
	profileType ColorProfileType
}

func NewImageInfo(width, height Scalar, colorType ColorType, alphaType AlphaType,
	profileType ColorProfileType) *ImageInfo {
	var imageInfo = &ImageInfo{
		width:       width,
		height:      height,
		colorType:   colorType,
		alphaType:   alphaType,
		profileType: profileType,
	}

	return imageInfo
}

func NewImageInfoN32(width, height Scalar, alphaType AlphaType, profileType ColorProfileType) *ImageInfo {
	return NewImageInfo(width, height, KColorTypeN32, alphaType, profileType)
}

func NewImageInfoN32Premul(width, height Scalar, profileType ColorProfileType) *ImageInfo {
	return NewImageInfo(width, height, KColorTypeN32, KAlphaTypePremul, profileType)
}

func NewImageInfoA8(width, height Scalar) *ImageInfo {
	return NewImageInfo(width, height, KColorTypeAlpha8, KAlphaTypePremul, KColorProfileTypeLinear)
}

func NewImageInfoUnknown(width, height Scalar) *ImageInfo {
	return NewImageInfo(width, height, KColorTypeUnknown, KAlphaTypeUnknown, KColorProfileTypeLinear)
}

func (ii *ImageInfo) CloneWH(width, height Scalar) *ImageInfo {
	var newInfo ImageInfo
	newInfo = *ii
	newInfo.width = width
	newInfo.height = height
	return &newInfo
}

func (ii *ImageInfo) Width() Scalar {
	return ii.width
}

func (ii *ImageInfo) Height() Scalar {
	return ii.height
}

func (ii *ImageInfo) ColorType() ColorType {
	return ii.colorType
}

func (ii *ImageInfo) SetColorType(colorType ColorType) {
	ii.colorType = colorType
}

func (ii *ImageInfo) AlphaType() AlphaType {
	return ii.alphaType
}

func (ii *ImageInfo) SetAlphaType(alphaType AlphaType) {
	ii.alphaType = alphaType
}

func (ii *ImageInfo) ProfileType() ColorProfileType {
	return ii.profileType
}

func (ii *ImageInfo) IsValid() bool {
	if ii.width < 0 || ii.height < 0 {
		return false
	}

	if !ii.colorType.IsVaild() {
		return false
	}

	if !ii.alphaType.IsValid() {
		return false
	}

	return true
}

func (ii *ImageInfo) IsEmpty() bool {
	return ii.width <= 0 || ii.height <= 0
}

func (ii *ImageInfo) IsOpaque() bool {
	return ii.alphaType.IsOpaque()
}

func (ii *ImageInfo) IsLinear() bool {
	return ii.profileType == KColorProfileTypeLinear
}

func (ii *ImageInfo) IsSRGB() bool {
	return ii.profileType == KColorProfileTypeSRGB
}

func (ii *ImageInfo) ComputeOffset(x, y int, rowBytes uint) (uint, error) {
	if uint(x) >= uint(ii.width) || uint(y) >= uint(ii.height) {
		return 0, fmt.Errorf("OOR: ggk.ImageInfo(0x%x).ComputeOffset(%d, %d, %d)",
			ii, x, y, rowBytes)
	}

	return ii.colorType.ComputeOffset(x, y, rowBytes), nil
}

func (ii *ImageInfo) Equal(other *ImageInfo) bool {
	var equal = false

	equal = (ii.colorType == other.colorType)
	equal = equal && (ii.alphaType == other.alphaType)
	equal = equal && (ii.profileType == other.profileType)
	equal = equal && (ii.width == other.width)
	equal = equal && (ii.height == other.height)

	return equal
}

func (ii *ImageInfo) BytesPerPixel() int {
	return ii.colorType.BytesPerPixel()
}

func (ii *ImageInfo) MinRowBytes64() int64 {
	var minRowBytes64 int64 = int64(ii.width) * int64(ii.BytesPerPixel())
	return minRowBytes64
}

func (ii *ImageInfo) MinRowBytes() int {
	return int(ii.MinRowBytes64())
}

func (ii *ImageInfo) ValidRowBytes(rowBytes int) bool {
	return int64(rowBytes) >= ii.MinRowBytes64()
}

func (ii *ImageInfo) SafeSize64(rowBytes int) uint64 {
	if ii.height == 0 {
		return 0
	}

	return uint64(ii.height-1)*uint64(rowBytes) +
		uint64(int(ii.width)*ii.BytesPerPixel())
}

func (ii *ImageInfo) SafeSize(rowBytes int) uint {
	var size uint64 = ii.SafeSize64(rowBytes)
	if size != uint64(uint(size)) {
		return 0
	}
	return uint(size)
}

// ReadPixelsHlp is helper to package and trim the parameters passed to
// ReadPixels()
type tReadPixelsHlp struct {
	Pixels   []byte
	RowBytes int
	Info     *ImageInfo
	X        Scalar
	Y        Scalar
}

func newReadPixelsHlp(info *ImageInfo, pixels []byte, rowBytes int, x, y Scalar) *tReadPixelsHlp {
	var hlp = new(tReadPixelsHlp)
	hlp.Info = info
	hlp.Pixels = pixels
	hlp.RowBytes = rowBytes
	hlp.X = x
	hlp.Y = y
	return hlp
}

// On nil, may have modified its fields (except RowBytes) to make it a legal
// subset of the specified src width/height.
//
// On error, leaves self unchanged, but indicates that it does not overlap src,
// or is not valid (e.g. bad Info) for ReadPixels().
func (h *tReadPixelsHlp) Trim(srcWidth, srcHeight Scalar) error {
	var ct = h.Info.ColorType()
	if ct == KColorTypeUnknown || ct == KColorTypeIndex8 {
		return errors.New("tReadPixelsHlp Trim: bad color type")
	}

	if h.Pixels == nil || h.RowBytes < h.Info.MinRowBytes() {
		return errors.New("tReadPixelsHlp Trim: bad pixels or rowBytes")
	}

	if h.Info.Width() == 0 || h.Info.Width() == 0 {
		return errors.New("tReadPixelsHlp Trim: bad width or height")
	}

	var x, y = h.X, h.Y
	var srcRect = MakeRect(x, y, h.Info.Width(), h.Info.Height())
	if !srcRect.IntersectXYWH(0, 0, srcWidth, srcHeight) {
		return errors.New("tReadPixelsHlp Trim: bad srcRect")
	}

	// if x or y are negative, then we have to adjust pixels.
	if x > 0 {
		x = 0
	}

	if y > 0 {
		y = 0
	}

	// x, y are either 0 or negative.
	var idx = int(0 - (y*Scalar(h.RowBytes) + x*Scalar(h.Info.BytesPerPixel())))

	// the intersect may have shrunk info's logical size.
	h.Pixels = h.Pixels[idx:]
	h.Info = h.Info.CloneWH(srcRect.Width(), srcRect.Height())
	h.X, h.Y = srcRect.X(), srcRect.Y()

	return nil
}
