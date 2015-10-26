package main

import (
	"fmt"
)

type View interface {
	ViewName() string
}

type BaseView struct {
	View
}

func (b *BaseView) PrintViewInfo() {
	fmt.Printf("BaseView：%v\n", b.ViewName())
}

type ImageView struct {
	BaseView
}

func NewImageView() *ImageView {
	var i = new(ImageView)
	i.View = i
	return i
}

func (i *ImageView) ViewName() string {
	return "ImageView"
}

func main() {
	var obj = NewImageView()
	obj.PrintViewInfo()
}
