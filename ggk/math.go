package ggk

// Return a*b/255, rounding any fractional bits.
// Only valid if a and b are unsigned and <= 0x7fff
func MulDiv255Round(a, b uint16) uint8 {
	var prod uint8 = uint8(a*b + 128)
	return (prod + (prod >> 8)) >> 8
}
