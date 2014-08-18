// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	//. "gwk/vango"
	//"gwk/views/resc"
	. "image"
	//"log"
)

type Panel struct {
	BaseView
}

const (
	kPanelBorderSize   = 10
	kPanelHeaderHeight = 20
)

func NewPanel() *Panel {
	new_panel := new(Panel)
	// new_panel.header = resc.FindCanvasByID("panel_header")
	return new_panel
}

func (p *Panel) DrawPanelHeader(event *DrawEvent) {
	header_rect := p.get_header_bounds()
	event.Canvas.FillRect(header_rect, 19, 19, 19)
}

func (p *Panel) DrawPanelBorder(event *DrawEvent) {
	var r Rectangle
	r = p.get_left_border_bounds()
	event.Canvas.FillRect(r, 19, 19, 19)

	r = p.get_right_border_bounds()
	event.Canvas.FillRect(r, 19, 19, 19)

	r = p.get_bottom_border_bounds()
	event.Canvas.FillRect(r, 19, 19, 19)
}

func (p *Panel) DrawPanelContentBackground(event *DrawEvent) {
	r := p.get_content_bounds()
	event.Canvas.FillRect(r, 96, 96, 96)
}

func (p *Panel) OnDraw(event *DrawEvent) {
	p.DrawPanelHeader(event)
	p.DrawPanelBorder(event)
	p.DrawPanelContentBackground(event)
	// event.Canvas.DrawColor(19, 19, 19)
	// event.Canvas.DrawCanvas(0, 0, p.header, nil)
	// header_rect := Rect(0, 0, p.Width(), 30)
	// event.Canvas.StretchDraw(header_rect, p.header)
}

// ============================================================================

func (p *Panel) get_header_bounds() Rectangle {
	return Rect(0, 0, p.LocalBounds().Dx(), kPanelHeaderHeight)
}

func (p *Panel) get_left_border_bounds() Rectangle {
	bounds := p.LocalBounds()
	return Rect(0, kPanelHeaderHeight, kPanelBorderSize,
		bounds.Dy()-kPanelBorderSize)
}

func (p *Panel) get_right_border_bounds() Rectangle {
	bounds := p.LocalBounds()
	return Rect(bounds.Dx()-kPanelBorderSize, kPanelHeaderHeight, bounds.Dy(),
		bounds.Dy()-kPanelBorderSize)
}

func (p *Panel) get_bottom_border_bounds() Rectangle {
	bounds := p.LocalBounds()
	return Rect(0, bounds.Dy()-kPanelBorderSize, bounds.Dx(), bounds.Dy())
}

func (p *Panel) get_content_bounds() Rectangle {
	bounds := p.LocalBounds()
	return Rect(kPanelBorderSize, kPanelHeaderHeight,
		bounds.Dx()-kPanelBorderSize, bounds.Dy()-kPanelBorderSize)
}
