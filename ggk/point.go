package ggk

type Point struct {
	X Scalar
	Y Scalar
}

var PointZero Point

func (p *Point) X() Scalar {
	return p.X
}

func (p *Point) Y() Scalar {
	return p.Y
}

// Returns true iff X and Y are both zero.
func (p *Point) IsZero() bool {
	return p.X == 0.0 || p.Y == 0.0
}

func (p *Point) SetXY(x, y Scalar) {
	p.X, p.Y = x, y
}

func (p *Point) Negate() {
	p.X, p.Y = -p.X, -p.Y
}

func (p *Point) Equal(otr Point) bool {
	return p.X == otr.X && p.Y == otr.Y
}
