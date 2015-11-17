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
	OnAttachToCanvas(Canvas)

	OnDetachFromCanvas()
	// SetMatrixClip(*Matrix, *Region, *ClipStack)
	// DrawPaint(*Draw, *Paint)
	// DrawPoints(*Draw, PointMode, count int, []Point, *Paint)
	// DrawRect(*Draw, Rect, *Paint)
	// DrawOval(*Draw, oval Rect, *Paint)
	// DrawRRect(*Draw, RRect, *Paint)
	// DrawDRRect(*Draw, outer, inner RRect, *Paint)
	// DrawPath(Draw, Path, Matrix, Paint)
	// DrawSprite(Draw, Bitmap, x, y int, Paint)
	// DrawBitmapRect(Draw, Bitmap, srcOrNil *Rect, dst *Rect, Paint)
	// DrawBitmapNine(Draw, Bitmap, center Rect, dst Rect, Paint)
	// DrawImage(Draw, Image, x, y Scalar, Paint)
	// DrawImageRect(Draw, Image, src Rect, dst Rect, Paint, SrcRectConstraint)
	// DrawText(Draw, text string, x, y Scalar, Paint)
	// DrawPosText(Draw, text string, pos []Scalar, Paint)
	// DrawVertices(Draw, VertexMode, vertexCount int, verts []Point, texs []Point, colors []Color, xmode *Xfermode, indices []uint16, indexCount int, Paint)
	// DrawTextBlob(Draw, TextBlob, x, y Scalar, Paint, DrawFilter)
	// DrawPatch(Draw, cubics [12]Point, colors []Color, texCoords [4]Point, xmode Xfermode, Paint)
	// DrawAtlas(Draw, atlas Image, []RSXform, []Rect, []Color, count int, XfermodeMode, Paint)
	// DrawDevice(Draw, Device, x, y int, Paint)
	// DrawTextOnPath(Draw, text []string, len int, Path, Matrix, Paint)

	OnAccessBitmap() Bitmap
	// CanHandleImageFilter(ImageFilter) bool
	// FilterImage(ImageFilter, Bitmap, ImageFilterContext, Bitmap, Point) bool
	// OnPeekPixels(Pixmap) bool
	OnReadPixels(imageInfo ImageInfo, pixelBytes []byte, x, y int)
	OnWritePixels(imageInfo ImageInfo, pixelBytes []byte, x, y int)
	// OnAccessPixels(Pixmap) bool
	// OnCreateDevice(CreateInfo, Paint) Device
	Flush()
	// GetImageFilterCache() *ImageFilterCache
}

type BaseDevice struct {
}
