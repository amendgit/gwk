package ggk

type Canvas struct {
	// surfaceProps SurfaceProps
	saveCount int
	// metaData *MetaData
	// baseSurface *BaseSurface
	// mcStack *Dequeue
	// clipStack ClipStack
	// mcRec *MCRec // points to top of the stack

	deviceCMDirty bool

	cachedLocalClipBounds      Rect
	cachedLocalClipBoundsDirty bool

	allowSoftClip          bool
	allowSimplifyClip      bool
	conservativeRasterClip bool
}

func (c *Canvas) readPixels(dstInfo ImageInfo, dstData []byte, rowBytes int,
	x, y int) error {
	// TOIMPL
	return nil
}
