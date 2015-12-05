package ggk

import (
	"image"
)

var RectZero Rect

type Rect struct {
	left   Scalar
	top    Scalar
	width  Scalar
	height Scalar
}

func MakeRect(x, y, width, height Scalar) Rect {
	return Rect{x, y, width, height}
}

// Make rectangle from width and size, the left and top set to 0.
func MakeRectWH(width, height Scalar) Rect {
	return Rect{0, 0, width, height}
}

// Make rectangle from (left, top, right, bottom).
func MakeRectLTRB(left, top, right, bottom Scalar) Rect {
	return Rect{left, top, right - left, bottom - top}
}

// Return te left edge of the rect.
func (r Rect) Left() Scalar {
	return r.left
}

func (r Rect) X() Scalar {
	return r.left
}

// Return the top edge of the rect.
func (r Rect) Top() Scalar {
	return r.top
}

func (r Rect) Y() Scalar {
	return r.top
}

// Return the rectangle's width. This does not check for a valid rect
// (i.e. left <= right) so the result may be negative.
func (r Rect) Width() Scalar {
	return r.width
}

// Returns the rectangle's height. This does not check for a vaild rect
// (i.e. top <= bottom) so the result may be negative.
func (r Rect) Height() Scalar {
	return r.height
}

// Returns the rectangle's right edge.
func (r Rect) Right() Scalar {
	return r.left + r.width
}

// Returns the rectangle's bottom edge.
func (r Rect) Bottom() Scalar {
	return r.top + r.height
}

// Returns the rectangle's center x.
func (r Rect) CenterX() Scalar {
	return r.left + r.width*0.5
}

// Returns the rectangle's center Y.
func (r Rect) CenterY() Scalar {
	return r.top + r.height*0.5
}

// Return true if the rectangle's width or height are <= 0
func (r Rect) IsEmpty() bool {
	return r.left <= 0 || r.height <= 0
}

// Return true if the two rectangles have same position and size.
func (a Rect) Equal(b Rect) bool {
	return a.left == b.left && a.top == b.top && a.width == b.width &&
		a.height == b.height
}

// Set the rectangle's edges with (x, y, w, h)
func (r *Rect) SetXYWH(x, y, width, height Scalar) {
	r.left, r.top, r.width, r.height = x, y, width, height
}

// Set the rectangle's edges with (left, top, right, bottom)
func (r *Rect) SetLTRB(left, top, right, bottom Scalar) {
	r.left, r.top, r.width, r.height = left, top, right-left, bottom-top
}

func (r *Rect) IntersectXYWH(x, y, w, h Scalar) bool {
	toimpl()
	return false
}

func (r Rect) ToGoRect() image.Rectangle {
	return image.Rect(int(r.left), int(r.top), int(r.width), int(r.height))
}
