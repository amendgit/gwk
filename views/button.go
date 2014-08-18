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
	img_normal *Canvas
}

func NewButton() *Button {
	var b = new(Button)
	b.SetID("button")
	b.img_normal = resc.FindCanvasByID("button_normal")
	b.SetBounds(b.img_normal.Bounds())
	return b
}

func (b *Button) OnDraw(event *DrawEvent) {
	// event.Canvas.DrawCanvas(0, 0, b.img_normal)
	// event.Canvas.DrawColor(0, 0, 250)
	event.Canvas.AlphaBlend(0, 0, b.img_normal)
}
