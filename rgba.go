package img

import (
	"image"
	"image/color"
)

//ColorImage describes a RGBA image which implements draw.Image
type ColorImage [][]color.Color

//ColorModel returns the color model used by this implementation
func (i ColorImage) ColorModel() color.Model {
	return color.RGBAModel
}

//Bounds returns the size of the image
func (i ColorImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(i), len(i[0]))
}

//At returns a color in the position [x][y].
func (i ColorImage) At(x, y int) color.Color {
	return i[x][y]
}

//Set a color into the position [x][y]
func (i ColorImage) Set(x, y int, c color.Color) {
	i[x][y] = c
}
