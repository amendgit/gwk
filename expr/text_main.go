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
				"type":   "image_view",
				"left":   2,
				"top":    2,
				"width":  216,
				"height": 218,
				"color":  0xffffff,
				"delegate": views.UIMap{
					"on_mouse_enter": func(event *views.MouseEvent) {
						_, ok := &views.ImageView{}, false
						if _, ok = event.Owner.(*views.ImageView); !ok {
							return
						}
						// iv.SetColorRGB(byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)))
						// iv.ScheduleDraw()
					},
					"on_draw": func(event *views.DrawEvent) {
						ctxt := views.GlobalDrawContext()
						// ctxt.SelectCanvas(event.Canvas)
						ctxt.SetFontSize(25)
						ctxt.SetFontColor(0, 0, 255)
						ctxt.DrawText("Hello, GWK!", image.Rect(30, 50, 200, 200))
						ctxt.SetStrokeColor(255, 0, 0)
						ctxt.StrokeRect(image.Rect(20, 30, 180, 70))
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
