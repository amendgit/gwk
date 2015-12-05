package ggk

// Pixmap pairs ImageInfo with actual pixels and rowbytes. This class does not
// try to manage the lifetime of the pixel memory (nor the colortable if
// provided).
type Pixmap struct {
	info       ImageInfo
	colorTable *ColorTable
	pixelsData []byte
	rowBytes   int
}

func (p *Pixmap) Width() Scalar {
	return p.info.Width()
}

func (p *Pixmap) Height() Scalar {
	return p.info.Height()
}
