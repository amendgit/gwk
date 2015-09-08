package ggk

// 8-bit type for an alpha value. 0xff is 100% opaque, 0x00 is 100% transparent.
// type Alpha uint8

// 32-bit ARGB color value, not permultiplied. The color components are alwarys in
// a known order. This is different from PMColor, which has its bytes in a configuration
// dependent order, to match the format of kARGB32 bitmaps, Color is the type used to
// specify colors in Paint and in gradient.
type Color uint32

// common colors
const (
	Color_AlphaTransparent = 0x00 // transparent Alpha value
	Color_AlphaOpaque      = 0xff // opaque Alpha value

	Color_Transparent = 0x00000000 // transparent Color value

	Color_Black  = 0xff000000 // black Color value
	Color_Dkgray = 0xff444444 // dark gray Color value
	Color_Gray   = 0xff888888 // gray Color value
	Color_Ltgray = 0xffcccccc // light gray Color value
	Color_White  = 0xffffffff // white Color value

	Color_Red     = 0xffff0000 // red Color value
	Color_Green   = 0xff00ff00 // green Color value
	Color_Blue    = 0xff0000ff // blue Color value
	Color_Yellow  = 0xffffff00 // yellow Color value
	Color_Cyan    = 0xff00ffff // cyan Color value
	Color_Magenta = 0xffff00ff // magenta Color value
)

// Return a Color from 8-bit component values.
func ColorWithARGB(a, r, g, b uint8) Color {
	return Color((uint32(a) << 24) | (uint32(r) << 16) | (uint32(g) << 8) | uint32(b))
}

// Return a Color value from 8-bit component values, with an implied value
// of 0xff for alpha (fully opaque)
func ColorWithRGB(r, g, b uint8) Color {
	return ColorWithARGB(0xff, r, g, b)
}

// Return the alpha byte from a Color value
func (color Color) Alpha() uint8 {
	return uint8((color >> 24) & 0xff)
}

// Return the red byte from a Color value
func (color Color) Red() uint8 {
	return uint8((color >> 16) & 0xff)
}

// Return the green byte from a Color value
func (color Color) Green() uint8 {
	return uint8((color >> 8) & 0xff)
}

// Return the blue byte from a Color value
func (color Color) Blue() uint8 {
	return uint8(color & 0xff)
}

// Return the alpha red green blue bytes in order from a Color value
func (color Color) ARGB() (uint8, uint8, uint8, uint8) {
	var a, r, g, b uint8
	a = uint8((color >> 24) & 0xff)
	r = uint8((color >> 16) & 0xff)
	g = uint8((color >> 8) & 0xff)
	b = uint8(color & 0xff)
	return a, r, g, b
}

// Set the alpha byte to a Color value.
func (color *Color) SetAlpha(alpha uint8) {
	*color = Color(uint32(*color&0x00ffffff) | (uint32(alpha) << 24))
}
