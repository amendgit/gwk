// Copyright 2014 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resc

import (
	"bytes"
	"encoding/xml"
	. "gwk/vango"
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

type Resc interface {
}

type CanvasResc struct {
	id       string
	canvas   *Canvas
	filename string
}

type Canvas3Resc struct {
	id       string
	canvas3  *Canvas3
	filename string
}

// type ColorResc struct {
// 	id    string
// 	color Color
// }

var g_id_resc_map map[string]Resc = make(map[string]Resc)

func LoadRescFile(filename string) {
	file_content, err := ioutil.ReadFile("./resc/resc.xml")
	d := xml.NewDecoder(bytes.NewBuffer(file_content))
	t, err := d.Token()

	load_image_resc := func(start xml.StartElement) {
		var id string
		for _, attr := range start.Attr {
			if attr.Name.Local == "id" {
				id = attr.Value
			}
		}
		var filename string
		t, _ := d.Token()
		if char_data, ok := t.(xml.CharData); ok {
			filename = string([]byte(char_data))
		}
		var resc CanvasResc
		resc.filename = filename
		g_id_resc_map[id] = &resc
		t, _ = d.Token()
		if end, ok := t.(xml.EndElement); ok {
			if end.Name.Local != "Image" {
				// resturn error
			}
		}
	}

	for err == nil {
		switch token := t.(type) {
		case xml.StartElement:
			if token.Name.Local == "Image" {
				load_image_resc(token)
			}
		default:
		}
		t, err = d.Token()
	}
}

func FindCanvasByID(id string) *Canvas {
	resc := g_id_resc_map[id]
	if canvas_resc, ok := resc.(*CanvasResc); ok {
		if canvas_resc.canvas == nil {
			canvas_resc.canvas = LoadCanvas("resc/" + canvas_resc.filename)
		}
		return canvas_resc.canvas
	}
	return nil
}

func LoadCanvas(filename string) *Canvas {
	var fd, err = os.Open(filename)
	if err != nil {
		return nil
	}
	defer fd.Close()
	return LoadCanvasFile(fd)
}

func LoadCanvasFile(fd *os.File) *Canvas {
	var png, err = png.Decode(fd)
	if err != nil {
		return nil
	}
	if nrgba, ok := png.(*image.NRGBA); ok {
		var canvas = NewCanvas(nrgba.Rect.Dx(), nrgba.Rect.Dy())
		canvas.DrawImageNRGBA(0, 0, nrgba, nil)
		canvas.SetOpaque(true)
		return canvas
	}
	if rgba, ok := png.(*image.RGBA); ok {
		var canvas = NewCanvas(rgba.Rect.Dx(), rgba.Rect.Dy())
		canvas.DrawImageRGBA(0, 0, rgba, nil)
		return canvas
	}
	return nil
}
