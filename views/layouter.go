// Copyright 2014 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package views

import (
	"log"
)

// ============================================================================

type Layouter interface {
	Layout(view View)
}

// ============================================================================

var (
	g_vertical_layouter   Layouter
	g_horizontal_layouter Layouter
)

func init_layout() {
	g_vertical_layouter = NewFuncLayouter(VerticalLayoutFunc)
}

// ============================================================================

type LayoutFunc func(view View)

type FuncLayouter struct {
	layout_func LayoutFunc
}

func NewFuncLayouter(layout_func LayoutFunc) *FuncLayouter {
	layouter := new(FuncLayouter)
	layouter.layout_func = layout_func
	return layouter
}

func (f *FuncLayouter) Layout(view View) {
	f.layout_func(view)
}

// ============================================================================

func VerticalLayoutFunc(view View) {
	bounds := view.Bounds()
	children := view.Children()

	if children == nil {
		return
	}

	log.Printf("vertical %v", bounds)

	height := bounds.Dy() / len(children)
	margin_top := 0
	for _, child := range children {
		w := bounds.Dx()
		h := height
		x := 0
		y := margin_top
		child.SetXYWH(x, y, w, h)
		log.Printf("x y w h %v %v %v %v", x, y, w, h)
		margin_top = margin_top + height
	}
}
