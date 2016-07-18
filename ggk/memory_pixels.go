package ggk

// Implement interface PixelsFactory.
type MemoryPixelsFactory struct {
}

func MemoryPixelsDefaultFactory() *MemoryPixelsFactory {
	// TOIMPL
	return nil
}

func (f *MemoryPixelsFactory) Create(info *ImageInfo, rowBytes int, ct *ColorTable) *Pixels {
	return NewMemoryPixelsAlloc(info, rowBytes)
}

type MemoryPixels struct {
	*Pixels
	storage []byte
}

// NewMemoryPixels return a new instance with the provided stroage, rowBytes,
// and optional colortable. The caller is responsible for managing the
// lifetime of the pixel storage buffer, as this pixels will not try
// to delete it.
//
// The pixels will ref the colortable (if not nil)
//
// Returns nil on failture.
func NewMemoryPixelsDirect(storage []byte) *MemoryPixels {
	var p MemoryPixels
	p.Pixels = NewPixels()
	p.storage = storage
	p.SetPrelocked(storage)
	return &p
}

func NewMemoryPixelsAlloc(info *ImageInfo, rowBytes int) *MemoryPixels {
	var p MemoryPixels
	p.Pixels = NewPixels()
	p.storage = make([]byte, info.SafeSize(rowBytes))
	p.SetPrelocked(p.storage)
	return &p
}

func isImageInfoValid(info *ImageInfo) bool {
	if info != nil && info.Width() < 0 || info.Height() < 0 ||
		!info.ColorType().IsVaild() || !info.AlphaType().IsValid() {
		return false
	}
	return true
}
