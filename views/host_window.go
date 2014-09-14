// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	. "gwk/vango"
	"image"
)

func (h *HostWindow) OnHostPaint(native_context NativeContext, dirty_rect image.Rectangle) {
	native_canvas := NewNativeCanvas(dirty_rect)
	defer native_canvas.Release()

	draw_rect := dirty_rect.Sub(dirty_rect.Min)
	draw_context := GlobalDrawContext()
	draw_context.SetCanvas(native_canvas.Canvas)
	draw_context.DrawCanvas(draw_rect.Min.X, draw_rect.Min.Y,
		h.root_view.Canvas(), dirty_rect)

	native_canvas.BlitToNativeContext(native_context, dirty_rect.Min.X,
		dirty_rect.Min.Y, &draw_rect)
}

func (h *HostWindow) SetRootView(root_view *RootView) {
	h.root_view = root_view
	if root_view.HostWindow() != h {
		root_view.SetHostWindow(h)
	}
}

func (h *HostWindow) RootView() *RootView {
	return h.root_view
}
