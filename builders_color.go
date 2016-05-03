package img

import (
	"image"
	"image/color"
	"sync"
)

//NewColorFromImage recieves an image and returns a drawable RGBA image
func NewColorFromImage(img image.Image) ColorImage {
	borders := img.Bounds()
	maxX, maxY := borders.Max.X, borders.Max.Y
	ci := make(ColorImage, maxX)
	wg := new(sync.WaitGroup)
	wg.Add(maxX)
	for idx := range ci {
		ci[idx] = make([]color.Color, maxY)
		go readColorRow(wg, img, ci, idx, maxY)
	}
	wg.Wait()
	return ci
}

//NewEmptyImage creates a rbga image which size is equal to img
func NewEmptyImage(borders image.Rectangle) ColorImage {
	maxX, maxY := borders.Max.X, borders.Max.Y
	ci := make([][]color.Color, maxX)
	for idx := range ci {
		ci[idx] = make([]color.Color, maxY)
		for idy := range ci[idx] {
			ci[idx][idy] = color.RGBA{}
		}
	}
	return ci
}

func readColorRow(wg *sync.WaitGroup, from image.Image, to ColorImage, x, maxY int) {
	for y := 0; y < maxY; y++ {
		to[x][y] = to.ColorModel().Convert(from.At(x, y))
	}
	wg.Done()
}
