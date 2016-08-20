package ggk

// Paint holds the style and color information about how to draw geometries, text
// and bitmaps.
type Paint struct {
	flags   uint16
	hinting uint8
}

// PaintHinting specifies the level of hinting to be performed. These names are
// taken from the Gnome/Cairo names for the same. They are translated into
// Freetype concepts the same as in cairo-ft-font.c:
// KPaintHintingNo     -> FT_LOAD_NO_HINTING
// KPaintHintingSlight -> FT_LOAD_TARGET_LIGHT
// KPaintHintingNormal -> <default, no option>
// KPaintHintingFull   -> <same as KPaintHintingNormal, unelss we are rendering
//                         subpixel glyphs, in which case TARGET_LCD or
//                         TARGET_LCD_V is used>
type PaintHinting int

const (
	KPaintHintingNo     = 0
	KPaintHintingSlight = 1
	KPaintHintingNormal = 2 // this is the default.
	KPaintHintingFull   = 3
)

func (p *Paint) Hinting() PaintHinting {
	return PaintHinting(p.hinting)
}

func (p *Paint) SetHinting(hinting PaintHinting) {
	p.hinting = uint8(hinting)
}

func (p *Paint) Looper() *DrawLooper {
	return p.looper
}

func (p *Paint) SetLooper(looper *DrawLooper) {
	p.looper = looper
}

type PaintFlags int

const (
	KPaintFlagAntiAlias          = 0x01
	KPaintFlagDither             = 0x04
	KPaintFlagUnderline          = 0x08
	KPaintFlagStrikeThruText     = 0x10
	KPaintFlagFakeBoldText       = 0x20
	KPaintFlagLinearText         = 0x40
	KPaintFlagSubpixelText       = 0x80
	KPaintFlagDevKernText        = 0x100
	KPaintFlagLCDRenderText      = 0x200
	KPaintFlagEmbeddedBitmapText = 0x400
	KPaintFlagAutoHinting        = 0x800
	KPaintFlagVerticalText       = 0x1000

	// hack for GDI -- do not use if you can help it when adding extra flags,
	// note that the flags member is specified with a bit-width and you'll have
	// expand it.
	KPaintFlagGenA8FromLCD = 0x2000

	KPaintFlagAllFlags = 0xFFFF
)

func (p *Paint) Flags() PaintFlags {
	return PaintFlags(p.flags)
}

func (p *Paint) SetFlags(flags PaintFlags) {
	p.flags = uint16(flags)
}

func (p *Paint) CanComputeFastBounds() bool {
	if p.Looper() != nil {
		return p.Looper().CanComputeFastBounds()
	}
	if p.ImageFilter() != nil && p.ImageFilter.CanComputeFastBounds() {
		return false
	}
	return !p.Rasterizer()
}