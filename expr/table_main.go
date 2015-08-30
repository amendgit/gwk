// Copyright 2014 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package main

import (
	"fmt"
)

func NewWorkAreaView() string {
	return "WorkAreaView"
}

func layout_root_view(a int) {

}

type UIMap map[string]interface{}

func makeLeftPanel() UIMap {
	var m = UIMap{
		"type":   "grid_view",
		"layout": "vertical",
	}
	return m
}

func main() {
	var ui = UIMap{
		"layout": "horizontal",
		"children": []UIMap{
			{
				"type":   "button",
				"layout": "horizontal",
			},
			{
				"type":   "image_view",
				"layout": "horizontal",
			},
			{
				"type":        "custom_view",
				"custom_view": NewWorkAreaView(),
			},
			{
				"type":   "list_view",
				"layout": layout_root_view,
			},
			makeLeftPanel(),
		},
	}

	fmt.Printf("%v", ui)
}
