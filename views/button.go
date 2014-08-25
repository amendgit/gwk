// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	. "gwk/vango"
	"gwk/views/resc"
)

type Button struct {
	BaseView
	image_normal *Canvas `view:"image_normal"`
}

func NewButton() *Button {
	var b = new(Button)
	b.SetID("button")
	b.image_normal = resc.FindCanvasByID("button_normal")
	b.SetBounds(b.image_normal.Bounds())
	return b
}

func (b *Button) OnDraw(event *DrawEvent) {
	// event.Canvas.DrawCanvas(0, 0, b.image_normal)
	// event.Canvas.DrawColor(0, 0, 250)
	event.Canvas.AlphaBlend(0, 0, b.image_normal)
}
