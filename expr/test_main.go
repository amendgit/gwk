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
					/// row 0
					{
						"type":   "test_view",
						"id":     "view-0-0",
						"left":   10,
						"top":    10,
						"width":  220,
						"height": 200,
						// "children": []views.UIMap{
						// 	{
						// 		"type":   "test_view",
						// 		"id":     "view-0-0-0",
						// 		"left":   10,
						// 		"width":  200,
						// 		"top":    10,
						// 		"height": 180,
						// 	},
						// },
					},
					/// row 1
					{
						"type":   "test_view",
						"id":     "view-0-1",
						"left":   240,
						"top":    10,
						"width":  532,
						"height": 200,
					},
					{
						"type":   "test_view",
						"id":     "view-1-0",
						"left":   10,
						"top":    220,
						"width":  300,
						"height": 200,
					},
					{
						"type":   "test_view",
						"id":     "view-1-1",
						"left":   320,
						"top":    220,
						"width":  252,
						"height": 200,
					},
					/// row 2
					{
						"type":   "test_view",
						"id":     "view-1-2",
						"left":   582,
						"top":    220,
						"width":  190,
						"height": 200,
					},
					{
						"type":   "test_view",
						"id":     "view-2-1",
						"left":   10,
						"top":    430,
						"width":  762,
						"height": 120,
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
