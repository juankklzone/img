package img

import (
	"image"
	"image/color"
)

//GrayImage describes a grayscale image which implements draw.Image
type GrayImage [][]color.Gray

//ColorModel returns the color model used by this implementation
func (gi GrayImage) ColorModel() color.Model {
	return color.GrayModel
}

//Bounds returns the size of the image
func (gi GrayImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(gi), len(gi[0]))
}

//At returns a color in the position [x][y]
func (gi GrayImage) At(x, y int) color.Color {
	r, g, b, a := gi[x][y].RGBA()
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

//SetGray is a convenience method to avoid Color wrapping
func (gi GrayImage) SetGray(x, y int, gray uint8) {
	gi[x][y].Y = gray
}

//AtGray is a convenience method to avoid Color wrapping
func (gi GrayImage) AtGray(x, y int) uint8 {
	return gi[x][y].Y
}

//Set a color into the position [x][y]
func (gi GrayImage) Set(x, y int, c color.Color) {
	r, g, b, _ := gi[x][y].RGBA()
	gi[x][y] = color.Gray{uint8((r + g + b) / 3)}
}

//Clone returns a new instance of the gray image itself.
func (gi GrayImage) Clone() (i GrayImage) {
	x, y := gi.Bounds().Max.X, gi.Bounds().Max.Y
	i = make(GrayImage, x)
	for ix := range gi {
		i[ix] = make([]color.Gray, y)
		for iy := range gi[ix] {
			i[ix][iy].Y = gi[ix][iy].Y
		}
	}
	return i
}
