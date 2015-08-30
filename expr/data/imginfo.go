package main

import (
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
	var filename = os.Args[1]
	fd, err := os.Open(filename)
	if err != nil {
		log.Printf("os.Open(%v) %v", filename, err)
		return
	}

	img, err := png.Decode(fd)
	if err != nil {
		log.Printf("image.Decode %v", err)
	}

	switch img.(type) {
	case *image.RGBA:
		log.Printf("RGBA")
	case *image.NRGBA:
		log.Printf("NRGBA")
	case *image.YCbCr:
		log.Printf("YCbCr")
	default:
		log.Printf("unknown")
	}

	return
}
