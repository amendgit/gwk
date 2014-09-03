// Copyright 2014 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package main

import (
	"gwk"
	"gwk/vango"
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
				"left":   2,
				"top":    2,
				"width":  216,
				"height": 218,
				"color":  0xffffff,
				"delegate": views.UIMap{
					"on_mouse_enter": func(event *views.MouseEvent) {
						iv, ok := &views.ImageView{}, false
						if iv, ok = event.Owner.(*views.ImageView); !ok {
							return
						}
						iv.SetColorRGB(byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)))
						iv.ScheduleDraw()
					},
					"on_draw": func(event *views.DrawEvent) {
						ctxt := vango.NewContext()
						ctxt.SelectCanvas(event.Canvas)
						ctxt.DrawText("GWK", image.Rect(0, 0, 50, 50))
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
