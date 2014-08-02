package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	ft "gwk/freetype"
)

func main() {
	const (
		n = 17
		r = 256 * 80
	)
	s := ft.Fix32(r * math.Sqrt(2) / 2)
	t := ft.Fix32(r * math.Tan(math.Pi/8))

	m := image.NewRGBA(image.Rect(0, 0, 800, 600))
	draw.Draw(m, m.Bounds(), image.NewUniform(color.RGBA{63, 63, 63, 255}), image.ZP, draw.Src)
	mp := ft.NewRGBADrawer(m)
	mp.SetColor(image.Black)
	z := ft.NewRast(800, 600)

	for i := 0; i < n; i++ {
		cx := ft.Fix32(25600 + 51200*(i%4))
		cy := ft.Fix32(2560 + 32000*(i/4))
		c := ft.RastPoint{X: cx, Y: cy}
		theta := math.Pi * (0.5 + 0.5*float64(i)/(n-1))
		dx := ft.Fix32(r * math.Cos(theta))
		dy := ft.Fix32(r * math.Sin(theta))
		d := ft.RastPoint{X: dx, Y: dy}
		// Draw a quarter-circle approximated by two quadratic segments,
		// with each segment spanning 45 degrees.
		z.Start(c)
		z.Add1(c.Add(ft.RastPoint{X: r, Y: 0}))
		z.Add2(c.Add(ft.RastPoint{X: r, Y: t}), c.Add(ft.RastPoint{X: s, Y: s}))
		z.Add2(c.Add(ft.RastPoint{X: t, Y: r}), c.Add(ft.RastPoint{X: 0, Y: r}))
		// Add another quadratic segment whose angle ranges between 0 and 90 degrees.
		// For an explanation of the magic constants 22, 150, 181 and 256, read the
		// comments in the freetype/ft package.
		dot := 256 * d.Dot(ft.RastPoint{X: 0, Y: r}) / (r * r)
		multiple := ft.Fix32(150 - 22*(dot-181)/(256-181))
		z.Add2(c.Add(ft.RastPoint{X: dx, Y: r + dy}.Mul(multiple)), c.Add(d))
		// Close the curve.
		z.Add1(c)
	}
	z.Rast(mp)

	for i := 0; i < n; i++ {
		cx := ft.Fix32(25600 + 51200*(i%4))
		cy := ft.Fix32(2560 + 32000*(i/4))
		for j := 0; j < n; j++ {
			theta := math.Pi * float64(j) / (n - 1)
			dx := ft.Fix32(r * math.Cos(theta))
			dy := ft.Fix32(r * math.Sin(theta))
			m.Set(int((cx+dx)/256), int((cy+dy)/256), color.RGBA{255, 255, 0, 255})
		}
	}

	// Save that RGBA image to disk.
	f, err := os.Create("round.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println("Wrote round.png OK.")
}
