// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	. "image"
	// . "gwk/vango"
	// . "gwk/views/resc"
)

type MainFrame struct {
	BaseView

	left_panel  View
	right_panel View
	main_panel  View
	toolbar     View
}

func NewMainFrame() *MainFrame {
	var v = new(MainFrame)
	v.SetID("main_frame")
	v.SetLayouter(v)
	v.SetXYWH(0, 0, 50, 50)
	// v.canvas_bkg = resc.LoadCanvas("data/texture.png")
	// v.Canvas().DrawTexture(v.Canvas().Bounds(), v.canvas_bkg, v.canvas_bkg.Bounds())
	return v
}

func (v *MainFrame) MockUp(ui UIMap) {
	if left_panel, ok := ui.UIMap("left_panel"); ok {
		v.left_panel = MockUp(left_panel)
		v.AddChild(v.left_panel)
	}

	if right_panel, ok := ui.UIMap("right_panel"); ok {
		v.right_panel = MockUp(right_panel)
		v.AddChild(v.right_panel)
	}

	if main_panel, ok := ui.UIMap("main_panel"); ok {
		v.main_panel = MockUp(main_panel)
		v.AddChild(v.main_panel)
	}

	if toolbar, ok := ui.UIMap("toolbar"); ok {
		v.toolbar = MockUp(toolbar)
		v.AddChild(v.toolbar)
	}
}

func (m *MainFrame) Layout(parent View) {
	// m.Canvas().DrawTexture(m.Canvas().Bounds(), m.canvas_bkg, m.canvas_bkg.Bounds())
	var r Rectangle

	if m.left_panel != nil {
		r = m.get_left_panel_bounds()
		m.left_panel.SetBounds(r)
	}

	if m.right_panel != nil {
		r = m.get_right_panel_bounds()
		m.right_panel.SetBounds(r)
	}

	if m.main_panel != nil {
		r = m.get_main_panel_bounds()
		m.main_panel.SetBounds(r)
	}

	if m.toolbar != nil {
		r = m.get_toolbar_bounds()
		m.toolbar.SetBounds(r)
	}
}

// ============================================================================

func (m *MainFrame) get_left_panel_bounds() Rectangle {
	bounds := m.LocalBounds()
	const panel_width = 200

	y := m.get_toolbar_height()
	w := panel_width
	h := bounds.Dy()

	return Rect(0, y, w, h)
}

func (m *MainFrame) get_toolbar_height() int {
	return 30
}

func (m *MainFrame) get_toolbar_bounds() Rectangle {
	bounds := m.LocalBounds()
	return Rect(0, 0, bounds.Dx(), 30)
}

func (m *MainFrame) get_right_panel_bounds() Rectangle {
	bounds := m.LocalBounds()
	const panel_width = 200

	x := bounds.Dx() - panel_width
	y := m.get_toolbar_height()

	return Rect(x, y, bounds.Dx(), bounds.Dy())
}

func (m *MainFrame) get_main_panel_bounds() Rectangle {
	bounds := m.LocalBounds()
	const panel_width = 200

	x := panel_width
	y := m.get_toolbar_height()
	w := bounds.Dx() - panel_width
	h := bounds.Dy()

	return Rect(x, y, w, h)
}
