package ggk

// MatrixTypeMask is the enum of the bit fields for the mask return by Type()
// Use this to identify the complexity of the matrix.
type MatrixTypeMask int

const (
	KMatrixTypeMaskIdentity    = 0
	KMatrixTypeMaskTranslate   = 0x01 // set if the matrix has translation
	KMatrixTypeMaskScale       = 0x02 // set if the matrix has X or Y scale
	KMatrixTypeMaskAffine      = 0x04 // set if the matrix skews or rotates
	KMatrixTypeMaskPerspective = 0x08 // set if the matrix is in perspective
	KMatrixTypeMaskUnknown     = 0x80
)

// Matrix holds a 3x3 matrix for transforming coordinates. Matrix does not have
// a constructor, so it must be explicitly initialized using either reset() - to
// construct an identity matrix, or one of the set functions (e.g. SetTranslate,
// SetRotate, etc.).
type Matrix struct {
	mat      [9]Scalar
	typeMask uint32
}

const (
	KMScaleX = iota
	KMSkewX
	KMTransX
	KMSkewY
	KMScaleY
	KMTransY
	KMPersp0
	KMPersp1
	KMPersp2
)

// Affine arrays are in column major order
// because that's how PDF and XPS like it.
const (
	KAScaleX = iota
	KASkewY
	KASkewX
	KAScaleY
	KATransX
	KATransY
);

func (m *Matrix) SetTypeMask(mask int) {
	m.typeMask = mask
}

func (m *Matrix) TypeMask() MatrixTypeMask {
	if (m.typeMask & KMatrixTypeMaskUnknown) != 0 {
		// m.typeMask = m.ComputeTypeMask()
	}
	// only return the public masks.
	return MatrixTypeMask(m.typeMask & 0xF)
}

// [scale-x    skew-x      trans-x]   [X]   [X']
// [skew-y     scale-y     trans-y] * [Y] = [Y']
// [persp-0    persp-1     persp-2]   [1]   [1 ]
func (m *Matrix) Reset() {
	m.mat[KMScaleX], m.mat[KMSkewX ], m.mat[KMTransX] = 1, 0, 0
	m.mat[KMSkewY ], m.mat[KMScaleY], m.mat[KMTransY] = 0, 1, 0
	m.mat[KMPersp0], m.mat[KMPersp1], m.mat[KMPersp2] = 0, 0, 1
	// setTypeMask(kIdentity_Mask | kRectStaysRect_Mask)
}