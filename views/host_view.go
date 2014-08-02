// Copyright 2013 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	"image"
)

type HostView struct {
	*HostWindow
	*RootView
}

func NewHostView(bounds image.Rectangle) *HostView {
	var hv HostView

	hv.HostWindow = NewHostWindow(bounds)
	hv.RootView = NewRootView(hv.ClientBounds())
	hv.HostWindow.SetRootView(hv.RootView)

	return &hv
}

func (hv *HostView) Show() {
	hv.HostWindow.Show()
}
