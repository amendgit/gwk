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

func (m *Matrix) TypeMask() MatrixTypeMask {
	if (m.typeMask & KMatrixTypeMaskUnknown) != 0 {
		// m.typeMask = m.ComputeTypeMask()
	}
	// only return the public masks.
	return MatrixTypeMask(m.typeMask & 0xF)
}
