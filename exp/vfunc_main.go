package main

import (
	"fmt"
)

type View interface {
	OnDraw()
}

type BaseView struct {
}

func (v *BaseView) OnDraw() {
	fmt.Printf("BaseView.OnDraw()\n")
}

type ImageView struct {
	BaseView
}

func (v *ImageView) OnDraw() {
	fmt.Printf("ImageView.OnDraw()\n")
}

func main() {
	image_view := new(ImageView)
	view := View(image_view)
	view.OnDraw()

	view = &(image_view.BaseView)
	view.OnDraw()
}
