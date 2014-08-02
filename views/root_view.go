// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	. "gwk/vango"
	. "image"
	"log"
)

type RootView struct {
	View
	host_window *HostWindow

	mouse_move_handler Viewer
}

func NewRootView(bounds Rectangle) *RootView {
	var v = new(RootView)
	v.SetID("cn.ustc.edu/gwk/root_view")
	v.SetBounds(bounds)
	return v
}

func (r *RootView) AddChild(child Viewer) {
	log.Printf("RootView.AddChild %v + %v", r.ID(), child.ID())
	r.children = append(r.children, child)
	child.SetParent(r)
}

func (r *RootView) SetHostWindow(h *HostWindow) {
	r.host_window = h
	if h.RootView() != r {
		h.SetRootView(r)
	}
}

func (r *RootView) ToAbsPoint(pt Point) Point {
	return pt
}

func (r *RootView) ToAbsRect(rc Rectangle) Rectangle {
	return rc
}

func (r *RootView) ToDevicePoint(pt Point) Point {
	log.Printf("NOTIMPLEMENT")
	return pt
}

func (r *RootView) ToDeviceRect(rc Rectangle) Rectangle {
	log.Printf("NOTIMPLEMENT")
	return rc
}

func (r *RootView) HostWindow() *HostWindow {
	return r.host_window
}

func (r *RootView) Canvas() *Canvas {
	if r.canvas == nil {
		r.canvas = NewCanvas(r.W(), r.H())
	}
	canvas_bounds := r.canvas.Bounds()
	if canvas_bounds.Dx() < r.W() || canvas_bounds.Dy() < r.H() {
		r.canvas = NewCanvas(r.W(), r.H())
	} else {
		return r.canvas.SubCanvas(Rect(0, 0, r.W(), r.H()))
	}
	return r.canvas
}

func (r *RootView) DispatchDraw(dirty_rect Rectangle) {
	children := r.Children()
	if children != nil && len(children) < 1 {
		return
	}
	var dispatch_draw func(event *DrawEvent)
	dispatch_draw = func(event *DrawEvent) {
		view := event.Owner
		bounds := view.Bounds()
		dirty_rect := event.DirtyRect.Intersect(bounds.Sub(bounds.Min))
		if dirty_rect.Empty() {
			return
		}

		if view == nil {
			return
		}

		view.OnDraw(event)

		view_canvas := event.Canvas
		for _, child := range view.Children() {
			// Caculate the child dirty rectangle.
			child_dirty_rect := dirty_rect.Intersect(child.Bounds())
			if child_dirty_rect.Empty() {
				continue
			}
			child_dirty_rect = child_dirty_rect.Sub(child_dirty_rect.Min)

			// Clip the canvas to child bounds.
			child_canvas := view_canvas.SubCanvas(child.Bounds())

			// Make a new draw event.
			child_draw_event := &DrawEvent{
				Owner:     child,
				DirtyRect: child_dirty_rect,
				Canvas:    child_canvas,
			}

			// Dispatch draw.
			dispatch_draw(child_draw_event)
		}
	}

	// RootView only have one child. That's the MainFrame.
	event := &DrawEvent{
		Owner:     children[0],
		DirtyRect: dirty_rect,
		Canvas:    r.Canvas(),
	}
	dispatch_draw(event)
}

func DispatchLayout(v Viewer) {
	if v.Layouter() != nil {
		v.Layouter().Layout(v)
	}

	for _, child := range v.Children() {
		DispatchLayout(child)
	}
}

func (r *RootView) DispatchLayout() {
	new_rect := r.Bounds()

	if r.Children() == nil {
		return
	}

	r.Children()[0].SetXYWH(0, 0, new_rect.Dx(), new_rect.Dy())
	DispatchLayout(r.Children()[0])

	r.DispatchDraw(r.Bounds())
}

func get_event_handler_for_point(v Viewer, pt Point) Viewer {
	pt.X, pt.Y = pt.X-v.X(), pt.Y-v.Y()

	for _, child := range v.Children() {
		rect := child.Bounds()
		if rect.Min.X < pt.X && rect.Min.Y < pt.Y && rect.Max.X > pt.X &&
			rect.Max.Y > pt.Y {
			return get_event_handler_for_point(child, pt)
		}
	}

	return v
}

func (r *RootView) DispatchMouseMove(pt Point) {
	v := get_event_handler_for_point(r, pt)

	// for v != r.mouse_move_handler {
	// 	v = v.Parent()
	// }

	mouse_event := NewMouseEvent(pt)

	if v != nil && v != r && v != r.mouse_move_handler {
		old_handler := r.mouse_move_handler
		r.mouse_move_handler = v
		if old_handler != nil {
			old_handler.OnMouseLeave(mouse_event)
		}

		if r.mouse_move_handler != nil {
			r.mouse_move_handler.OnMouseEnter(mouse_event)
		}
	}

	if r.mouse_move_handler != nil {
		// r.mouse_move_handler.OnMouseMove()
	}
}

func (r *RootView) UpdateRect(rect Rectangle) {
	rect = rect.Intersect(r.Bounds())
	r.DispatchDraw(rect)
	r.host_window.InvalidateRect(rect)
}
