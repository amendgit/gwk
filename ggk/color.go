package ggk

import "errors"

// 8-bit type for an alpha value. 0xff is 100% opaque, 0x00 is 100% transparent.
// type Alpha uint8

// Color is 32-bit ARGB color value, not permultiplied. The color components are
// alwarys in a known order. This is different from PMColor, which has its bytes
// in a configuration dependent order, to match the format of kARGB32 bitmaps,
// Color is the type used to specify colors in Paint and in gradient.
type Color uint32

// common colors
const (
	KColorAlphaTransparent Color = 0x00 // transparent Alpha value
	KColorAlphaOpaque            = 0xff // opaque Alpha value

	KColorTransparent = 0x00000000 // transparent Color value

	KColorBlack  = 0xff000000 // black Color value
	KColorDkgray = 0xff444444 // dark gray Color value
	KColorGray   = 0xff888888 // gray Color value
	KColorLtgray = 0xffcccccc // light gray Color value
	KColorWhite  = 0xffffffff // white Color value

	KColorRed     = 0xffff0000 // red Color value
	KColorGreen   = 0xff00ff00 // green Color value
	KColorBlue    = 0xff0000ff // blue Color value
	KColorYellow  = 0xffffff00 // yellow Color value
	KColorCyan    = 0xff00ffff // cyan Color value
	KColorMagenta = 0xffff00ff // magenta Color value
)

// ColorWithARGB return a Color from 8-bit component values.
func ColorWithARGB(a, r, g, b uint8) Color {
	return Color((uint32(a) << 24) | (uint32(r) << 16) | (uint32(g) << 8) | uint32(b))
}

// ColorWithRGB return a Color value from 8-bit component values, with an
// implied value of 0xff for alpha (fully opaque)
func ColorWithRGB(r, g, b uint8) Color {
	return ColorWithARGB(0xff, r, g, b)
}

// Alpha return the alpha byte from a Color value
func (color Color) Alpha() uint8 {
	return uint8((color >> 24) & 0xff)
}

// Red return the red byte from a Color value
func (color Color) Red() uint8 {
	return uint8((color >> 16) & 0xff)
}

// Green return the green byte from a Color value
func (color Color) Green() uint8 {
	return uint8((color >> 8) & 0xff)
}

// Blue return the blue byte from a Color value
func (color Color) Blue() uint8 {
	return uint8(color & 0xff)
}

// ARGB return the alpha red green blue bytes in order from a Color value
func (color Color) ARGB() (uint8, uint8, uint8, uint8) {
	var a, r, g, b uint8
	a = uint8((color >> 24) & 0xff)
	r = uint8((color >> 16) & 0xff)
	g = uint8((color >> 8) & 0xff)
	b = uint8((color >> 0) & 0xff)
	return a, r, g, b
}

// SetAlpha set the alpha byte to a Color value.
func (color *Color) SetAlpha(alpha uint8) {
	*color = Color(uint32(*color&0x00ffffff) | (uint32(alpha) << 24))
}

// PremultipliedColor is 32-bit ARGB color value, premultiplied. The byte order
// for this value is configuration dependent, matching the format of kARGB32
//  bitmaps. This is different from Color, which is nonpremultipled, and is
// always in the same byte order.
type PremulColor uint32

// ErrARGBIsNotPremultipled represent the ARGB is not premultiplied error.
var ErrARGBIsNotPremultipled = errors.New("ARGB is not premultiplied")

// PremulColorFromARGB32 pack the components into a PremulColor,
// checking (in the debug version) that the component are 0..255, and are
// already premultiplied (i.e. alpha >= color)
func PremulColorFromARGB32(a, r, g, b uint8) (PremulColor, error) {
	if r > a || g > a || b > a {
		return 0, ErrARGBIsNotPremultipled
	}
	return PremulColor((a << KARGB32ShiftA) | (r << KARGB32ShiftR) |
		(g << KARGB32ShiftG) | (b << KARGB32ShiftB)), nil
}

// PremultiplyARGB return a PremultipliedColor value from unpremultiplied 8-bit
// component values.
func PremultiplyARGB(a, r, g, b uint8) (PremulColor, error) {
	if a != 255 {
		r = MulDiv255Round(uint16(r), uint16(a))
		g = MulDiv255Round(uint16(r), uint16(a))
		b = MulDiv255Round(uint16(r), uint16(a))
	}
	return PremulColorFromARGB32(a, r, g, b)
}

// PremulColor return a PremultipliedColor value from a Color value. This
// is done by mutiplying the color components by the color's alpha, and by
// arranging the bytes in a configuration dependent order, to match the format
// of ARGB32 bitmaps.
func PremultiplyColor(c Color) (PremulColor, error) {
	var a, r, g, b uint8 = c.ARGB()
	return PremultiplyARGB(a, r, g, b)
}
