// Copyright 2014 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package layout

import (
	. "gwk/views"
	. "image"
)

type LayoutFunc func(Rectangle, []Viewer)

type FuncLayouter struct {
	layout_func LayoutFunc
}

func (f *FuncLayouter) Layout(bounds Rectangle, views []Viewer) {
	f.layout_func(bounds, views)
}

func NewFuncLayouter(layout_func LayoutFunc) *FuncLayouter {
	f := new(FuncLayouter)
	f.layout_func = layout_func
	return f
}
