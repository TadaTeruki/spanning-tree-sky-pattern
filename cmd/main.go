package main

import (
	"image/color"
	"image/png"
	"math"
	"os"

	skypattern "github.com/TadaTeruki/spanning-tree-sky-pattern"
)

func main() {

	image := skypattern.GeneratePattern(
		0,
		1024, 1024,
		&color.RGBA{70, 140, 200, 255},
		&color.RGBA{230, 245, 255, 255},
		func(depth int) float64 {
			return 1.0 - math.Sin(float64(depth)/500.)
		},
	)

	file, _ := os.Create("image.png")
	png.Encode(file, image)

	file.Close()
}
