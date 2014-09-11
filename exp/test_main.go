// Copyright 2014 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package main

import (
	"gwk"
	// "gwk/vango"
	"gwk/views"
	"image"
	"math/rand"
	"time"
)

func make_main_ui_map() views.UIMap {
	var main_ui_map = views.UIMap{
		"type": "base_view",
		"children": []views.UIMap{
			views.UIMap{
				"type":   "test_view",
				"id":     "frame",
				"left":   2,
				"top":    2,
				"width":  780,
				"height": 558,
				"children": []views.UIMap{
					{
						"type":   "test_view",
						"id":     "view0",
						"left":   10,
						"top":    10,
						"width":  100,
						"height": 80,
					},
					{
						"type":   "test_view",
						"id":     "view1",
						"left":   200,
						"top":    2,
						"width":  500,
						"height": 200,
					},
					{
						"type":   "test_view",
						"id":     "view2",
						"left":   500,
						"top":    300,
						"width":  200,
						"height": 200,
					},
					{
						"type":   "test_view",
						"id":     "view3",
						"left":   20,
						"top":    250,
						"width":  300,
						"height": 250,
					},
				},
			},
		},
	}

	return main_ui_map
}

func main() {
	rand.Seed(time.Now().Unix())
	gwk.Init()

	host_view := views.NewHostView(image.Rect(0, 0, 800, 600))
	host_view.RootView.AddChild(views.MockUp(make_main_ui_map()))
	host_view.Show()

	views.MainUIEventLoop().Run()
}
