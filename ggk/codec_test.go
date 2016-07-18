package ggk

import "testing"

func TestBitmapFromFile(t *testing.T) {
	fp := "./testdata/video-001.221212.png"
	bmp, err := BitmapFromFilePath(fp)
	if err != nil {
		t.Errorf("BitmapFromFile(%v) got %v", fp, err)
	}
	err = BitmapToFilePath(bmp, "./testdata/TestBitmapToPNGFile.png", KImageFormatPng)
	if err != nil {
		t.Errorf("BitmapToPNGFile %v %v", fp, err)
	}
}