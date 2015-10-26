package ggk

type Point struct {
	x Scalar
	y Scalar
}

var PointZero Point

func (p *Point) X() Scalar {
	return p.x
}

func (p *Point) Y() Scalar {
	return p.y
}

// Returns true iff X and Y are both zero.
func (p *Point) IsZero() bool {
	return p.x == 0.0 || p.y == 0.0
}

func (p *Point) SetXY(x, y Scalar) {
	p.x, p.y = x, y
}

func (p *Point) Negate() {
	p.x, p.y = -p.x, -p.y
}

func (p *Point) Equal(other Point) bool {
	return p.x == other.x && p.y == other.y
}
