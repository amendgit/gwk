package ggk

import (
	"bufio"
	"image"
	"image/png"
	"io"
	"os"
)

type ImageFormatType int
const (
	KImageFormatTypeUnknown = iota
	KImageFormatTypeJpeg
	KImageFormatTypePng
	KImageFormatTypeWebp
)

func BitmapToGoImage(bmp *Bitmap) (image.Image, error) {
	var img image.Image
	switch bmp.ColorType() {
	case KColorTypeRGBA8888:
		img = &image.RGBA{
			Pix:    bmp.PixelsBytes(),
			Stride: bmp.RowBytes(),
			Rect:   bmp.Bounds().ToGoRect(),
		}
	default:
		return nil, errorf(`color type %v is not support.`, bmp.ColorType())
	}
	return img, nil
}

func BitmapFromReader(r io.Reader) (*Bitmap, error) {
	var img, _, err = image.Decode(r)
	if err != nil {
		return nil, err
	}
	var (
		pixels   []uint8
		ct       ColorType = KColorTypeRGBA8888
		at       AlphaType
		rowBytes int
	)
	switch t := img.(type) {
	case *image.RGBA:
		ct = KColorTypeRGBA8888
		at = KAlphaTypeUnpremul
		pixels = t.Pix
		rowBytes = t.Stride
	default:
		return nil, errorf(`color type %v is not support.`, t)
	}
	var info ImageInfo
	info.width, info.height = Scalar(img.Bounds().Dx()), Scalar(img.Bounds().Dy())
	info.SetColorType(ct)
	info.SetAlphaType(at)
	var bmp Bitmap
	bmp.InstallPixels(info, pixels, rowBytes, nil)
	return &bmp, nil
}

func BitmapToWriter(bmp *Bitmap, w io.Writer, format ImageFormatType) error {
	var img, err = BitmapToGoImage(bmp)
	if err != nil {
		return err
	}
	switch format {
	case KImageFormatTypePng:
		return png.Encode(w, img)
	default:
		return errorf(`image format %v is not support`, format)
	}
	return nil
}

func BitmapFromFilePath(fp string) (*Bitmap, error) {
	var fd, err = os.Open(fp)
	if err != nil {
		return nil, err
	}
	return BitmapFromReader(bufio.NewReader(fd))
}

func BitmapToFilePath(bmp *Bitmap, fp string, format ImageFormatType) error {
	var fd, err = os.Create(fp)
	if err != nil {
		return err
	}
	return BitmapToWriter(bmp, bufio.NewWriter(fd), format)
}
