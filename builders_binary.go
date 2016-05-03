package img

import (
	"image"
	"sync"
)

//NewBinaryImage returns a BinaryImage with img's size
//if the color is greater than div, it's set up to true
func NewBinaryImage(img GrayImage, div uint8) BinaryImage {
	borders := img.Bounds()
	maxX, maxY := borders.Max.X, borders.Max.Y
	bi := make([][]bool, maxX)
	wg := new(sync.WaitGroup)
	wg.Add(maxX)
	for x := range bi {
		bi[x] = make([]bool, maxY)
		go readBinaryRow(wg, img, bi, x, maxY, div)
	}
	wg.Wait()
	return bi
}

func readBinaryRow(wg *sync.WaitGroup, from GrayImage, to BinaryImage, x, maxY int, div uint8) {
	for y := 0; y < maxY; y++ {
		to[x][y] = from.AtGray(x, y) < div
	}
	wg.Done()
}

//NewBinary builds a binary image from an image.Image
func NewBinary(img image.Image, div uint8) BinaryImage {
	borders := img.Bounds()
	x, y := borders.Max.X, borders.Max.Y
	bi := make(BinaryImage, x)
	wg := new(sync.WaitGroup)
	wg.Add(x)
	for x := range bi {
		bi[x] = make([]bool, y)
		readBinary(wg, img, bi, x, y, div)
	}
	wg.Wait()
	return bi
}

func readBinary(wg *sync.WaitGroup, from image.Image, to BinaryImage, x, maxY int, div uint8) {
	for y := 0; y < maxY; y++ {
		r, g, b, _ := from.ColorModel().Convert(from.At(x, y)).RGBA()
		gray := uint8(0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b))
		if gray >= div {
			to[x][y] = true
		}
	}
	wg.Done()
}
