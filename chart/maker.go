package chart

import (
	"image"
	"image/color"
	"image/gif"
	"io"
)

// Maker is chart maker
type Maker interface {
	Make(List, io.Writer)
}

func abs(x float64) float64 {
	if x < 0 {
		return 0 - x
	}
	return x
}

// GifMaker make a gif file of trend metric
type GifMaker struct {
	width  float64
	height float64
}

// Make creates the concrete file
func (d GifMaker) Make(list List, out io.Writer) {
	points := list.CreatePoints(d.width, d.height)
	anim := gif.GIF{LoopCount: -1}
	palette := []color.Color{color.White, color.Black}

	for n := range points {
		rect := image.Rect(0, 0, int(d.width), int(d.height))
		img := image.NewPaletted(rect, palette)

		points[0 : n+1].connect(func(x, y float64) {
			img.SetColorIndex(int(x), int(y), 1)
		})
		anim.Delay = append(anim.Delay, 10)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}

// NewGifMaker creates a new GifMaker
func NewGifMaker(width, height float64) GifMaker {
	return GifMaker{width, height}
}
