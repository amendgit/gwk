package ggk

// Region encapsulates the geometric region used to specify clippint areas for
// drawing.
type Region struct {
	bounds Rect
}

type RegionOp int

func (r *Region) Op(rgn *Region, op RegionOp) bool {
	return false
}
