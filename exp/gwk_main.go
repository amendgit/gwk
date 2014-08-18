// Copyright 2014 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package main

import (
	"gwk"
	. "gwk/views"
	"image"
	"math/rand"
	"time"
)

func makeMainUIMap() UIMap {
	var main_uimap = UIMap{
		"type": "main_frame",
		"toolbar": UIMap{
			"type":   "toolbar",
			"width":  "fill_parent",
			"height": "20",
		},
		"left_panel": UIMap{
			"type":   "view",
			"id":     "left_panels_container",
			"width":  "fill_parent",
			"height": "fill_parent",
			"layout": "vertical",
			"children": []UIMap{
				{
					"type": "panel",
				},
				{
					"type": "panel",
				},
			},
		},
		"main_panel": UIMap{
			"type":       "image_view",
			"color":      0x606060,
			"image_resc": "image_show_cat",
			// "children": []UIMap{
			// 	{
			// 		"type":   "image_view",
			// 		"left":   50,
			// 		"top":    120,
			// 		"width":  570,
			// 		"height": 300,
			// 	},
			// },
		},
		"right_panel": UIMap{
			"type":   "view",
			"id":     "left_panels_container",
			"width":  "fill_parent",
			"height": "fill_parent",
			"layout": "vertical",
			"children": []UIMap{
				{
					"type": "panel",
				},
				{
					"type": "panel",
				},
				{
					"type": "button",
					// "action": func() {
					// 	return
					// },
				},
			},
		},
	}

	return main_uimap
}

func main() {
	rand.Seed(time.Now().Unix())

	gwk.Init()

	var host_view = NewHostView(image.Rect(0, 0, 1159, 687))

	host_view.RootView.AddChild(MockUp(makeMainUIMap()))

	host_view.Show()

	CurrentUIEventLoop().Run()
}
