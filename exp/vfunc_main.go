package main

import (
	"log"
)

type Viewer interface {
	OnDraw()
}

type View struct {
}

func (v *View) OnDraw() {
	log.Printf("View.OnDraw()")
}

type ImageView struct {
	View
}

func (v *ImageView) OnDraw() {
	log.Printf("ImageView.OnDraw()")
}

func main() {
	image_view := new(ImageView)
	v := Viewer(image_view)
	v.OnDraw()
}
