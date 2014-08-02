package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	ft "gwk/freetype"
)

var fontfile = flag.String("fontfile", "./luxisr.ttf", "filename of the ttf font")

func printBounds(b ft.Bounds) {
	fmt.Printf("XMin:%d YMin:%d XMax:%d YMax:%d\n", b.XMin, b.YMin, b.XMax, b.YMax)
}

func printGlyph(g *ft.Glyph) {
	printBounds(g.Rect)
	fmt.Print("Points:\n---\n")
	e := 0
	log.Printf("%v", g.AllPoints)
	for i, p := range g.AllPoints {
		fmt.Printf("%4d, %4d", p.X, p.Y)
		if p.Flag&0x01 != 0 {
			fmt.Print("  on\n")
		} else {
			fmt.Print("  off\n")
		}
		if i+1 == int(g.EndIndexArray[e]) {
			fmt.Print("---\n")
			e++
		}
	}
}

func main() {
	flag.Parse()
	fmt.Printf("Loading fontfile %q\n", *fontfile)
	b, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	font, err := ft.Parse(b)
	if err != nil {
		log.Println(err)
		return
	}
	if font == nil {
		log.Println("Parse error.")
		return
	}
	fupe := font.FUnitsPerEm()
	printBounds(font.Bounds(fupe))
	fmt.Printf("FUnitsPerEm:%d\n\n", fupe)

	c0, c1 := 'A', 'V'

	i0 := font.Index(c0)
	hm := font.HMetric(fupe, i0)
	g := ft.NewGlyph()
	err = g.Load(font, fupe, i0, nil)

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("'%c' glyph\n", c0)
	fmt.Printf("AdvanceWidth:%d LeftSideBearing:%d\n", hm.AdvanceWidth, hm.LeftSideBearing)
	printGlyph(g)
	i1 := font.Index(c1)
	fmt.Printf("\n'%c', '%c' Kerning:%d\n", c0, c1, font.Kerning(fupe, i0, i1))
}
