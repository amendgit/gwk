package ggk

import (
	"image"
	"io"
)

func ImageFromReader(r io.Reader) *Bitmap {
	var img, _, err = image.Decode(r)
	if err != nil {
		return nil
	}

	var (
		pixels []uint8
		ct     ColorType
		at     AlphaType
	)
	switch i := img.(type) {
	case *image.RGBA:
		ct = KColorTypeRGBA8888
		at = KAlphaTypeUnpremul
		pixels = i.Pix
	}

	var info ImageInfo
	info.width, info.height = Scalar(img.Bounds().Dx()), Scalar(img.Bounds().Dy())
	info.SetColorType(ct)
	info.SetAlphaType(at)

	var bmp Bitmap
	bmp.InstallPixels(info, pixels, info.MinRowBytes(), nil)

	return &bmp
}
