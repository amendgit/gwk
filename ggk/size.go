package ggk

type Size struct {
	width  Scalar
	height Scalar
}

var SizeZero Size

func MakeSize(w, h Scalar) Size {
	return Size{w, h}
}

func (s *Size) SetWH(w, h Scalar) {
	s.width, s.height = w, h
}

// Returns true if either width or height are <= 0.
func (s *Size) IsEmpty() bool {
	return s.width <= 0 || s.height <= 0
}

func (s *Size) Width() Scalar {
	return s.width
}

func (s *Size) Height() Scalar {
	return s.height
}

func (s *Size) Equal(o Size) bool {
	return s.width == o.width && s.height == o.height
}
