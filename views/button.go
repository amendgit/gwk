// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	. "gwk/vango"
	"gwk/views/resc"
)

type Button struct {
	View
	canvas_enable *Canvas
}

func NewButton() *Button {
	var b = new(Button)
	b.SetID("button")
	b.canvas_enable = resc.LoadCanvas("data/button.png")
	b.SetBounds(b.canvas_enable.Bounds())
	b.Canvas().DrawCanvas(0, 0, b.canvas_enable, nil)
	return b
}
