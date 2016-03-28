package img

import (
	"image"
	"image/color"
)

//NewColorFromImage recieves an image and returns a drawable RGBA image
func NewColorFromImage(img image.Image) ColorImage {
	borders := img.Bounds()
	maxX, maxY := borders.Max.X, borders.Max.Y
	ci := make(ColorImage, maxX)
	transf := ci.ColorModel()
	for idx := range ci {
		ci[idx] = make([]color.Color, maxY)
		for idy := range ci[idx] {
			ci[idx][idy] = transf.Convert(img.At(idx, idy))
		}
	}
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
	for idx := range ci {
		ci[idx] = make([]color.Gray, maxY)
		for idy := range ci[idx] {
			r, g, b, _ := img.ColorModel().Convert(img.At(idx, idy)).RGBA()
			ci[idx][idy] = color.Gray{
				Y: uint8(0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b)),
			}
		}
	}
	return ci
}

//NewBinaryImage returns a BinaryImage with img's size
//if the color is greater than div, it's set up to true
func NewBinaryImage(img GrayImage, div uint8) BinaryImage {
	borders := img.Bounds()
	maxX, maxY := borders.Max.X, borders.Max.Y
	bi := make([][]bool, maxX)
	for x := range bi {
		bi[x] = make([]bool, maxY)
		for y := range bi[x] {
			if img.AtGray(x, y) > div {
				bi[x][y] = true
			}
		}
	}
	return bi
}
