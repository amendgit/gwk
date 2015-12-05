package ggk

import (
	"bufio"
	"image"
	"image/png"
	"io"
	"os"
)

func BitmapFromReader(r io.Reader) (*Bitmap, error) {
	var img, _, err = image.Decode(r)
	if err != nil {
		return nil, err
	}

	var (
		pixels   []uint8
		ct       ColorType
		at       AlphaType
		rowBytes int
	)

	switch t := img.(type) {
	case *image.RGBA:
		ct = KColorTypeRGBA8888
		at = KAlphaTypeUnpremul
		pixels = t.Pix
		rowBytes = t.Stride
	}

	var info ImageInfo
	info.width, info.height = Scalar(img.Bounds().Dx()), Scalar(img.Bounds().Dy())
	info.SetColorType(ct)
	info.SetAlphaType(at)

	var bmp Bitmap
	bmp.InstallPixels(info, pixels, rowBytes, nil)

	return &bmp, nil
}

func BitmapFromFile(name string) (*Bitmap, error) {
	var f, err = os.Open(name)
	if err != nil {
		return nil, err
	}

	return BitmapFromReader(bufio.NewReader(f))
}

func BitmapToGoImage(bmp *Bitmap) (image.Image, error) {
	var img image.Image

	switch bmp.ColorType() {
	case KColorTypeRGBA8888:
		img = &image.RGBA{
			Pix:    bmp.PixelsData(),
			Stride: bmp.RowBytes(),
			Rect:   bmp.Bounds().ToGoRect(),
		}
	}

	return img, nil
}

func BitmapToPNGWriter(bmp *Bitmap, w io.Writer) error {
	var img, err = BitmapToGoImage(bmp)
	if err != nil {
		return err
	}

	return png.Encode(w, img)
}

func BitmapToPNGFile(bmp *Bitmap, name string) error {
	var f, err = os.Open(name)
	if err != nil {
		return err
	}

	return BitmapToPNGWriter(bmp, bufio.NewWriter(f))
}
