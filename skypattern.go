/*

MIT License

Copyright (c) 2022 Tada Teruki (多田 瑛貴)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package skypattern

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"

	pq "github.com/TadaTeruki/PriorityQueueGo/PriorityQueue"
)

var direction = [][]int{
	{0, -1}, {0, 1}, {-1, 0}, {1, 0},
}

type Path struct {
	isx, isy, iex, iey int
}

func mapColorUnit(m float64, a, b uint8) uint8 {
	return uint8(math.Max(math.Min(float64(a)*m+float64(b)*(1.-m), 1.0), 0.0))
}

func mapColor(m float64, ca, cb *color.RGBA) *color.RGBA {
	nc := new(color.RGBA)
	nc.R = mapColorUnit(m, ca.R, cb.R)
	nc.G = mapColorUnit(m, ca.G, cb.G)
	nc.B = mapColorUnit(m, ca.B, cb.B)
	nc.A = mapColorUnit(m, ca.A, cb.A)
	return nc
}

func getWeight(ord int) float64 {
	return 1.0 - math.Sin(float64(ord)/500)
}

func GetPattern(seed int64, width, height int, color_a, color_b *color.RGBA) {

	rand.Seed(seed)

	depth := make([][]int, height)

	for iy := 0; iy < height; iy++ {
		depth[iy] = make([]int, width)
		for ix := 0; ix < width; ix++ {
			depth[iy][ix] = -1
		}
	}

	startx, starty := 0, 0

	depth[starty][startx] = 0

	var nextpath pq.PriorityQueue
	nextpath.Push(pq.MakeObject(Path{startx, starty, startx, starty}, rand.Float64()))

	for {
		if nextpath.GetSize() == 0 {
			break
		}

		p := nextpath.GetFront().Value.(Path)

		nextpath.PopFront()

		for i := 0; i < len(direction); i++ {

			nx := p.iex + direction[i][0]
			ny := p.iey + direction[i][1]

			if nx < 0 {
				nx = width - 1
			}
			if ny < 0 {
				ny = height - 1
			}
			if nx >= width {
				nx = 0
			}
			if ny >= height {
				ny = 0
			}

			if depth[ny][nx] < 0 {
				nextpath.Push(pq.MakeObject(Path{p.iex, p.iey, nx, ny}, rand.Float64()))
			}

		}

		depth[p.iey][p.iex] = depth[p.isy][p.isx] + 1
	}

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	var wg sync.WaitGroup

	wg.Add(width * height)

	for iy := 0; iy < height; iy++ {
		for ix := 0; ix < width; ix++ {
			go func(x, y int) {
				weight := math.Min(math.Max(getWeight(depth[y][x]), 0.0), 1.0)
				color := mapColor(weight, color_a, color_b)
				img.Set(x, y, color)
				wg.Done()
			}(ix, iy)
		}
	}

	wg.Wait()

	file, _ := os.Create("image.png")
	png.Encode(file, img)

	file.Close()
}
