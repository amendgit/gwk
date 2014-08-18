// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	. "gwk/vango"
	. "image"
)

type DrawEvent struct {
	Owner     View
	Canvas    *Canvas
	DirtyRect Rectangle
}

type MouseEvent struct {
	Owner    View
	Location Point
}

func NewMouseEvent(pt Point) *MouseEvent {
	mouse_event := new(MouseEvent)
	mouse_event.Location = pt
	return mouse_event
}
