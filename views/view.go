// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	. "gwk/vango"
	. "image"
)

type View interface {
	ID() string
	SetID(id string)

	AddChild(child View)
	Children() []View

	Parent() View
	SetParent(parent View)

	Canvas() *Canvas
	SetCanvas(canvas *Canvas)

	ToAbsPoint(pt Point) Point
	ToDevicePoint(pt Point) Point
	ToAbsRect(rc Rectangle) Rectangle
	ToDeviceRect(rc Rectangle) Rectangle

	OnDraw(event *DrawEvent)

	OnMouseEnter(event *MouseEvent)
	OnMouseLeave(event *MouseEvent)

	ScheduleDraw()
	ScheduleDrawRect(dirty Rectangle)

	X() int
	Y() int
	W() int
	H() int
	XYWH() (x, y, w, h int)
	SetXYWH(x, y, w, h int)

	Left() int
	Top() int
	Width() int
	Height() int
	SetLeft(left int)
	SetTop(top int)
	SetWidth(width int)
	SetHeight(height int)

	Bounds() Rectangle
	LocalBounds() Rectangle
	SetBounds(bounds Rectangle)

	UIMap() UIMap
	SetUIMap(ui_map UIMap)
	MockUp(ui UIMap)

	Layouter() Layouter
	SetLayouter(l Layouter)

	SetDelegate(delegate ViewDelegate)
	Delegate() ViewDelegate
}
