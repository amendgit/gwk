package ggk

type RasterClip struct {
	bw                     *Region
	isBW                   bool
	forceConservativeRects bool
	aaclip                 *AAClip

	isEmpty bool
	isRect  bool
}

func NewRasterClip(conservativeRasterClip bool) *RasterClip {
	toimpl()
	return nil
}

func (r *RasterClip) IsEmpty() bool {
	return r.isEmpty
}

func (r *RasterClip) IsBW() bool {
	return r.isBW
}

func (r *RasterClip) SetRect(rect Rect) {
	toimpl()
}
