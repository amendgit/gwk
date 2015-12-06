package ggk

type Draw struct {
	dst        *Pixmap
	mat        *Matrix
	region     *Region
	rasterClip *RasterClip
	clipStack  *ClipStack
	dev        *BaseDevice
	// procs      *DrawProcs
}

func (d *Draw) DrawPaint(paint *Paint) {
	if d.rasterClip.IsEmpty() {
		return
	}

	// var devRect Rect = MakeRect(0, 0, d.dst.Width(), d.dst.Height())

	if d.rasterClip.IsBW() {
		// TOIMPL
	}

	// normal case: use a blitter
	// ScanFillRect(devRect, p.rasterClip, blitter)
}

func (d *Draw) DrawRect(rect Rect, paint *Paint) {
	toimpl()
	return
}
