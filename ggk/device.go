package ggk

type Device interface {
	// ImageInfo returns ImageInfo for this device, If the canvas is not backed
	// by pixels (cpu or gpu), then the info's ColorType will be
	// KColorTypeUnknown.
	ImageInfo() ImageInfo

	// AccessGPURenderTarget return the device's gpu render target or nil.
	// AccessGPURenderTarget() *GPURenderTarget

	// OnAttachToCanvas is invoked whenever a device is installed in a canvas
	// (i.e., SetDevice, SaveLayer (for the new device created by the save),
	// and Canvas' BaseDevice & Bitmap - taking ctors). It allows the device
	// to prepare for drawing (e.g., locking their pixels, etc.)
	OnAttachToCanvas(*Canvas)

	OnDetachFromCanvas()
	SetMatrixClip(*Matrix, *Region, *ClipStack)
	DrawPaint(*Draw, *Paint)
	// DrawPoints(draw *Draw, mode PointMode, count int, pts []Point, paint *Paint)
	DrawRect(draw *Draw, rect Rect, paint *Paint)
	DrawOval(draw *Draw, oval Rect, paint *Paint)
	// DrawRRect(draw *Draw, RRect, *Paint)
	// DrawDRRect(*Draw, outer, inner RRect, *Paint)
	DrawPath(draw *Draw, path *Path, mat *Matrix, paint *Paint)
	DrawSprite(draw *Draw, bmp *Bitmap, x, y int, paint *Paint)
	DrawBitmapRect(draw *Draw, bmp *Bitmap, srcOrNil *Rect, dst Rect, paint *Paint) (finalDst Rect)
	DrawBitmapNine(draw *Draw, bmp *Bitmap, center Rect, dst Rect, paint *Paint)
	DrawImage(draw *Draw, image *Image, x, y Scalar, paint *Paint)
	// DrawImageRect(draw *Draw, image *Image, src Rect, dst Rect, paint *Paint, SrcRectConstraint)
	DrawText(draw *Draw, text string, x, y Scalar, paint *Paint)
	DrawPosText(draw *Draw, text string, pos []Scalar, paint *Paint)
	// DrawVertices(Draw, VertexMode, vertexCount int, verts []Point, texs []Point, colors []Color, xmode *Xfermode, indices []uint16, indexCount int, Paint)
	// DrawTextBlob(Draw, TextBlob, x, y Scalar, Paint, DrawFilter)
	// DrawPatch(Draw, cubics [12]Point, colors []Color, texCoords [4]Point, xmode Xfermode, Paint)
	// DrawAtlas(Draw, atlas Image, []RSXform, []Rect, []Color, count int, XfermodeMode, Paint)
	DrawDevice(draw *Draw, dev Device, x, y int, paint *Paint)
	DrawTextOnPath(draw *Draw, texts []string, len int, path *Path, mat *Matrix, paint *Paint)

	OnAccessBitmap() *Bitmap
	CanHandleImageFilter(*ImageFilter) bool
	FilterImage(filter *ImageFilter, bmp *Bitmap, ctxt *ImageFilterContext) (resultBmp *Bitmap, offset Point, ok bool)
	OnPeekPixels(pixmap *Pixmap) bool
	OnReadPixels(imageInfo ImageInfo, pixelBytes []byte, x, y int)
	OnWritePixels(imageInfo ImageInfo, pixelBytes []byte, x, y int)
	OnAccessPixels(pixmap *Pixmap) bool
	// OnCreateDevice(CreateInfo, Paint) Device
	Flush()
	// GetImageFilterCache() *ImageFilterCache
}

type BaseDevice struct {
	Device Device
}
