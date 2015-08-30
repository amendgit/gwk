package ggk

// type Bitmap struct {
//	row_bytes uint32
//	flags     uint32
//	pixels    []byte
//	info      ImageInfo
//  color_table ColorTable
// }

// func (pslf *Bitmap) Swap(other *Bitmap) {}

// func (slf *Bitmap) Info() ImageInfo {}
// func (pslf *Bitmap) Width() int {
//	return pslf.info.Width()
// }

// func (pslf *Bitmap) Height() int {
//	return slfp.info.Height()
// }

// func (slf *Bitmap) ColorType() ColorType {}
// func (slf *Bitmap) AlphaType() AlphaType {}
// func (slf *Bitmap) ProfileType() ProfileType {}
// func (pslf *Bitmap) BytesPerPixel()     {}
// func (pslf *Bitmap) RowBytesAsPixels()  {}
// func (pslf *Bitmap) ShiftPerPixel() int {}
// func (pslf *Bitmap) Empty() bool        {}
// func (pslf *Bitmap) IsNil() bool        {}
// func (pslf *Bitmap) DrawNothing() bool  {}
// func (pslf *Bitmap) RowBytes() int      {}

// func (slf *Bitmap) SetAlphaType(alpha_type AlphaType) bool {}
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
// func (pslf *Bitmap) Reset()                   {}
// func (pslf *Bitmap) ComputeIsOpaque() bool    {}

// func (pslf *Bitmap) Bounds() Rect {}
// func (pslf *Bitmap) Dimensions() Size {}
// func (pslf *Bitmap) SetInfo(image_info *ImageInfo) {}

// func (pslf *Bitmap) AllocPixels(image_info *ImageInfo, row_bytes int) {}
// func (pslf *Bitmap) AllocPixels(image_info *ImageInfo) {}
// func AllocN32Pixels(width, height int, isOpaque bool) {}
// func InstallPixels(image_info *ImageInfo, pixels []byte, row_bytes int, color_table *ColorTable, context interface{}) bool {
// }

// func CanCopy(color_type ColorType) bool {}
// func DeepCopy() *Bitmap                 {}
// func Copy(color_type ColorType) *Bitmap {}

// func SetPixels(pixels []byte, colort_table ColorTable)                          {}
// func PeekPixels(pixmap Pixmap)                                                  {}
// func CopyPixels(dst_size int, dst_row_bytes int, preserved_dst_pad bool) []byte {}
// func ColorAtIndex(x, y int) Color                                               {}
