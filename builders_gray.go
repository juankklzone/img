package img

import (
	"image"
	"image/color"
	"sync"
)

//NewEmptyGrayImage creates an grayscale image with a size equals to the recieved Rectangle
func NewEmptyGrayImage(borders image.Rectangle) GrayImage {
	x, y := borders.Max.X, borders.Max.Y
	gi := make(GrayImage, x)
	for ix := range gi {
		gi[ix] = make([]color.Gray, y)
	}
	return gi
}

//NewGrayImage recieves an image and returns a drawable grayscale image
func NewGrayImage(img image.Image) GrayImage {
	borders := img.Bounds()
	maxX, maxY := borders.Max.X, borders.Max.Y
	ci := make([][]color.Gray, maxX)
	wg := new(sync.WaitGroup)
	wg.Add(maxX)
	for idx := range ci {
		ci[idx] = make([]color.Gray, maxY)
		go readGrayRow(wg, img, ci, idx, maxY)
	}
	wg.Wait()
	return ci
}

func readGrayRow(wg *sync.WaitGroup, from image.Image, to GrayImage, x, maxY int) {
	for y := 0; y < maxY; y++ {
		r, g, b, _ := from.ColorModel().Convert(from.At(x, y)).RGBA()
		to[x][y] = color.Gray{
			Y: uint8(0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b)),
		}
	}
	wg.Done()
}
