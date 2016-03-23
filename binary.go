package img

import (
	"image"
	"image/color"
)

//BinaryImage is an image with logical values
type BinaryImage [][]bool

//At returns the related Color to the x,y position
func (bi BinaryImage) At(x, y int) color.Color {
	if bi[x][y] {
		return color.White
	}
	return color.Black
}

//ColorModel returns the model of color used in BinaryImage
func (bi BinaryImage) ColorModel() color.Model {
	return color.GrayModel
}

//Bounds returns the size of the binary image
func (bi BinaryImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(bi), len(bi[0]))
}

//Set wheter a black or white color into the image
func (bi BinaryImage) Set(x, y int, c color.Color) {
	switch c {
	case color.Black:
		bi[x][y] = false
	case color.White:
		bi[x][y] = true
	}
}

//SetBinary is a convenience method to set directly a logical value
func (bi BinaryImage) SetBinary(x, y int, color bool) {
	bi[x][y] = color
}

//Clone returns a new instance of the binary image itself
func (bi BinaryImage) Clone() (clone BinaryImage) {
	x, y := bi.Bounds().Max.X, bi.Bounds().Max.Y
	clone = make(BinaryImage, x)
	for ix := range bi {
		clone[ix] = make([]bool, y)
		for iy := range bi[ix] {
			clone[ix][iy] = bi[ix][iy]
		}
	}
	return
}
