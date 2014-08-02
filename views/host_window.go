// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	. "gwk/vango"
	"image"
)

func (h *HostWindow) OnHostPaint(ctxt NativeContext, dirty_rect image.Rectangle) {
	var canvas = NewNativeCanvas(dirty_rect)
	defer canvas.Release()

	var canvas_rect = dirty_rect.Sub(dirty_rect.Min)

	canvas.DrawCanvas(canvas_rect.Min.X, canvas_rect.Min.Y,
		h.root_view.Canvas(), &dirty_rect)
	canvas.BlitToContext(ctxt, dirty_rect.Min.X, dirty_rect.Min.Y, &canvas_rect)
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
