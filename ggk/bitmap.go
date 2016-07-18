package ggk

import (
	"errors"
	"sync/atomic"
)

type Bitmap struct {
	rowBytes int
	flags    uint8

	info       *ImageInfo
	colorTable *ColorTable

	pixels         *Pixels
	pixelOrigin    Point
	pixelLockCount int32
}

// Swap the fields of the two bitmaps. This routine is guaranteed to never fail or throw.
func (bmp *Bitmap) Swap(otr *Bitmap) {
	*bmp, *otr = *otr, *bmp
}

func (bmp *Bitmap) Info() *ImageInfo {
	return bmp.info
}

func (bmp *Bitmap) Width() Scalar {
	return bmp.info.Width()
}

func (bmp *Bitmap) Height() Scalar {
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
// KColorTypeUnknown, then 0 is returend.
func (bmp *Bitmap) BytesPerPixel() int {
	return bmp.info.BytesPerPixel()
}

// Return the rowBytes expressed as a number of pixels (like width and height).
// If the colortype is KColorTypeUnknown, then 0 is returend.
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
	return bmp.pixels == nil
}

var ErrBitmapIsNotValid = errors.New(`error: bitmap is not valid.`)

// IsValid return true iff the bitmap has valid imageInfo, pixels and colorTable
func (bmp *Bitmap) IsValid() bool {
	if !bmp.info.IsValid() {
		return false
	}
	if !bmp.info.ValidRowBytes(bmp.rowBytes) {
		return false
	}
	if bmp.info.ColorType() == KColorTypeRGB565 &&
		bmp.info.AlphaType() != KAlphaTypeOpaque {
		return false
	}
	// TOIMPL
	if bmp.pixels != nil {
		if bmp.pixelLockCount <= 0 &&
			//    !bmp.pixels.IsLock() &&
			bmp.rowBytes < bmp.info.MinRowBytes() &&
			bmp.pixelOrigin.X() < 0 &&
			bmp.pixelOrigin.Y() < 0 &&
			bmp.info.Width() < bmp.Width()+bmp.pixelOrigin.X() &&
			bmp.info.Height() < bmp.Height()+bmp.pixelOrigin.Y() {
			return false
		}
	} else {
		if bmp.colorTable != nil {
			return false
		}
	}
	return true
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
	}
	return true
}

func (bmp *Bitmap) Pixels() *Pixels {
	return bmp.pixels
}

func (bmp *Bitmap) PixelsBytes() []byte {
	bmp.pixels.LockPixels()
	var bytes = bmp.pixels.Bytes()
	return bytes
}

func (bmp *Bitmap) InstallPixels(requestedInfo ImageInfo, pixelsBytes []byte, rowBytes int, ct *ColorTable) bool {
	if !bmp.SetInfo(requestedInfo, rowBytes) {
		// release pixels
		bmp.Reset()
		return false
	}
	if pixelsBytes == nil {
		// release pixels
		return true // we behaved as if they called setInfo()
	}
	var pixels = NewMemoryPixelsDirect(pixelsBytes)
	if pixels == nil {
		bmp.Reset()
		return false
	}
	bmp.pixels = pixels.Pixels
	// since we're already allocated, we LockPixels right away.
	bmp.LockPixels()
	if !bmp.IsValid() {
		// 	log.Printf(`xyz`)
	}
	return true
}

func (bmp *Bitmap) Reset() {
	bmp.freePixels()
	var zero Bitmap
	*bmp = zero
}

func (bmp *Bitmap) Bounds() Rect {
	var (
		x      = bmp.pixelOrigin.X()
		y      = bmp.pixelOrigin.Y()
		width  = bmp.info.Width()
		height = bmp.info.Height()
	)
	return MakeRect(x, y, width, height)
}

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
	bmp.info = imageInfo.MakeAlphaType(alphaType)
	bmp.rowBytes = rowBytes
	return true
}

var ErrAllocPixels = errors.New(`ERROR: bad imageInfo, rowBytes. or allocate failed`)

func (bmp *Bitmap) AllocPixels(requestedInfo ImageInfo, rowBytes int) error {
	if requestedInfo.ColorType() == KColorTypeIndex8 {
		bmp.Reset()
		return ErrAllocPixels
	}
	if !bmp.SetInfo(requestedInfo, rowBytes) {
		bmp.Reset()
		return ErrAllocPixels
	}
	// SetInfo may have corrected info (e.g. 565 is always opaque).
	var correctedInfo = bmp.Info()
	// SetInfo may have computed a valid rowBytes if 0 were passed in
	rowBytes = bmp.RowBytes()
	// Allocate memories.
	var pixels = NewMemoryPixelsAlloc(correctedInfo, rowBytes)
	if pixels == nil {
		bmp.Reset()
		return ErrAllocPixels
	}
	bmp.pixels = pixels.Pixels
	if bmp.LockPixels() != nil {
		bmp.Reset()
		return ErrAllocPixels
	}
	return ErrAllocPixels
}

// Assign a pixels and origin to the bitmap. Pixels are reference.
// so the existing one (if any) will be unref'd and the new one will be
// ref'd. (x,y) specify the offset within the pixelRef's pixels for the
// top/left corner of the bitmap. For a bitmap that encompass the entire
// pixels of the pixel ref, these will be (0,0).
func (bmp *Bitmap) SetPixels(pixels *Pixels, origin Point) {
	toimpl()
}

// Call this to ensure that the bitmap points to the current pixel address
// in the pixels. Balance it with a call to UnlockPixels(). These calls
// are harmless if there is no pixelRef.
func (bmp *Bitmap) LockPixels() error {
	if bmp.pixels != nil && atomic.AddInt32(&bmp.pixelLockCount, 1) == 1 {
		bmp.pixels.LockPixels()
	}
	return nil
}

// When you are finished access the pixel memory, call this to balance a
// previous call to LockPixels(). This allows pixelRefs that implement
// cached/deferred image decoding to know when there are active clients of
// a given image.
func (bmp *Bitmap) UnlockPixels() error {
	if bmp.pixels != nil && atomic.AddInt32(&bmp.pixelLockCount, -1) == 0 {
		bmp.pixels.UnlockPixels()
	}
	return nil
}

// Unreference any pixels or colorTables.
func (bmp *Bitmap) freePixels() {
	if bmp.pixels != nil {
		if bmp.pixelLockCount > 0 {
			bmp.UnlockPixels()
		}
		bmp.pixels = nil
		bmp.pixelOrigin = PointZero
	}
	bmp.pixelLockCount = 0
	bmp.pixels = nil
	bmp.colorTable = nil
}
