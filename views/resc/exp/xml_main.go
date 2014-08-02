package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type CanvasResc struct {
	id       string
	canvas   *Canvas
	filename string
}

func main() {
	file_content, err := ioutil.ReadFile("./resc.xml")
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
		g_id_canvas_map[id] = &resc
		t, _ = d.Token()
		if end, ok := t.(xml.EndElement); ok {
			if end.Name.Local != "Image" {

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
