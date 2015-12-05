package ggk

// ColorTable holds an array PMColors (premultiplied 32-bit colors) used by
// 8-bit bitmaps, where the bitmap bytes are interpreted as indices into the color table.
//
// ColorTable is thread-safe.
type ColorTable struct {
	colors []PremulColor
	count  int
}

func NewColorTable(colors []PremulColor) *ColorTable {
	toimpl()
	return nil
}

func (ct *ColorTable) Count() int {
	return ct.count
}

func (ct *ColorTable) ReadColors() []PremulColor {
	return ct.colors
}

func (ct *ColorTable) At(index int) PremulColor {
	return ct.colors[index]
}
