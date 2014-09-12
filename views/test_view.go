// Copyright 2014 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	"fmt"
	"log"
)

type TestView struct {
	BaseView
}

func NewTestView() View {
	return new(TestView)
}

func (v *TestView) OnDraw(event *DrawEvent) {
	/// draw text and bounds
	ctxt := GraphicContext()
	ctxt.DrawColor(0xff, 0xff, 0xff)
	text := fmt.Sprintf("id: %v xywh: %v %v %v %v", v.ID(), v.X(), v.Y(), v.W(), v.H())
	ctxt.DrawText(text, v.LocalBounds())
	ctxt.SetStrokeColor(0x00, 0x00, 0xff)
	ctxt.StrokeRect(v.LocalBounds())
}

func (v *TestView) OnMouseEnter(event *MouseEvent) {
	log.Printf("OnMouseEnter: %v", v.ID())
}

func (v *TestView) OnMouseLeave(event *MouseEvent) {
	log.Printf("OnMouseLeave %v", v.ID())
}
