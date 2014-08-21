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
		"type": "base_view",
		"children": []views.UIMap{
			views.UIMap{
				"type":   "image_view",
				"left":   10,
				"top":    10,
				"width":  200,
				"height": 200,
				"color":  0xffffff,
				"delegate": views.UIMap{
					"on_mouse_enter": func(event *views.MouseEvent) {
						var funny_task func()
						funny_task = func() {
							iv, _ := event.Owner.(*views.ImageView)
							iv.SetColorRGB(byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)))
							iv.ScheduleDraw()
							views.MainUIEventLoop().PostDelayedTask(funny_task, 1000)
						}
						views.MainUIEventLoop().PostDelayedTask(funny_task, 1000)
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

	host_view := views.NewHostView(image.Rect(0, 0, 235, 260))
	host_view.RootView.AddChild(views.MockUp(make_main_ui_map()))
	host_view.Show()

	views.MainUIEventLoop().Run()
}
