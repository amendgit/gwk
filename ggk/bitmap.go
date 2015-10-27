package ggk

import "sync/atomic"

type Bitmap struct {
	rowBytes       int
	flags          uint8
	info           ImageInfo
	colorTable     *ColorTable // cached from pixelRef
	pixels         []byte      // cached from pixelRef
	pixelRef       *PixelRef
	pixelRefOrigin Point
	pixelLockCount int32
}

// Swap the fields of the two bitmaps. This routine is guaranteed to never fail or throw.
func (bmp *Bitmap) Swap(other *Bitmap) {
	// bmp.colorTable, other.colorTable = other.colorTable, bmp.colorTable
	bmp.pixels, other.pixels = other.pixels, bmp.pixels
	bmp.info, other.info = other.info, bmp.info
	bmp.flags, other.flags = other.flags, bmp.flags
	bmp.rowBytes, other.rowBytes = other.rowBytes, bmp.rowBytes
}

func (bmp *Bitmap) Info() ImageInfo {
	return bmp.info
}

func (bmp *Bitmap) Width() int {
	return bmp.info.Width()
}

func (bmp *Bitmap) Height() int {
	return bmp.info.Height()
}

func (bmp *Bitmap) ColorType() ColorType {
	return bmp.info.colorType
}

func (bmp *Bitmap) AlphaType() AlphaType {
	return bmp.info.alphaType
}

func (bmp *Bitmap) ProfileType() ColorProfileType {
	return bmp.info.profileType
}

// Return the number of bytes per pixel based on the colortype. If the colortype is
// ColorType_Unknown, then 0 is returend.
func (bmp *Bitmap) BytesPerPixel() int {
	return bmp.info.BytesPerPixel()
}

// Return the rowBytes expressed as a number of pixels (like width and height).
// If the colortype is ColorType_Unknown, then 0 is returend.
func (bmp *Bitmap) RowBytesAsPixels() int {
	return bmp.rowBytes >> uint(bmp.ShiftPerPixel())
}

// Return the shift amount per pixel (i.e. 0 for 1-byte per pixel, 1 for 2-bytes per pixel
// colortypes. 2 for 4-bytes per pixel colortypes). Returns 0 for ColorType_Unknown.
func (bmp *Bitmap) ShiftPerPixel() int {
	return bmp.BytesPerPixel() >> 1
}

// IsEmpty returns true iff the bitmap has empty dimensions.
// Hey!  Before you use this, see if you really want to know DrawNothing() intead.
func (bmp *Bitmap) IsEmpty() bool {
	return bmp.info.IsEmpty()
}

// Return true iff the bitmap has no pixelref. Note: this can return true even if the
// dimensions of the bitmap are > 0 (see IsEmpty()).
// Hey! Before you use this, see if you really want to know DrawNothing() intead.
func (bmp *Bitmap) IsNil() bool {
	return bmp.pixelRef == nil
}

// Return true iff drawing the bitmap has no effect.
func (bmp *Bitmap) DrawNothing() bool {
	return bmp.IsEmpty() || bmp.IsNil()
}

// Return the number of bytes between subsequent rows of the bitmap.
func (bmp *Bitmap) RowBytes() int {
	return bmp.rowBytes
}

// Set the bitmap's alphaType, returing true on success. If false is
// returned, then the specified new alphaType is incompatible with the
// colortype, and the current alphaType is unchanged.
//
// Note: this changes the alpahType for the underlying types, which means
// that all bitmaps that might be sharing (subsets of) the pixels will
// be affected.
func (bmp *Bitmap) SetAlphaType(alphaType AlphaType) bool {
	alphaType, err := bmp.info.colorType.ValidateAlphaType(alphaType)
	if err != nil {
		return false
	}

	if bmp.info.alphaType != alphaType {
		bmp.info.SetAlphaType(alphaType)
		if bmp.pixelRef != nil {
			bmp.pixelRef.SetAlphaType(alphaType)
		}
	}

	return true
}

func (bmp *Bitmap) updatePixelsFromRef() {
	if bmp.pixelRef == nil {
		return
	}

	if bmp.pixelLockCount > 0 {
		var p = bmp.pixelRef.Pixels()
		if p != nil {
			var idx = int(bmp.pixelRefOrigin.Y())*bmp.rowBytes +
				int(bmp.pixelRefOrigin.X())*bmp.info.BytesPerPixel()
			p = p[idx:]
		}
		bmp.pixels = p
		bmp.colorTable = bmp.pixelRef.ColorTable()
	} else {
		bmp.pixels = nil
		bmp.colorTable = nil
	}
}

// func (pslf *Bitmap) Pixels()                  {}
// func (pslf *Bitmap) Size() int                {}
// func (pslf *Bitmap) SafeSize() int            {}
// func (pslf *Bitmap) ComputeSize64() int64     {}
// func (pslf *Bitmap) ComputeSafeSize64() int64 {}
// func (pslf *Bitmap) Immutable() bool          {}
// func (pslf *Bitmap) SetImmutale()             {}
// func (pslf *Bitmap) Opaque()                  {}
// func (pslf *Bitmap) IsVolatile() bool         {}
// func (pslf *Bitmap) SetIsVolatile()           {}
func (bmp *Bitmap) Reset() {
	bmp.freePixels()
	var zero Bitmap
	*bmp = zero
}

// func (pslf *Bitmap) ComputeIsOpaque() bool    {}

// func (pslf *Bitmap) Bounds() Rect {}
// func (pslf *Bitmap) Dimensions() Size {}

func (bmp *Bitmap) SetInfo(imageInfo ImageInfo, rowBytes int) bool {
	alphaType, err := imageInfo.ColorType().ValidateAlphaType(imageInfo.AlphaType())
	if err != nil {
		bmp.Reset()
		return false
	}

	// alphaType is the real value.

	var minRowBytes int64 = imageInfo.MinRowBytes64()
	if int64(int32(minRowBytes)) != minRowBytes {
		bmp.Reset()
		return false
	}

	if imageInfo.Width() < 0 || imageInfo.Height() < 0 {
		bmp.Reset()
		return false
	}

	if imageInfo.ColorType() == KColorTypeUnknown {
		rowBytes = 0
	} else if rowBytes == 0 {
		rowBytes = int(minRowBytes)
	} else if !imageInfo.ValidRowBytes(rowBytes) {
		bmp.Reset()
		return false
	}

	bmp.freePixels()

	bmp.info.SetAlphaType(alphaType)
	bmp.rowBytes = rowBytes

	return true
}

func (bmp *Bitmap) AllocPixels(requestedInfo ImageInfo, rowBytes int) bool {
	if requestedInfo.ColorType() == KColorTypeIndex8 {
		bmp.Reset()
		return false
	}

	if !bmp.SetInfo(requestedInfo, rowBytes) {
		bmp.Reset()
		return false
	}

	// SetInfo may have corrected info (e.g. 565 is always opaque).
	var correctedInfo = bmp.Info()

	// SetInfo may have computed a valid rowBytes if 0 were passed in
	rowBytes = bmp.RowBytes()

	var pixelRef *PixelRef = MallocPixelRefDefaultFactory().Create(correctedInfo, rowBytes, nil)
	if pixelRef == nil {
		bmp.Reset()
		return false
	}

	bmp.SetPixelRef(pixelRef, PointZero)

	if !bmp.LockPixels() {
		bmp.Reset()
		return false
	}

	return true
}

// Assign a pixelRef and origin to the bitmap. PixelRefs are reference.
// so the existing one (if any) will be unref'd and the new one will be
// ref'd. (x,y) specify the offset within the pixelRef's pixels for the
// top/left corner of the bitmap. For a bitmap that encompass the entire
// pixels of the pixelref, these will be (0,0).
func (bmp *Bitmap) SetPixelRef(pixelRef *PixelRef, origin Point) *PixelRef {
	// TOIMPL
	return bmp.pixelRef
}

// Call this to ensure that the bitmap points to the current pixel address
// in the pixelRef. Balance it with a call to UnlockPixels(). These calls
// are harmless if there is no pixelRef.
func (bmp *Bitmap) LockPixels() bool {
	if bmp.pixelRef != nil && atomic.AddInt32(&bmp.pixelLockCount, 1) == 0 {
		bmp.pixelRef.LockPixels()
		bmp.updatePixelsFromRef()
	}
	return true
}

// When you are finished access the pixel memory, call this to balance a
// previous call to LockPixels(). This allows pixelRefs that implement
// cached/deferred image decoding to know when there are active clients of
// a given image.
func (bmp *Bitmap) UnlockPixels() {
	// TOIMPL
}

// func AllocN32Pixels(width, height int, isOpaque bool) {}
// func InstallPixels(image_info *ImageInfo, pixels []byte, row_bytes int, color_table *ColorTable, context interface{}) bool {
// }

// unreference any pixelRefs or colorTables.
func (bmp *Bitmap) freePixels() {
	if bmp.pixelRef != nil {
		if bmp.pixelLockCount > 0 {
			bmp.pixelRef.UnlockPixels()
		}
		bmp.pixelRef = nil
		bmp.pixelRefOrigin = PointZero
	}

	bmp.pixelLockCount = 0
	bmp.pixels = nil
	bmp.colorTable = nil
}

// func CanCopy(color_type ColorType) bool {}
// func DeepCopy() *Bitmap                 {}
// func Copy(color_type ColorType) *Bitmap {}

// func SetPixels(pixels []byte, colort_table ColorTable)                          {}
// func PeekPixels(pixmap Pixmap)                                                  {}
// func CopyPixels(dst_size int, dst_row_bytes int, preserved_dst_pad bool) []byte {}
// func ColorAtIndex(x, y int) Color                                               {}
