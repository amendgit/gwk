// Copyright 2014 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package main

import (
	"gwk"
	"gwk/views"
	"image"
	"math/rand"
	"time"
)

func make_main_ui_map() views.UIMap {
	var main_ui_map = views.UIMap{
		"type":  "image_view",
		"color": 0xffffff,
		"children": []views.UIMap{
			views.UIMap{
				"type": "button",
				"left": 50,
				"top":  50,
			},
		},
	}

	return main_ui_map
}

func main() {
	rand.Seed(time.Now().Unix())
	gwk.Init()

	host_view := views.NewHostView(image.Rect(0, 0, 235, 260))
	host_view.RootView.AddChild(views.MockUp(make_main_ui_map()))
	host_view.Show()

	views.MainUIEventLoop().Run()
}
